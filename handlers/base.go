package handlers

import "os"

type BaseHandler struct {
	secret   string
	cookieID string
}

func (self *BaseHandler) Secret() string {
	if self.secret != "" {
		return self.secret
	}
	secret := os.Getenv("secret")
	self.secret = secret
	return secret
}

func (self *BaseHandler) CookieID() string {
	if self.cookieID != "" {
		return self.cookieID
	}
	cookieID := os.Getenv("cookieID")
	self.cookieID = cookieID
	return cookieID
}
