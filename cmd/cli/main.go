package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bodokaiser/gerenuk/httpd"
	"github.com/bodokaiser/gerenuk/parser/html"
	"github.com/bodokaiser/gerenuk/store"
)

var url string

var pool = httpd.NewPool()

var list = store.NewList()

func main() {
	flag.StringVar(&url, "url", "", "URL to crawl.")
	flag.Parse()

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	list.Add(req.URL.String())
	pool.Add(req)
	pool.Run()

	for {
		req, res, err := pool.Get()

		if err != nil {
			log.Fatal()
		}

		s := bufio.NewScanner(res.Body)
		s.Split(html.ScanHref)

		for s.Scan() {
			t := s.Text()

			if strings.HasPrefix(t, "/") {
				req.URL.Path = t
				t = req.URL.String()
			}
			if strings.HasPrefix(t, "http") && !list.Has(t) {
				req, _ := http.NewRequest("GET", t, nil)

				list.Add(t)
				pool.Add(req)

				fmt.Printf("Found url: %s\n", t)
			}
		}

		res.Body.Close()
	}
}
