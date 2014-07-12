package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/bodokaiser/gerenuk/httpd"
	"github.com/bodokaiser/gerenuk/parser/html"
)

func main() {
	request("http://www.satisfeet.me")

	time.Sleep(10 * time.Second)
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

		/*
			if bytes.IndexRune(r.Value, '/') == 0 {
				go request("http://" + cr.Host + r.String())
			}
		*/
		if bytes.HasPrefix(r.Value, []byte("http")) {
			go request(r.String())
		}
	}
}
