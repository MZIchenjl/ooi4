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
	"time"
)

type KancolleAuth struct {
	http.Client
	header       http.Header
	worldID      int
	loginID      string
	password     string
	dmmToken     string
	token        string
	idKey        string
	owner        string
	osapiURL     string
	pwdKey       string
	worldIP      string
	apiToken     string
	apiStartTime string
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
		fmt.Println(err, j)
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
