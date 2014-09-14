package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"time"

	"github.com/bodokaiser/crawler"
	"github.com/bodokaiser/crawler/http/event"
)

var rate int
var path, addr string

func main() {
	flag.IntVar(&rate, "rate", 0, "")
	flag.StringVar(&path, "path", "", "")
	flag.StringVar(&addr, "addr", "", "")
	flag.Parse()

	http.Handle("/", Handler())
	http.ListenAndServe(addr, nil)
}

func Handler() http.Handler {
	c := crawler.New()

	fs := http.FileServer(http.Dir(path))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/crawl":
			m := make(map[string]string)

			if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
				panic(err)
			}

			r, err := crawler.NewRequest(m["origin"])
			if err != nil {
				panic(err)
			}

			c.Add(r)
		case "/events":
			e, err := event.NewStream(r, w)
			if err != nil {
				panic(err)
				return
			}

			for _, r := range c.Get() {
				b, err := json.Marshal(r)
				if err != nil {
					panic(err)
					return
				}

				e.Emit(string(b))

				for _, u := range r.Refers {
					r, _ := crawler.NewRequest(u.String())

					c.Add(r)
				}

				if rate > 0 {
					time.Sleep(time.Duration(rate) * time.Millisecond)
				}
			}
		default:
			fs.ServeHTTP(w, r)
		}
	})
}
