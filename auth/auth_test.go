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

func Test_getOSAPIURL(t *testing.T) {
	auth.getDMMTokens()
	auth.getAjaxToken()
	err := auth.getOSAPIURL()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("OK")
	}
}

func Test_getWorld(t *testing.T) {
	auth.getDMMTokens()
	auth.getAjaxToken()
	auth.getOSAPIURL()
	err := auth.getWorld()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("OK")
	}
}

func Test_getAPIToken(t *testing.T) {
	auth.getDMMTokens()
	auth.getAjaxToken()
	auth.getOSAPIURL()
	auth.getWorld()
	err := auth.getAPIToken()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("OK")
	}
}
