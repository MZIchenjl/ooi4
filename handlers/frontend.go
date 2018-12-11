package handlers

import (
	"net/http"
	"strconv"

	"github.com/MZIchenjl/ooi4/auth"
	"github.com/MZIchenjl/ooi4/templates"
	"github.com/gorilla/sessions"
)

type FrontEndHandler struct{}

type TmplParams struct {
	Mode      int
	StartTime int64
	OSAPIURL  string
	Host      string
	ErrMsg    string
	Token     string
}

func (self *FrontEndHandler) clearCookie(w http.ResponseWriter, r *http.Request) {
	session := sessions.NewSession(cookieStore, cookieName)
	session.Save(r, w)
}

func (self *FrontEndHandler) Form(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, cookieName)
	if err != nil {
		self.clearCookie(w, r)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	mode := session.Values["mode"]
	if mode == nil {
		mode = 1
		session.Values["mode"] = 1
		session.Save(r, w)
	}
	templates.Form.Execute(w, TmplParams{Mode: mode.(int)})
}

func (self *FrontEndHandler) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	session, err := cookieStore.Get(r, cookieName)
	if err != nil {
		self.clearCookie(w, r)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	loginID := r.Form.Get("login_id")
	password := r.Form.Get("password")
	mode, err := strconv.ParseInt(r.Form.Get("mode"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	session.Values["mode"] = int(mode)
	if loginID != "" && password != "" {
		kancolle := auth.New(loginID, password)
		switch mode {
		case 1, 2, 3:
			err, _ := kancolle.GetEntry()
			if err != nil {
				templates.Form.Execute(w, TmplParams{
					ErrMsg: err.Error(),
					Mode:   int(mode),
				})
				return
			}
			session.Values["world_ip"] = kancolle.WorldIP
			session.Values["api_token"] = kancolle.APIToken
			session.Values["api_starttime"] = kancolle.APIStartTime
			session.Save(r, w)
			switch mode {
			case 1:
				http.Redirect(w, r, "/kancolle", http.StatusFound)
			case 2:
				http.Redirect(w, r, "/kcv", http.StatusFound)
			case 3:
				http.Redirect(w, r, "/poi", http.StatusFound)
			}
		case 4:
			err, osapiURL := kancolle.GetOSAPI()
			if err != nil {
				templates.Form.Execute(w, TmplParams{
					ErrMsg: err.Error(),
					Mode:   int(mode),
				})
				return
			}
			session.Values["osapi_url"] = osapiURL
			session.Save(r, w)
			http.Redirect(w, r, "/connector", http.StatusFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		}
	} else {
		templates.Form.Execute(w, TmplParams{
			ErrMsg: "请输入完整的登录ID和密码",
			Mode:   int(mode),
		})
	}
}

func (self *FrontEndHandler) Normal(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, cookieName)
	if err != nil {
		self.clearCookie(w, r)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	apiStartTime := session.Values["api_starttime"]
	apiToken := session.Values["api_token"]
	worldIP := session.Values["world_ip"]
	if apiStartTime != nil && apiToken != nil && worldIP != nil {
		templates.Normal.Execute(w, TmplParams{
			Host:      r.Host,
			Token:     apiToken.(string),
			StartTime: apiStartTime.(int64),
		})
		return
	} else {
		self.clearCookie(w, r)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}
}

func (self *FrontEndHandler) Flash(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, cookieName)
	if err != nil {
		self.clearCookie(w, r)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	apiStartTime := session.Values["api_starttime"]
	apiToken := session.Values["api_token"]
	worldIP := session.Values["world_ip"]
	if apiStartTime != nil && apiToken != nil && worldIP != nil {
		templates.Flash.Execute(w, TmplParams{
			Host:      r.Host,
			Token:     apiToken.(string),
			StartTime: apiStartTime.(int64),
		})
		return
	} else {
		self.clearCookie(w, r)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}
}

func (self *FrontEndHandler) KCV(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, cookieName)
	if err != nil {
		self.clearCookie(w, r)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	apiStartTime := session.Values["api_starttime"]
	apiToken := session.Values["api_token"]
	worldIP := session.Values["world_ip"]
	if apiStartTime != nil && apiToken != nil && worldIP != nil {
		templates.KCV.Execute(w, nil)
		return
	} else {
		self.clearCookie(w, r)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}
}

func (self *FrontEndHandler) Poi(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, cookieName)
	if err != nil {
		self.clearCookie(w, r)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	apiStartTime := session.Values["api_starttime"]
	apiToken := session.Values["api_token"]
	worldIP := session.Values["world_ip"]
	if apiStartTime != 0 && apiToken != "" && worldIP != "" {
		templates.Poi.Execute(w, TmplParams{
			Host:      r.Host,
			Token:     apiToken.(string),
			StartTime: apiStartTime.(int64),
		})
		return
	} else {
		self.clearCookie(w, r)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}
}

func (self *FrontEndHandler) Connector(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, cookieName)
	if err != nil {
		self.clearCookie(w, r)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	osapiURL := session.Values["osapi_url"]
	if osapiURL != nil {
		templates.Connector.Execute(w, TmplParams{
			OSAPIURL: osapiURL.(string),
		})
		return
	} else {
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (self *FrontEndHandler) Logout(w http.ResponseWriter, r *http.Request) {
	self.clearCookie(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
