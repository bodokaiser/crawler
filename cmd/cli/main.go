package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bodokaiser/crawler"
)

var worker int

var origin string

func main() {
	flag.IntVar(&worker, "worker", 1, "")
	flag.StringVar(&origin, "origin", "", "")
	flag.Parse()

	r, err := crawler.NewRequest(origin)
	if err != nil {
		log.Fatalf("Error on initial request: %s.\n", err)

		return
	}

	ch := make(chan *crawler.Request)

	c := crawler.New()
	c.Do(r)
	c.Run(worker)

	go wait(r, ch)

	for r := range ch {
		fmt.Println(r.Origin)

		for _, u := range r.Refers {
			fmt.Println(r)

			r, _ := crawler.NewRequest(u.String())

			c.Do(r)

			go wait(r, ch)
		}
	}
}

func wait(r *crawler.Request, out chan<- *crawler.Request) {
	<-r.Done

	out <- r
}
