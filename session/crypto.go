package session

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
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
	mac := val[spIdx+1:]
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

func PKCS5Padding(val []byte, blockSize int) []byte {
	padding := blockSize - len(val)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(val, padtext...)
}

func PKCS5UnPadding(val []byte) []byte {
	length := len(val)
	if length == 0 {
		return val
	}
	unpadding := int(val[length-1])
	return val[:(length - unpadding)]
}

func decrypt(val []byte, secret string) []byte {
	block, _ := aes.NewCipher(padding(secret))
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, []byte(secret[:blockSize]))
	ret := make([]byte, len(val))
	blockMode.CryptBlocks(ret, val)
	ret = PKCS5UnPadding(ret)
	return ret
}

func encrypt(val []byte, secret string) []byte {
	block, _ := aes.NewCipher(padding(secret))
	blockSize := block.BlockSize()
	val = PKCS5Padding(val, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte(secret[:blockSize]))
	ret := make([]byte, len(val))
	blockMode.CryptBlocks(ret, val)
	return ret
}
