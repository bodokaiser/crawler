package main

import (
	"log"

	"github.com/bodokaiser/go-crawler/conf"
	"github.com/bodokaiser/go-crawler/httpd"
	"github.com/bodokaiser/go-crawler/robot"
)

func main() {
	c := conf.New()
	r := robot.New()
	h := httpd.New()

	if err := c.Flags(); err != nil {
		log.Fatal(err)
	}

	if err := r.Open(c); err != nil {
		log.Fatal(err)
	}

	h.Listen(c)
}
