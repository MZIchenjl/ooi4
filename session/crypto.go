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

var COOKIE_INVALID = errors.New("Cookie 签名不合法")

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
	str := []byte(val[:spIdx])
	mac := []byte(val[spIdx:])
	var strBytes []byte
	var macBytes []byte
	var err error
	_, err = base64.StdEncoding.Decode(strBytes, str)
	_, err = base64.StdEncoding.Decode(macBytes, mac)
	if err != nil {
		return COOKIE_INVALID, nil
	}
	nmacBytes := sign(strBytes, secret)
	if hmac.Equal(macBytes, nmacBytes) {
		return nil, strBytes
	}
	return COOKIE_INVALID, nil
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
