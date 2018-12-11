package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/MZIchenjl/ooi4/session"
	"github.com/gorilla/mux"
)

type APIHandler struct {
	baseHandler
	worlds map[string][]byte
	mu     sync.Mutex
}

func (self *APIHandler) WorldImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	size := vars["size"]
	sess := session.GetSession(r, self.CookieID, self.Secret)
	if sess.WorldIP != "" {
		ipSections := strings.Split(sess.WorldIP, ".")
		for i, v := range ipSections {
			ipSections[i] = fmt.Sprintf("%03s", v)
		}
		imageName := fmt.Sprintf("%s_%s", strings.Join(ipSections, "_"), size)
		if self.worlds[imageName] == nil {
			u := fmt.Sprintf("http://203.104.209.102/kcs/resources/image/world/%s.png", imageName)
			coro, err := http.Get(u)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(http.StatusText(http.StatusBadRequest)))
				return
			}
			defer coro.Body.Close()
			self.mu.Lock()
			self.worlds[imageName], err = ioutil.ReadAll(coro.Body)
			self.mu.Unlock()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(http.StatusText(http.StatusBadRequest)))
				return
			}
		}
		w.Write(self.worlds[imageName])
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}
}

func (self *APIHandler) API(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	sess := session.GetSession(r, self.CookieID, self.Secret)
	if sess.WorldIP != "" {
		referer := r.Referer()
		referer = strings.Replace(referer, r.Host, sess.WorldIP, 1)
		referer = strings.Replace(referer, "https://", "http://", 1)
		u, err := url.Parse(fmt.Sprintf("http://%s/kcsapi/%s", sess.WorldIP, action))
		req, err := http.NewRequest(r.Method, u.String(), r.Body)
		req.Header = r.Header
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko")
		req.Header.Set("Origin", strings.Replace(r.Header.Get("Origin"), r.Host, sess.WorldIP, 1))
		req.Header.Set("Referer", referer)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}
		defer res.Body.Close()
		w.Header().Set("Date", res.Header.Get("Date"))
		w.Header().Set("Connection", res.Header.Get("Connection"))
		w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
		w.Header().Set("Content-Encoding", res.Header.Get("Content-Encoding"))
		buf := make([]byte, chunkSize)
		for {
			n, err := res.Body.Read(buf)
			if err != nil && err != io.EOF {
				return
			}
			w.Write(buf[:n])
			if err == io.EOF {
				break
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}
}
