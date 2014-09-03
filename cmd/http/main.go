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

var fs = http.FileServer(http.Dir("share"))

var ev = httpd.NewEventHandler(connect)

func main() {
	flag.StringVar(&host, "host", ":3000", "The host to listen on.")
	flag.Parse()

	http.Handle("/", fs)
	http.Handle("/events", ev)

	http.ListenAndServe(host, nil)
}

func connect(w http.ResponseWriter, r *http.Request, f http.Flusher) {
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
		pool := httpd.NewPool()
		list.Add(url)
		pool.Add(req)
		pool.Run()

		for {
			_, res, err := pool.Get()

			if err != nil {
				log.Fatal(err)

				return
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

						httpd.SendEvent(w, t[i:])
					}
				}

				res.Body.Close()
			}
		}

	}
}
