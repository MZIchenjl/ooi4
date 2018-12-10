package session

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func getSalt(l int) string {
	var j []byte
	for i := 0; i < l; i++ {
		r := rand.Intn(26)
		j = append(j, alphabet[r])
	}
	return strings.Join(strings.Split(string(j), ""), ".")
}
