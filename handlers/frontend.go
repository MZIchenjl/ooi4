package handlers

import (
	"github.com/MZIchenjl/ooi4/auth"
	"github.com/MZIchenjl/ooi4/session"
	"github.com/MZIchenjl/ooi4/template"
	"net/http"
	"strconv"
)

type FrontEndHandler struct {
	BaseHandler
}

type TmplParams struct {
	Mode      int
	StartTime int64
	OSAPIURL  string
	Schema    string
	Host      string
	ErrMsg    string
	Token     string
}

func setCookie(w http.ResponseWriter, sess *session.Session, name, secret string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    session.GetCooke(sess, secret),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func clearCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}

func (self *FrontEndHandler) Form(w http.ResponseWriter, r *http.Request) {
	var mode int
	sess := session.GetSession(r, self.CookieID(), self.Secret())
	if sess.Mode != 0 {
		mode = sess.Mode
	} else {
		mode = 1
		sess.Mode = mode
		setCookie(w, sess, self.CookieID(), self.Secret())
	}
	template.Form.Execute(w, TmplParams{Mode: mode})
}

func (self *FrontEndHandler) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sess := session.GetSession(r, self.CookieID(), self.Secret())
	loginID := r.Form.Get("login_id")
	password := r.Form.Get("password")
	mode, err := strconv.ParseInt(r.Form.Get("mode"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sess.Mode = int(mode)
	if loginID != "" && password != "" {
		kancolle := auth.New(loginID, password)
		switch mode {
		case 1, 2, 3:
			err, _ := kancolle.GetEntry()
			if err != nil {
				template.Form.Execute(w, TmplParams{
					ErrMsg: err.Error(),
					Mode:   sess.Mode,
				})
				return
			}
			sess.APIStartTime = kancolle.APIStartTime
			sess.APIToken = kancolle.APIToken
			sess.WorldIP = kancolle.WorldIP
			setCookie(w, sess, self.CookieID(), self.Secret())
			switch mode {
			default:
				http.Redirect(w, r, "/kancolle", http.StatusFound)
			case 2:
				http.Redirect(w, r, "/kcv", http.StatusFound)
			case 3:
				http.Redirect(w, r, "/poi", http.StatusFound)
			}
		case 4:
			err, osapiURL := kancolle.GetOSAPI()
			if err != nil {
				template.Form.Execute(w, TmplParams{
					ErrMsg: err.Error(),
					Mode:   sess.Mode,
				})
				return
			}
			sess.OSAPIURL = osapiURL
			setCookie(w, sess, self.CookieID(), self.Secret())
			http.Redirect(w, r, "/connector", http.StatusFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		template.Form.Execute(w, TmplParams{
			ErrMsg: "请输入完整的登录ID和密码",
			Mode:   sess.Mode,
		})
	}
}

func (self *FrontEndHandler) Normal(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r, self.CookieID(), self.Secret())
	if sess.APIStartTime != 0 && sess.APIToken != "" && sess.WorldIP != "" {
		template.Normal.Execute(w, TmplParams{
			Schema:    r.URL.Scheme,
			Host:      r.Host,
			Token:     sess.APIToken,
			StartTime: sess.APIStartTime,
		})
		return
	} else {
		clearCookie(w, self.CookieID())
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (self *FrontEndHandler) Flash(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r, self.CookieID(), self.Secret())
	if sess.APIStartTime != 0 && sess.APIToken != "" && sess.WorldIP != "" {
		template.Flash.Execute(w, TmplParams{
			Schema:    r.URL.Scheme,
			Host:      r.Host,
			Token:     sess.APIToken,
			StartTime: sess.APIStartTime,
		})
		return
	} else {
		clearCookie(w, self.CookieID())
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (self *FrontEndHandler) Poi(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r, self.CookieID(), self.Secret())
	if sess.APIStartTime != 0 && sess.APIToken != "" && sess.WorldIP != "" {
		template.Poi.Execute(w, TmplParams{
			Schema:    r.URL.Scheme,
			Host:      r.Host,
			Token:     sess.APIToken,
			StartTime: sess.APIStartTime,
		})
		return
	} else {
		clearCookie(w, self.CookieID())
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (self *FrontEndHandler) Connector(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r, self.CookieID(), self.Secret())
	if sess.OSAPIURL != "" {
		template.Connector.Execute(w, TmplParams{
			OSAPIURL: sess.OSAPIURL,
		})
		return
	} else {
		clearCookie(w, self.CookieID())
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (self *FrontEndHandler) Logout(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, self.CookieID())
	http.Redirect(w, r, "/", http.StatusFound)
}
