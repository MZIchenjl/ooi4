package templates

import "html/template"

var Form = template.Must(template.New("Form").Parse(`<html lang="zh-Hans-CN">
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
    <div id="ooi-form" class="uk-width-medium-2-5">
      {{if .ErrMsg ne ""}}
      <div class="uk-panel uk-panel-box uk-panel-box-primary uk-margin-bottom">
        <p class="uk-text-danger">{{.ErrMsg}}</p>
      </div>
      {{end}}
      <form  method="post" class="uk-form uk-form-stacked">
        <div class="uk-form-row">
          <label class="uk-form-label" for="login_id">DMM登录ID：</label>
          <div class="uk-form-controls">
            <input id="login_id" maxlength="80" name="login_id" type="email" class="uk-form-large uk-form-width-large">
          </div>
        </div>
        <div class="uk-form-row">
          <label class="uk-form-label" for="password">DMM登录密码：</label>
          <div class="uk-form-controls">
            <input id="password" maxlength="80" name="password" type="password" class="uk-form-large uk-form-width-large">
          </div>
        </div>
        <div class="uk-form-row">
          <span class="uk-form-label">游戏方式：</span>
          <p>
            <input type="radio" id="mode1" name="mode" value="1"{{if .Mode eq 1}} checked{{end}}>
            <label for="mode1">在浏览器中运行</label>
          </p>
          <p>
            <input type="radio" id="mode2" name="mode" value="2"{{if .Mode eq 2}} checked{{end}}>
            <label for="mode2">在KCV/74EO中运行</label>
          </p>
          <p>
            <input type="radio" id="mode3" name="mode" value="3"{{if .Mode eq 3}} checked{{end}}>
            <label for="mode3">在poi中运行</label>
          </p>
          <p>
            <input type="radio" id="mode4" name="mode" value="4"{{if .Mode eq 4}} checked{{end}}>
            <label for="mode4">登录器直连模式</label>
          </p>
        </div>
        <div class="uk-form-row">
          <button class="uk-button uk-button-primary uk-button-large">登录游戏</button>
        </div>
      </form>
    </div>
    <div id="ooi-announcement" class="uk-width-medium-3-5">
      <h2 class="uk-text-primary">提示</h2>
      <p>本站采用了OOI系统，OOI (Online Objects Integration)是针对DMM网页游戏《艦隊これくしょん -艦これ-》的在线缓存系统。</p>
      <p>本站只支持DMM账号登录，不支持Facebook和Google账号登录。本网站不会保存任何用户信息，并对本网站服务可能造成的任何后果不负责任。</p>
      <p>本站的源代码已更新到OOI3，该项目代码采用AGPLv3许可证开源，发布于 <a href="https://github.com/acgx/ooi3" target="_blank">此版本库</a>。</p>
      <p>如果希望资助我们的项目，可以通过<a href="http://bbs.kancolle.tv/thread-4812-1-1.html" target="_blank">购买赞助商产品</a>的方式支持我们。</p>
      <p>OOI用户QQ群：337440297 | <a href="https://github.com/acgx/ooi3/issues" target="_blank">汇报bug</a></p>
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
