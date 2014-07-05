package httpd

import (
	"net/http"
	"text/template"
)

type ViewHandle struct {
	templates *template.Template
}

func NewViewHandle() *ViewHandle {
	t := template.Must(template.ParseGlob("httpd/templates/*"))

	return &ViewHandle{t}
}

func (v *ViewHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.templates.ExecuteTemplate(w, "index", nil)
}
