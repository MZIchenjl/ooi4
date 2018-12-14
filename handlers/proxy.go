package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/MZIchenjl/ooi4/auth"
)

func Proxy(w http.ResponseWriter, r *http.Request) {
	host := auth.WorldIPList[0]
	u := *r.URL
	u.Scheme = "http"
	u.Host = host
	req, _ := http.NewRequest(r.Method, u.String(), r.Body)
	referer := r.Referer()
	referer = strings.Replace(referer, r.Host, host, 1)
	referer = strings.Replace(referer, "https://", "http://", 1)
	req.Header = r.Header
	req.Header.Set("User-Agent", auth.UserAgent)
	req.Header.Set("Origin", strings.Replace(r.Header.Get("Origin"), r.Host, host, 1))

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
		if 0 == n || err == io.EOF {
			break
		}
		w.Write(buf[:n])
	}
}
