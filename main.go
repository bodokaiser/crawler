package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bodokaiser/gerenuk/conf"
	"github.com/bodokaiser/gerenuk/robot"
)

func main() {
	c := conf.New()

	if err := c.Flags(); err != nil {
		log.Fatal(err)
	}

	go request(c["url"])

	time.Sleep(5 * time.Second)
}

func request(url string) error {
	r := robot.New()

	r.Handle(func(r *http.Response) {
		fmt.Printf("Got response: %v", r)
	})

	r.Open(url)
}
