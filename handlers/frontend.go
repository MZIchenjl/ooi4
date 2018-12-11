package handlers

import (
	"net/http"
	"strconv"

	"github.com/MZIchenjl/ooi4/auth"
	"github.com/MZIchenjl/ooi4/session"
	"github.com/MZIchenjl/ooi4/templates"
)

type FrontEndHandler struct {
	baseHandler
}

type TmplParams struct {
	Mode      int
	StartTime int64
	OSAPIURL  string
	Host      string
	ErrMsg    string
	Token     string
}

func setCookie(w http.ResponseWriter, sess *session.Session, name, secret string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    session.GetCooke(sess, secret),
		HttpOnly: true,
		MaxAge:   0,
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
	sess := session.GetSession(r, self.CookieID, self.Secret)
	if sess.Mode != 0 {
		mode = sess.Mode
	} else {
		mode = 1
		sess.Mode = mode
		setCookie(w, sess, self.CookieID, self.Secret)
	}
	templates.Form.Execute(w, TmplParams{Mode: mode})
}

func (self *FrontEndHandler) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	sess := session.GetSession(r, self.CookieID, self.Secret)
	loginID := r.Form.Get("login_id")
	password := r.Form.Get("password")
	mode, err := strconv.ParseInt(r.Form.Get("mode"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	sess.Mode = int(mode)
	if loginID != "" && password != "" {
		kancolle := auth.New(loginID, password)
		switch mode {
		case 1, 2, 3:
			err, _ := kancolle.GetEntry()
			if err != nil {
				templates.Form.Execute(w, TmplParams{
					ErrMsg: err.Error(),
					Mode:   sess.Mode,
				})
				return
			}
			sess.APIStartTime = kancolle.APIStartTime
			sess.APIToken = kancolle.APIToken
			sess.WorldIP = kancolle.WorldIP
			setCookie(w, sess, self.CookieID, self.Secret)
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
				templates.Form.Execute(w, TmplParams{
					ErrMsg: err.Error(),
					Mode:   sess.Mode,
				})
				return
			}
			sess.OSAPIURL = osapiURL
			setCookie(w, sess, self.CookieID, self.Secret)
			http.Redirect(w, r, "/connector", http.StatusFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		}
	} else {
		templates.Form.Execute(w, TmplParams{
			ErrMsg: "请输入完整的登录ID和密码",
			Mode:   sess.Mode,
		})
	}
}

func (self *FrontEndHandler) Normal(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r, self.CookieID, self.Secret)
	if sess.APIStartTime != 0 && sess.APIToken != "" && sess.WorldIP != "" {
		templates.Normal.Execute(w, TmplParams{
			Host:      r.Host,
			Token:     sess.APIToken,
			StartTime: sess.APIStartTime,
		})
		return
	} else {
		clearCookie(w, self.CookieID)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (self *FrontEndHandler) Flash(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r, self.CookieID, self.Secret)
	if sess.APIStartTime != 0 && sess.APIToken != "" && sess.WorldIP != "" {
		templates.Flash.Execute(w, TmplParams{
			Host:      r.Host,
			Token:     sess.APIToken,
			StartTime: sess.APIStartTime,
		})
		return
	} else {
		clearCookie(w, self.CookieID)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (self *FrontEndHandler) KCV(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r, self.CookieID, self.Secret)
	if sess.APIStartTime != 0 && sess.APIToken != "" && sess.WorldIP != "" {
		templates.KCV.Execute(w, nil)
		return
	} else {
		clearCookie(w, self.CookieID)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (self *FrontEndHandler) Poi(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r, self.CookieID, self.Secret)
	if sess.APIStartTime != 0 && sess.APIToken != "" && sess.WorldIP != "" {
		templates.Poi.Execute(w, TmplParams{
			Host:      r.Host,
			Token:     sess.APIToken,
			StartTime: sess.APIStartTime,
		})
		return
	} else {
		clearCookie(w, self.CookieID)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (self *FrontEndHandler) Connector(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r, self.CookieID, self.Secret)
	if sess.OSAPIURL != "" {
		templates.Connector.Execute(w, TmplParams{
			OSAPIURL: sess.OSAPIURL,
		})
		return
	} else {
		clearCookie(w, self.CookieID)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (self *FrontEndHandler) Logout(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, self.CookieID)
	http.Redirect(w, r, "/", http.StatusFound)
}
