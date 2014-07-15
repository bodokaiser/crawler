package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"

	ghttp "github.com/bodokaiser/gerenuk/net/http"
	gurl "github.com/bodokaiser/gerenuk/net/url"
	"github.com/bodokaiser/gerenuk/text/html"
)

var (
	fs = http.FileServer(http.Dir("srv"))

	ev = ghttp.NewEventHandler(connect)
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

		list := gurl.NewList()
		pool := ghttp.NewPool()
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

						ghttp.SendEvent(w, t[i:])
					}
				}

				res.Body.Close()
			}
		}

	}
}
