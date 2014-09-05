package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bodokaiser/gerenuk/httpd"
	"github.com/bodokaiser/gerenuk/parser/html"
	"github.com/bodokaiser/gerenuk/store"
)

var host string

// File Handler.
var files = http.FileServer(http.Dir("share"))

func main() {
	flag.StringVar(&host, "host", ":3000", "The host to listen on.")
	flag.Parse()

	http.Handle("/", files)
	http.Handle("/events", http.HandlerFunc(handle))
	http.ListenAndServe(host, nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	event, err := httpd.NewEventStream(r, w)
	if err != nil {
		log.Fatal(err)
	}

	ref, err := url.Parse(r.Header.Get("Referer"))
	if err != nil {
		log.Fatal(err)
	}

	url, err := url.QueryUnescape(ref.Query().Get("url"))
	if err != nil {
		log.Fatal(err)
	}

	if req, err := http.NewRequest("GET", url, nil); err == nil {
		log.Printf("Starting to crawl: %s\n", url)

		list := store.NewList()
		list.Add(url)

		pool := httpd.NewPool()
		pool.Add(req)
		pool.Run()

		for {
			_, res, err := pool.Get()
			if err != nil {
				log.Fatal(err)
			}

			if res != nil {
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
					}
					if strings.HasPrefix(t, "mailto:") {
						i := strings.IndexRune(t, ':') + 1

						event.Emit(t[i:])
					}
				}

				res.Body.Close()
			}
		}

	}
}
