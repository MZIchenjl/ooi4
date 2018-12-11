package handlers

import "github.com/gorilla/sessions"

type baseHandler struct {
	cookieName  string
	cookieStore *sessions.CookieStore
}

func (self *baseHandler) Init(secret, cookie string) {
	self.cookieName = cookie
	self.cookieStore = sessions.NewCookieStore([]byte(secret))
	self.cookieStore.Options.MaxAge = 0
}

const chunkSize = 1024
