package session

import (
	"fmt"
	"testing"
)

func Test_getSalt(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		getSalt(16)
	}
	fmt.Println(getSalt(16))
	t.Log("OK")
}
