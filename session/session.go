package session

import (
	"bytes"
	"encoding/gob"
)

type Session struct {
	Mode         int
	salt         string
	WorldIP      string
	OSAPIURL     string
	APIToken     string
	APIStartTime int64
}

func NewSession() *Session {
	session := new(Session)
	session.salt = getSalt(16)
	return session
}

func Encode(s *Session) []byte {
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(s)
	return buf.Bytes()
}

func Decode(b []byte) *Session {
	session := new(Session)
	buf := bytes.NewReader(b)
	gob.NewDecoder(buf).Decode(session)
	return session
}
