package httpd

import (
	"net/http"
	"text/template"
)

var (
	index = template.Must(template.ParseFiles("httpd/templates/index.html"))
)

func Index(w http.ResponseWriter, r *http.Request) {
	index.Execute(w, nil)
}
