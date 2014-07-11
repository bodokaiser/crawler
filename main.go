package main

import (
	"fmt"
	"log"

	"github.com/bodokaiser/gerenuk/httpd"
	"github.com/bodokaiser/gerenuk/parser/html"
)

func main() {
	request("http://www.satisfeet.me")
}

func request(u string) {
	c := httpd.NewClient(u)

	c.Handle(parse)
	c.Handle(follow)

	if err := c.Open(); err != nil {
		log.Fatal(err)
	}
}

func parse(cr *httpd.ClientResult) {
	p := html.NewHrefParser(cr.Body)

	for {
		r := p.Next()

		if r == nil {
			break
		}

		fmt.Printf("%s href: %s\n", cr.Host, r.String())
	}
}

func follow(cr *httpd.ClientResult) {
	p := html.NewHrefParser(cr.Body)

	for {
		r := p.Next()

		if r == nil {
			break
		}
	}
}
