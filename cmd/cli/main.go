package main

import (
	"log"

	. "github.com/bodokaiser/gerenuk"
	"github.com/bodokaiser/gerenuk/pipe"
	"github.com/bodokaiser/gerenuk/store"
	_ "github.com/bodokaiser/gerenuk/store/mysql"
)

var conf = &Config{}

var pages store.PageStore

func main() {
	conf.Parse()
	if err := conf.Check(); err != nil {
		log.Fatalf("Error parsing parameters: %s.\n", err)

		return
	}

	s, err := store.Open("mysql", conf.Store.Url)
	if err != nil {
		log.Fatalf("Error connecting to store: %s.\n", err)

		return
	}
	pages, err = s.Page()
	if err != nil {
		log.Fatalf("Error setting up store: %s.\n", err)

		return
	}

	p := Page{}
	p.SetOrigin(conf.Origin)

	c := NewCrawler()
	c.Put(p)

	for e := range c.Pipe.Pipe(pipe.StageFunc(insert)).Listen() {
		p := e.(Page)

		//fmt.Printf("%s: %s\n", p.Origin(), strings.Join(p.Refers(), ", "))

		for _, r := range p.Refers() {
			p := Page{}
			p.SetOrigin(r)

			c.Put(p)
		}
	}
}

func insert(in <-chan pipe.Event) <-chan pipe.Event {
	out := make(chan pipe.Event)

	go func(in <-chan pipe.Event, out chan<- pipe.Event) {
		for e := range in {
			p := e.(Page)

			if err := pages.Insert(&p); err != nil {
				log.Fatalf("Error inserting page: %s.\n", err)

				break
			} else {
				out <- p
			}
		}

		close(out)
	}(in, out)

	return out
}
