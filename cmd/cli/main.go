package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/bodokaiser/crawler"
)

var wg = new(sync.WaitGroup)

func main() {
	var entry string

	flag.StringVar(&entry, "entry", "", "")
	flag.Parse()

	c := crawler.New()
	c.Put(&crawler.Page{
		Origin: entry,
	})

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go crawl(c)
	}

	wg.Wait()
}

func crawl(c *crawler.Crawler) {
	for {
		p := c.Get()

		fmt.Printf("%s\n", p.Origin)

		for _, r := range p.Refers {
			p := &crawler.Page{}
			p.SetOrigin(r)

			c.Put(p)
		}
	}

	wg.Done()
}
