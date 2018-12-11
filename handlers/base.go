package handlers

type baseHandler struct {
	Secret   string
	CookieID string
}

func (self *baseHandler) Init(secret, cookie string) {
	self.Secret = secret
	self.CookieID = cookie
}

const chunkSize = 1024
