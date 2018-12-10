package handlers

type baseHandler struct {
	Secret   string
	CookieID string
}

const chunkSize = 1024
