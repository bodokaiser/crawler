package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bodokaiser/gerenuk/httpd"
	"github.com/bodokaiser/gerenuk/utils"
)

var (
	fs = http.FileServer(http.Dir("httpd/public"))

	ev = httpd.NewEventHandler(connect)
)

func main() {
	host := flag.String("host", ":3000", "The host to listen on.")

	flag.Parse()

	http.Handle("/", fs)
	http.Handle("/events", ev)

	http.ListenAndServe(*host, nil)
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

		p := httpd.NewPool()
		p.Add(req)
		p.Run()

		for {
			_, res, err := p.Get()

			if err != nil {
				log.Fatal(err)

				return
			}

			if res != nil {
				s := bufio.NewScanner(res.Body)
				s.Split(utils.ScanHref)

				for s.Scan() {
					httpd.SendEvent(w, "message", s.Text())

					if strings.HasPrefix(s.Text(), "http") {
						req, _ := http.NewRequest("GET", s.Text(), nil)

						p.Add(req)
					}

					f.Flush()
				}

				res.Body.Close()
			}
		}

	}
}
