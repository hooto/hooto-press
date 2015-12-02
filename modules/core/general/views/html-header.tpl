<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>{{if .__html_head_title__}}{{.__html_head_title__}} | {{end}}{{SysConfig "frontend_html_head_subtitle"}}</title>
  <link rel="stylesheet" href="{{HttpSrvBasePath "~/bootstrap/3.3/css/bootstrap.min.css"}}" type="text/css">
  <link rel="stylesheet" href="{{HttpSrvBasePath "~/l5s/css/main.css"}}" type="text/css">
  <link rel="shortcut icon" type="image/x-icon" href="{{HttpSrvBasePath "~/l5s/img/ap.ico"}}">
  <meta name="keywords" content="{{SysConfig "frontend_html_head_meta_keywords"}}">
  <meta name="description" content="{{SysConfig "frontend_html_head_meta_description"}}">
  <script src="{{HttpSrvBasePath "/~/lessui/js/sea.js"}}"></script>
  <script src="{{HttpSrvBasePath "/~/l5s/js/main.js"}}"></script>
  <script type="text/javascript">
    window._basepath = {{HttpSrvBasePath ""}};
    window.onload_hooks = [];
  </script>
</head>
<body>

