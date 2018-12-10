package session

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func GetSession(r *http.Request, name, secret string) *Session {
	cookie, err := r.Cookie(name)
	if err != nil {
		return NewSession()
	}
	val := cookie.Value
	if strings.Index(val, ".") == -1 {
		return NewSession()
	}
	err, rawb := unsign(val, secret)
	if err != nil {
		return NewSession()
	}
	raw := decrypt(rawb, secret)
	dec := Decode(raw)
	if dec == nil {
		return NewSession()
	}
	return dec
}

func GetCooke(s *Session, secret string) string {
	enc := encrypt(Encode(s), secret)
	mac := sign(enc, secret)
	return fmt.Sprintf("%s.%s", base64.StdEncoding.EncodeToString(enc), base64.StdEncoding.EncodeToString(mac))
}
