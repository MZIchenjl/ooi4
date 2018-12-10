package session

import (
	"bytes"
	"crypto/aes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/pkg/errors"
	"strings"
)

var CookieInvalid = errors.New("Cookie 签名不合法")

func sign(val []byte, secret string) []byte {
	macBuilder := hmac.New(sha256.New, []byte(secret))
	macBuilder.Write(val)
	return macBuilder.Sum(nil)
}

func padding(key string) []byte {
	buf := bytes.NewBuffer([]byte(key))
	padLength := 32 - buf.Len()
	for i := 0; i < padLength; i++ {
		buf.WriteByte(0x00)
	}
	return buf.Bytes()[:32]
}

func unsign(val, secret string) (error, []byte) {
	spIdx := strings.LastIndex(val, ".")
	str := val[:spIdx]
	mac := val[spIdx:]
	var strBytes []byte
	var macBytes []byte
	var err error
	strBytes, err = base64.StdEncoding.DecodeString(str)
	macBytes, err = base64.StdEncoding.DecodeString(mac)
	if err != nil {
		return CookieInvalid, nil
	}
	nmacBytes := sign(strBytes, secret)
	if hmac.Equal(macBytes, nmacBytes) {
		return nil, strBytes
	}
	return CookieInvalid, nil
}

func decrypt(val []byte, secret string) []byte {
	cipher, _ := aes.NewCipher(padding(secret))
	var ret []byte
	cipher.Decrypt(ret, val)
	return ret
}

func encrypt(val []byte, secret string) []byte {
	cipher, _ := aes.NewCipher(padding(secret))
	var ret []byte
	cipher.Encrypt(ret, val)
	return ret
}
