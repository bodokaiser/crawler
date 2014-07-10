package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bodokaiser/gerenuk/split/html"
)

func main() {
	c := make(chan string)

	go request("http://www.satisfeet.me", c)

	for {
		result, ok := <-c

		if !ok {
			break
		}

		fmt.Printf("%s\n", result)
	}
}

func request(url string, c chan<- string) {
	r, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	go scan(r.Body, c)
}

func scan(r io.Reader, c chan<- string) {
	s := bufio.NewScanner(r)
	s.Split(html.SplitHref)

	for s.Scan() {
		t := s.Text()

		if strings.HasPrefix(t, "/") {
			c <- fmt.Sprintf("Found local link: %s", t)
		}
		if strings.HasPrefix(t, "//") || strings.HasPrefix(t, "http") {
			c <- fmt.Sprintf("Found external link: %s", t)
		}

		if strings.HasPrefix(t, "mailto:") {
			c <- fmt.Sprintf("Found email: %s", t)
		}
	}

	close(c)
}
