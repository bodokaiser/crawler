package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	. "github.com/bodokaiser/gerenuk"
	"github.com/bodokaiser/gerenuk/httpd"
)

var host string

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
		log.Fatalf("Error establishing event stream: %s.\n", err)

		return
	}

	c, err := NewCrawlerFromRequest(r)
	if err != nil {
		log.Fatalf("Error creating crawler from request: %s.\n", err)

		return
	}

	for {
		url, err := c.Get()
		if err != nil {
			log.Fatalf("Error resolving pool result: %s.\n", err)

			return
		}

		if err := event.Emit(url); err != nil {
			log.Fatal("Error emitting event: %s.\n", err)

			return
		}
	}
}

func NewCrawlerFromRequest(r *http.Request) (*Crawler, error) {
	ref, err := url.Parse(r.Header.Get("Referer"))
	if err != nil {
		return nil, err
	}

	uri, err := url.QueryUnescape(ref.Query().Get("url"))
	if err != nil {
		return nil, err
	}

	return NewCrawler(uri), nil
}
