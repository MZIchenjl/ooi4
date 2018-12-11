package templates

import "html/template"

var Flash = template.Must(template.New("Flash").Parse(`<!DOCTYPE html>
<html lang="zh-Hans-CN">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="keywords" content="舰队collection,舰娘,艦隊これくしょん,艦これ">
  <title>OOI - 舰娘在线缓存系统</title>
  <link href="/static/css/uikit.min.css" rel="stylesheet">
  <link href="/static/css/uikit.almost-flat.min.css" rel="stylesheet">
  <link href="/static/css/ooi.css" rel="stylesheet">
  <script src="/static/js/jquery-2.1.4.min.js"></script>
  <script src="/static/js/uikit.min.js"></script>
</head>
<body>
  <div id="spacing_top" style="height:16px;"></div>
  <div id="flashWrap" style="width: 800px; margin: auto;">
    <embed id="externalswf" width="1280" height="800" wmode="opaque" quality="high" bgcolor="#000000" allowscriptaccess="always" base="{{.Schema}}://{{.Host}}/kcs2/" src="{{.Schema}}://{{.Host}}/kcs2/index.php?api_root=/kcsapi&voice_root=/kcs/sound&osapi_root=osapi.dmm.com&version=4.0.0.0&api_token={{.Token}}&amp;api_starttime={{.StartTime}}" style="display: block !important;"></embed>
  </div>
  <div class="statistics">
  </div>
</body>
</html>
`))
