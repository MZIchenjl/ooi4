package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/MZIchenjl/ooi4/auth"
	"github.com/gorilla/mux"
)

type APIHandler struct{}

func (self *APIHandler) WorldImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	size := vars["size"]
	session, err := cookieStore.Get(r, cookieName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	worldIP := session.Values["world_ip"]
	if worldIP != nil {
		ipSections := strings.Split(worldIP.(string), ".")
		for i, v := range ipSections {
			ipSections[i] = fmt.Sprintf("%03s", v)
		}
		imageName := fmt.Sprintf("%s_%s", strings.Join(ipSections, "_"), size)
		u := fmt.Sprintf("http://203.104.209.102/kcs2/resources/world/%s.png", imageName)
		coro, err := http.Get(u)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer coro.Body.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer coro.Body.Close()
		for key := range coro.Header {
			if !isExluded(key) {
				w.Header().Set(key, coro.Header.Get(key))
			}
		}
		buf := make([]byte, chunkSize)
		for {
			n, err := coro.Body.Read(buf)
			if err != nil && err != io.EOF {
				return
			}
			w.Write(buf[:n])
			if err == io.EOF {
				break
			}
		}
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (self *APIHandler) API(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, cookieName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	worldIP := session.Values["world_ip"]
	if worldIP != nil {
		WorldIP := worldIP.(string)
		referer := r.Referer()
		referer = strings.Replace(referer, r.Host, WorldIP, 1)
		referer = strings.Replace(referer, "https://", "http://", 1)
		u := *r.URL
		u.Scheme = "http"
		u.Host = WorldIP
		req, err := http.NewRequest(r.Method, u.String(), r.Body)
		req.Header = r.Header
		req.Header.Set("User-Agent", auth.UserAgent)
		req.Header.Set("Origin", strings.Replace(r.Header.Get("Origin"), r.Host, WorldIP, 1))
		req.Header.Set("Referer", referer)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer res.Body.Close()
		for key := range res.Header {
			if !isExluded(key) {
				w.Header().Set(key, res.Header.Get(key))
			}
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
