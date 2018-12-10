package templates

import "html/template"

var Poi = template.Must(template.New("Poi").Parse(`<!DOCTYPE html>
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
  <style type="text/css">
    html, body {
      overflow: hidden;
    }
    #externalswf {
      width: 800px;
      height: 480px;
    }
  </style>
</head>
<body>
<div id="flashWrap">
  <embed id="externalswf" width="1280" height="800" wmode="opaque" quality="high" bgcolor="#000000" allowscriptaccess="always" base="{{.Scheme}}://{{.Host}}/kcs2/"  src="{{.Scheme}}://{{.Host}}/kcs2/index.php?api_root=/kcsapi&voice_root=/kcs/sound&osapi_root=osapi.dmm.com&version=4.0.0.0&api_token={{.Token}}&api_starttime={{.StartTime}}" style="display: block !important;"></embed>
</div>
<script>
  $(function() {
    var swf = $("#externalswf");
    swf.css("width", "100%");
    swf.css("height", "100%");
    swf.css("position", "absolute");
    swf.css("margin", "0");
  });
</script>
<div class="statistics">
</div>
</body>
</html>
`))
