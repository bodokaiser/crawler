package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/bodokaiser/gerenuk/httpd"
	"github.com/bodokaiser/gerenuk/parser/html"
)

func main() {
	req, _ := http.NewRequest("GET", "http://www.google.com", nil)

	p := httpd.NewPool()

	p.Add(req)
	p.Run()

	for {
		req, res, err := p.Get()

		if err != nil {
			log.Fatal(err)
		}

		href := html.NewHrefParser(res.Body)

		for {
			r := href.Next()

			if r == nil {
				break
			}

			fmt.Printf("%s href: %s\n", req.Host, r.String())

			if bytes.HasPrefix(r.Value, []byte("http")) {
				req, _ := http.NewRequest("GET", r.String(), nil)

				p.Add(req)
			}
		}

		res.Body.Close()
	}
}
