package httpd

import (
	"net/http"

	"github.com/bodokaiser/go-crawler/conf"
)

type Httpd struct{}

func New() *Httpd {
	return &Httpd{}
}

func (h *Httpd) Listen(c conf.Conf) {
	http.HandleFunc("/", Index)

	http.ListenAndServe(c["addr"], nil)
}
