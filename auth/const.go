package auth

import "regexp"

var worldIPList = []string{
	"203.104.209.71",
	"203.104.209.87",
	"125.6.184.215",
	"203.104.209.183",
	"203.104.209.150",
	"203.104.209.134",
	"203.104.209.167",
	"203.104.248.135",
	"125.6.189.7",
	"125.6.189.39",
	"125.6.189.71",
	"125.6.189.103",
	"125.6.189.135",
	"125.6.189.167",
	"125.6.189.215",
	"125.6.189.247",
	"203.104.209.23",
	"203.104.209.39",
	"203.104.209.55",
	"203.104.209.102",
}

var urlLayouts = map[string]string{
	"login":        "https://accounts.dmm.com/service/login/password/=/",
	"ajax":         "https://accounts.dmm.com/service/api/get-token/",
	"auth":         "https://accounts.dmm.com/service/login/password/authenticate/",
	"game":         "http://www.dmm.com/netgame/social/-/gadgets/=/app_id=854854/",
	"make_request": "http://osapi.dmm.com/gadgets/makeRequest",
	"get_world":    "http://203.104.209.7/kcsapi/api_world/get_id/%s/1/%d",
	"get_entry":    "http://%s/kcsapi/api_auth_member/dmmlogin/%s/1/%d",
	"entry":        "http://%s/kcs2/index.php?api_root=/kcsapi&voice_root=/kcs/sound&osapi_root=osapi.dmm.com&version=4.0.0.0&api_token=%s&api_starttime=%d",
}

const userAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko"

var patterns = map[string]*regexp.Regexp{
	"dmm_token": regexp.MustCompile(`http-dmm-token" content="([\d|\w]+)"`),
	"token":     regexp.MustCompile(`token" content="([\d|\w]+)"`),
	"reset":     regexp.MustCompile(`認証エラー`),
	"osapi":     regexp.MustCompile(`URL\W+:\W+"(.*)",`),
}
