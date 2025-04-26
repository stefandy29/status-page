package main

const style_html = `
<style>
* {
    margin: 0px;
    padding: 0px;
    box-sizing: border-box;
    background-color:ghostwhite;
}

pre {
    white-space: pre-wrap;       /* Since CSS 2.1 */
    white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
    white-space: -pre-wrap;      /* Opera 4-6 */
    white-space: -o-pre-wrap;    /* Opera 7 */
    word-wrap: break-word;       /* Internet Explorer 5.5+ */
}

body {
    font-family: 'Calibri';
    font-weight: 100;
}

.presentation{
    margin-left: auto;
    margin-right: auto;
    align-items: center;
    width: 40%;
    padding-top: 125px;
    font-size: 20px;
}

footer{
    text-align: center;
}

@media screen and (max-width: 1280px)  {
    .presentation{
        width: 60%;
    }
}


@media screen and (max-width: 1024px) {
    .presentation{
        width: 70%;
    }
}

@media screen and (max-width: 720px) {
    .presentation{
        width: 95%;
    }
}

@media screen and (max-width: 480px) {
    .presentation{
        width: 100%;
    }
}

.flex-title{
    padding-bottom: 15px;
}

.flex-container {
  display: flex;
  flex-wrap: wrap;
  font-size: 30px;
  border-style: solid;
  border-width: 0px;
  border-bottom-width: 0.1px;
  padding-bottom: 5px;
  margin-bottom: 30px;
}

.flex-container:hover{
    opacity: 60%;
}

.flex-item-left {
  text-align: left;
  padding : 10px;
  flex: 50%;
  margin: auto;
  word-wrap: break-word;
  width: 70px;
}


.flex-item-right {
  text-align: right;
  padding: 10px;
  flex: 50%;
  margin: auto;
  word-wrap: break-word;
  width: 30px;
}
</style>
`

const skeleton = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="refresh" content="{{.Page.Reload}}">
    <title>{{.Page.Title}}</title>
    {{.Page.Style}}
</head>



<main>
    <section class="presentation">
        <div class="introduction">
        <h1>Updated at {{.Page.Date}}</h1>
        <hr>
			{{.Content}}
        </div>
    </section>
</main>
<footer>
    <div>Copyright © {{.Page.Year}}</div>
</footer>
</html>
`

const server_list_html = `<h1 class="flex-title">{{.ServerName}}</h1>{{range .ListMetrics}}{{.}}{{end}}`

const metric_html = `<div class="flex-container" style="background-color: rgba(255, 44, 44,  {{ Divide .Value}});">
	<div class="flex-item-left" style="background-color: rgba(255, 44, 44, 0);">{{.Key}}</div>
	<div class="flex-item-right" style="background-color: rgba(255, 44, 44, 0);">{{.Value}} {{.Type_Size}}</div>
</div>`
