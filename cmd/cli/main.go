package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bodokaiser/gerenuk/httpd"
	"github.com/bodokaiser/gerenuk/utils"
)

var (
	pool *httpd.Pool
)

func main() {
	req, err := request()

	if err != nil {
		log.Fatal(err)
	}

	pool = httpd.NewPool()
	pool.Add(req)
	pool.Run()

	for {
		_, res, err := pool.Get()

		if err != nil {
			log.Fatal()
		}

		scan(res.Body)

		res.Body.Close()
	}
}

func scan(r io.Reader) {
	s := bufio.NewScanner(r)
	s.Split(utils.ScanHref)

	for s.Scan() {
		t := s.Text()

		fmt.Printf("href: %s\n", t)

		if strings.HasPrefix(t, "http") {
			req, _ := http.NewRequest("GET", t, nil)

			pool.Add(req)
		}
	}
}

func request() (*http.Request, error) {
	url := flag.String("url", "", "URL to crawl.")

	flag.Parse()

	return http.NewRequest("GET", *url, nil)
}
