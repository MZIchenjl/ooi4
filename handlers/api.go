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
	BaseHandler
	worlds map[string][]byte
	mu     sync.Mutex
}

func (self *APIHandler) WorldImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	size := vars["size"]
	sess := session.GetSession(r, self.CookieID(), self.Secret())
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
				return
			}
			defer coro.Body.Close()
			self.mu.Lock()
			self.worlds[imageName], err = ioutil.ReadAll(coro.Body)
			self.mu.Unlock()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		w.Write(self.worlds[imageName])
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (self *APIHandler) API(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	sess := session.GetSession(r, self.CookieID(), self.Secret())
	if sess.WorldIP != "" {
		referer := r.Referer()
		referer = strings.Replace(referer, r.Host, sess.WorldIP, 1)
		referer = strings.Replace(referer, "https://", "http://", 1)
		u, err := url.Parse(fmt.Sprintf("http://%s/kcsapi/%s", sess.WorldIP, action))
		r.URL = u
		r.Host = u.Host
		r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko")
		r.Header.Set("Origin", fmt.Sprintf("http://%s/", sess.WorldIP))
		r.Header.Set("Referer", referer)
		res, err := http.DefaultClient.Do(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		defer res.Body.Close()
		w.Header().Set("Content-Type", "text/plain")
		buf := make([]byte, chunkSize)
		for {
			n, err := res.Body.Read(buf)
			if err != nil && err != io.EOF {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if 0 == n {
				break
			}
			w.Write(buf)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
