package main

import (
	"flag"

	"github.com/bodokaiser/gerenuk/httpd"
)

func main() {
	s := httpd.NewServer()

	e := flag.String("host", ":3000", "The host to listen on.")
	flag.Parse()

	s.Listen(*e)
}
