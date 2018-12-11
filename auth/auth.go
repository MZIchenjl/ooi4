package auth

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type KancolleAuth struct {
	http.Client
	header       http.Header
	worldID      int
	APIStartTime int64
	st           string
	loginID      string
	password     string
	dmmToken     string
	token        string
	idKey        string
	owner        string
	osapiURL     string
	pwdKey       string
	WorldIP      string
	APIToken     string
	entry        string
}

func New(id, pass string) *KancolleAuth {
	client := new(KancolleAuth)

	client.loginID = id
	client.password = pass

	jar, _ := cookiejar.New(nil)
	client.Jar = jar

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	client.Transport = transport

	client.Timeout = 10 * time.Second
	return client
}

func (self *KancolleAuth) getDMMTokens() error {
	req, _ := http.NewRequest(http.MethodGet, urlLayouts["login"], nil)
	req.Header.Set("User-Agent", userAgent)
	res, err := self.Do(req)
	if err != nil {
		return errors.New("连接DMM登录页失败")
	}
	defer res.Body.Close()
	var find [][]byte
	br := bufio.NewReader(res.Body)
	for {
		line, _, err := br.ReadLine()
		if self.dmmToken != "" && self.token != "" {
			return nil
		}
		if err == io.EOF {
			break
		}
		if self.dmmToken == "" {
			find = patterns["dmm_token"].FindSubmatch(line)
			if len(find) != 0 {
				self.dmmToken = string(find[1])
			}
		}
		if self.token == "" {
			find = patterns["token"].FindSubmatch(line)
			if len(find) != 0 {
				self.token = string(find[1])
			}
		}
	}
	if self.dmmToken == "" {
		return errors.New("获取DMM token失败")
	}
	if self.token == "" {
		return errors.New("获取token失败")
	}
	return nil
}

func (self *KancolleAuth) getAjaxToken() error {
	data, _ := json.Marshal(map[string]string{
		"token": self.token,
	})
	req, _ := http.NewRequest(http.MethodPost, urlLayouts["ajax"], bytes.NewReader(data))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", urlLayouts["login"])
	req.Header.Set("http-dmm-token", self.dmmToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://accounts.dmm.com")
	res, err := self.Do(req)
	if err != nil {
		return errors.New("DMM登录页AJAX请求失败")
	}
	defer res.Body.Close()
	j := new(struct {
		Header map[string]string `json:"header"`
		Body   map[string]string `json:"body"`
	})
	err = json.NewDecoder(res.Body).Decode(j)
	if err != nil || j.Body == nil {
		return errors.New("DMM修改登录机制了，请通知管理员处理")
	}
	self.token = j.Body["token"]
	self.idKey = j.Body["login_id"]
	self.pwdKey = j.Body["password"]
	if self.token == "" || self.idKey == "" || self.pwdKey == "" {
		return errors.New("DMM修改登录机制了，请通知管理员处理")
	}
	return nil
}

func (self *KancolleAuth) getOSAPIURL() error {
	data := make(url.Values)
	data.Set("login_id", self.loginID)
	data.Set("password", self.password)
	data.Set("token", self.token)
	data.Set("idKey", self.loginID)
	data.Set("pwKey", self.password)
	req, _ := http.NewRequest(http.MethodPost, urlLayouts["auth"], strings.NewReader(data.Encode()))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", urlLayouts["login"])
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://accounts.dmm.com")
	res, err := self.Do(req)
	if err != nil {
		return errors.New("连接DMM认证网页失败")
	}
	defer res.Body.Close()
	var findIdx []int
	br := bufio.NewReader(res.Body)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		findIdx = patterns["reset"].FindIndex(line)
		if len(findIdx) != 0 {
			return errors.New("DMM强制要求用户修改密码")
		}
	}
	req, _ = http.NewRequest(http.MethodGet, urlLayouts["game"], nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", urlLayouts["login"])
	req.Header.Set("Origin", "https://accounts.dmm.com")
	resp, err := self.Do(req)
	if err != nil {
		return errors.New("连接舰队collection游戏页面失败")
	}
	defer resp.Body.Close()
	var find [][]byte
	br = bufio.NewReader(resp.Body)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		find = patterns["osapi"].FindSubmatch(line)
		if len(find) != 0 {
			self.osapiURL = string(find[1])
			return nil
		}
	}
	return errors.New("用户名或密码错误，请重新输入")
}

type worldInfoResData struct {
	APIWorldId int `json:"api_world_id"`
}

type worldInfoRes struct {
	APIResult    int              `json:"api_result"`
	APIResultMsg string           `json:"api_result_msg"`
	APIData      worldInfoResData `json:"api_data"`
}

func (self *KancolleAuth) getWorld() error {
	qs, _ := url.Parse(self.osapiURL)
	q := qs.Query()
	self.owner = q.Get("owner")
	self.st = q.Get("st")
	u := fmt.Sprintf(urlLayouts["get_world"], self.owner, time.Now().UnixNano()/1e6)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", self.osapiURL)
	req.Header.Set("Origin", "https://accounts.dmm.com")
	res, err := self.Do(req)
	if err != nil {
		return errors.New("调查提督所在镇守府失败")
	}
	defer res.Body.Close()
	br := bufio.NewReader(res.Body)
	bf := bytes.NewBuffer(nil)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		bf.Write(line)
	}
	svdata := new(worldInfoRes)
	err = json.Unmarshal(bf.Bytes()[7:], svdata)
	if err != nil || svdata.APIResult != 1 {
		return errors.New("调查提督所在镇守府时发生错误")
	}
	self.worldID = svdata.APIData.APIWorldId
	self.WorldIP = worldIPList[self.worldID-1]
	return nil
}

type apiTokenRes struct {
	APIResult    int    `json:"api_result"`
	APIResultMsg string `json:"api_result_msg"`
	APIStarttime int64  `json:"api_starttime"`
	APIToken     string `json:"api_token"`
}

func (self *KancolleAuth) getAPIToken() error {
	u := fmt.Sprintf(urlLayouts["get_entry"], self.WorldIP, self.owner, time.Now().UnixNano()/1e6)
	data := make(url.Values)
	data.Set("url", u)
	data.Set("httpMethod", "GET")
	data.Set("authz", "signed")
	data.Set("st", self.st)
	data.Set("contentType", "JSON")
	data.Set("numEntries", "3")
	data.Set("getSummaries", "false")
	data.Set("signOwner", "true")
	data.Set("signViewer", "true")
	data.Set("gadget", "http://203.104.209.7/gadget.xml")
	data.Set("container", "dmm")
	req, _ := http.NewRequest(http.MethodPost, urlLayouts["make_request"], strings.NewReader(data.Encode()))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := self.Do(req)
	if err != nil {
		return errors.New("调查提督进入镇守府的口令失败")
	}
	defer res.Body.Close()
	br := bufio.NewReader(res.Body)
	bf := bytes.NewBuffer(nil)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		bf.Write(line)
	}
	resData := new(map[string]struct {
		RC      int               `json:"rc"`
		Body    string            `json:"body"`
		Headers map[string]string `json:"headers"`
	})
	err = json.Unmarshal(bf.Bytes()[27:], resData)
	rs := (*resData)[u]
	if err != nil || rs.RC != 200 {
		return errors.New("调查提督进入镇守府的口令失败")
	}
	svdata := new(apiTokenRes)
	err = json.NewDecoder(strings.NewReader(rs.Body[7:])).Decode(svdata)
	if err != nil || svdata.APIResult != 1 {
		return errors.New("调查提督进入镇守府的口令时发生错误")
	}
	self.APIToken = svdata.APIToken
	self.APIStartTime = svdata.APIStarttime
	self.entry = fmt.Sprintf(urlLayouts["entry"], self.WorldIP, self.APIToken, self.APIStartTime)
	return nil
}

func (self *KancolleAuth) GetOSAPI() (error, string) {
	var err error
	err = self.getDMMTokens()
	if err != nil {
		return err, ""
	}
	err = self.getAjaxToken()
	if err != nil {
		return err, ""
	}
	err = self.getOSAPIURL()
	return err, self.osapiURL
}

func (self *KancolleAuth) GetEntry() (error, string) {
	var err error
	err, _ = self.GetOSAPI()
	if err != nil {
		return err, ""
	}
	err = self.getWorld()
	if err != nil {
		return err, ""
	}
	err = self.getAPIToken()
	return err, self.entry
}
