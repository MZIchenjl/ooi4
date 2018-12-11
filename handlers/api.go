package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/MZIchenjl/ooi4/auth"
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
	session, err := self.cookieStore.Get(r, self.cookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	worldIP := session.Values["world_ip"].(string)
	if worldIP != "" {
		ipSections := strings.Split(worldIP, ".")
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
	session, err := self.cookieStore.Get(r, self.cookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	worldIP := session.Values["WorldIP"].(string)
	if worldIP != "" {
		referer := r.Referer()
		referer = strings.Replace(referer, r.Host, worldIP, 1)
		referer = strings.Replace(referer, "https://", "http://", 1)
		u, err := url.Parse(fmt.Sprintf("http://%s/kcsapi/%s", worldIP, action))
		req, err := http.NewRequest(r.Method, u.String(), r.Body)
		req.Header = r.Header
		req.Header.Set("User-Agent", auth.UserAgent)
		req.Header.Set("Origin", strings.Replace(r.Header.Get("Origin"), r.Host, worldIP, 1))
		req.Header.Set("Referer", referer)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}
		defer res.Body.Close()
		for key := range res.Header {
			w.Header().Set(key, res.Header.Get(key))
		}
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
