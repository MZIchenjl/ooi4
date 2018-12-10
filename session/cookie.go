package session

import (
	"bytes"
	"encoding/base64"
	"net/http"
)

func GetSession(r *http.Request, name, secret string) *Session {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil
	}
	val := cookie.Value
	err, rawb := unsign(val, secret)
	raw := decrypt(rawb, secret)
	return Decode(raw)
}

func GetCooke(s *Session, secret string) string {
	sbu := bytes.NewBuffer(nil)
	enc := encrypt(Encode(s), secret)
	sbu.Write(enc)
	sbu.WriteByte('.')
	mac := sign(enc, secret)
	sbu.Write(mac)
	return base64.StdEncoding.EncodeToString(sbu.Bytes())
}
