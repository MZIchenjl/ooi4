package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/MZIchenjl/ooi4/auth"
)

type ProxyHandler struct {
	baseHandler
}

func (self *ProxyHandler) Proxy(w http.ResponseWriter, r *http.Request) {
	host := auth.WorldIPList[0]
	u, _ := url.Parse(r.URL.String())
	u.Scheme = "http"
	u.Host = host
	req, _ := http.NewRequest(r.Method, u.String(), r.Body)
	referer := r.Referer()
	referer = strings.Replace(referer, r.Host, host, 1)
	referer = strings.Replace(referer, "https://", "http://", 1)
	req.Header = r.Header
	req.Header.Set("User-Agent", auth.UserAgent)
	req.Header.Set("Origin", strings.Replace(r.Header.Get("Origin"), r.Host, host, 1))
	req.Header.Set("Referer", referer)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	defer res.Body.Close()
	for key := range res.Header {
		if key != "Content-Length" {
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
}
