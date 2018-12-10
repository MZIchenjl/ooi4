package session

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func getSalt(l int) string {
	var j []byte
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < l; i++ {
		r := rand.Intn(26)
		j = append(j, alphabet[r])
	}
	return strings.Join(strings.Split(string(j), ""), ".")
}
