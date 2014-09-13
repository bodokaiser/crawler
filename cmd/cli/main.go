package main

import (
	"fmt"
	"log"

	"github.com/bodokaiser/crawler"
	_ "github.com/bodokaiser/crawler/store/mysql"
)

func main() {
	conf := &crawler.Config{}
	conf.Parse()

	if err := conf.Check(); err != nil {
		log.Fatalf("Error parsing parameters: %s.\n", err)

		return
	}

	c := crawler.New()
	c.Put(crawler.NewPage(conf.Origin))

	for {
		p := c.Get()

		fmt.Printf("Origin: %s\n", p.Origin())

		for _, r := range p.Refers() {
			fmt.Printf("Refer: %s\n", r)

			c.Put(crawler.NewPage(r))
		}
	}
}
