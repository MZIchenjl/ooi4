package handlers

import (
	"fmt"
	"github.com/MZIchenjl/ooi4/session"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type APIHandler struct {
	BaseHandler
	apiStart2 []byte
	worlds    map[string][]byte
	mu        sync.Mutex
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
