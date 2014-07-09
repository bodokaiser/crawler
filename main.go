package main

import (
	"log"
	"time"

	"github.com/bodokaiser/gerenuk/conf"
	"github.com/bodokaiser/gerenuk/parser"
	"github.com/bodokaiser/gerenuk/robots"
)

func main() {
	c := conf.New()

	if err := c.Flags(); err != nil {
		log.Fatal(err)
	}

	go request(c["url"])

	time.Sleep(5 * time.Second)
}

func request(url string) {
	r := robots.NewRobot()

	r.RegisterParser(&parser.EmailParser{})

	if err := r.Open(url); err != nil {
		log.Fatal(err)
	}
}
