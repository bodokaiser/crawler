package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	ghttp "github.com/bodokaiser/gerenuk/net/http"
	"github.com/bodokaiser/gerenuk/net/url"
	"github.com/bodokaiser/gerenuk/text/html"
)

var (
	pool = ghttp.NewPool()
	list = url.NewList()
)

func main() {
	url := flag.String("url", "", "URL to crawl.")
	flag.Parse()

	req, err := http.NewRequest("GET", *url, nil)

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
