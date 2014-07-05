package httpd

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bodokaiser/go-crawler/conf"
)

type Httpd struct {
	static http.Handler
}

func New() *Httpd {
	f := http.FileServer(http.Dir("./httpd/public"))

	return &Httpd{f}
}

func (h *Httpd) Listen(c conf.Conf) {
	r := mux.NewRouter()

	r.Handle("/", NewViewHandle())
	r.Handle("/events", NewEventHandle())

	r.PathPrefix("/stylesheets").Handler(h.static)
	r.PathPrefix("/javascripts").Handler(h.static)

	http.ListenAndServe(c["addr"], r)
}
