package httpd

import (
	"net/http"
	"text/template"
)

type Index struct {
	template *template.Template
}

func NewIndex() *Index {
	t := template.Must(template.ParseFiles("httpd/templates/index.html"))

	return &Index{t}
}

func (i *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.template.Execute(w, nil)
}
