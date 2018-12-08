package auth

import (
	"os"
	"testing"
)

var loginID = os.Getenv("loginID")
var password = os.Getenv("password")
var auth = New(loginID, password)

func Test_getDMMTokens(t *testing.T) {
	err := auth.getDMMTokens()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("OK")
	}
}

func Test_getAjaxToken(t *testing.T) {
	auth.getDMMTokens()
	err := auth.getAjaxToken()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("OK")
	}
}
