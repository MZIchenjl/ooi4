package template

import "html/template"

var Normal = template.Must(template.New("Normal").Parse(`<!DOCTYPE html>
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
<div id="ooi-page" class="uk-container uk-container-center">
  <div id="ooi-header" class="uk-grid uk-grid-small">
    <div id="ooi-logo" class="uk-width-small-1-10">
      <img src="/static/img/logo.png">
    </div>
    <div id="ooi-headline" class="uk-width-small-9-10">
      <h1 class="uk-text-primary">OOI - 舰娘在线缓存系统</h1>
      <hr>
    </div>
  </div>
  <div id="ooi-content" class="uk-grid">
    <div id="ooi-game" class="uk-width-1-1 uk-text-center">
      <embed id="externalswf" width="1280" height="800" wmode="opaque" quality="high" bgcolor="#000000" allowscriptaccess="always" base="{{.Scheme}}://{{.Host}}/kcs2/" src="{{.Scheme}}://{{.Host}}/kcs2/index.php?api_root=/kcsapi&voice_root=/kcs/sound&osapi_root=osapi.dmm.com&version=4.0.0.0&api_token={{.Token}}&amp;api_starttime={{.StartTime}}" style="display: block !important;"></embed>
      <a href="/logout" class="uk-button uk-button-primary">重新登录</a>
    </div>
  </div>
  <div id="ooi-footer" class="uk-text-center">
    <hr>
    <address><a href="http://kancolle.tv" target="_blank">海平线镇守府</a> &copy; 2014-2016</address>
  </div>
</div>
<div class="statistics">
</div>
</body>
</html>
`))
