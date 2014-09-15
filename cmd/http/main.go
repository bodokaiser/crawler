package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"time"

	"github.com/bodokaiser/crawler"
	"github.com/bodokaiser/crawler/http/event"
)

var rate, worker int
var path, addr string

func main() {
	flag.IntVar(&rate, "rate", 0, "")
	flag.IntVar(&worker, "worker", 100, "")
	flag.StringVar(&path, "path", "", "")
	flag.StringVar(&addr, "addr", ":3000", "")
	flag.Parse()

	http.Handle("/", Handler())
	http.ListenAndServe(addr, nil)
}

func Handler() http.Handler {
	ch := make(chan *crawler.Request)

	c := crawler.New()
	c.Run(worker)

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

			c.Do(r)

			go wait(r, ch)
		case "/events":
			e, err := event.NewStream(r, w)
			if err != nil {
				panic(err)
				return
			}

			for r := range ch {
				b, err := json.Marshal(r)
				if err != nil {
					panic(err)
					return
				}

				e.Emit(string(b))

				for _, u := range r.Refers {
					r, _ := crawler.NewRequest(u.String())

					c.Do(r)

					go wait(r, ch)
				}

				if rate > 0 {
					time.Sleep(time.Duration(rate) * time.Second)
				}
			}
		default:
			http.ServeFile(w, r, path+"/index.html")
		}
	})
}

func wait(r *crawler.Request, out chan<- *crawler.Request) {
	<-r.Done

	out <- r
}
