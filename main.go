package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bodokaiser/gerenuk/conf"
	"github.com/bodokaiser/gerenuk/httpd"
	"github.com/bodokaiser/gerenuk/utils"
)

func main() {
	c := conf.NewConf()
	s := httpd.NewServer()

	if err := c.Parse(os.Args); err != nil {
		log.Fatal(err)
	}

	if len(c.Entry) != 0 {
		crawl(c.Entry)
	}

	if len(c.Address) != 0 {
		s.Listen(c.Address)
	}
}

func crawl(url string) {
	req, _ := http.NewRequest("GET", url, nil)

	p := httpd.NewPool()

	p.Add(req)
	p.Run()

	for {
		req, res, err := p.Get()

		if err != nil {
			log.Fatal(err)
		}

		s := bufio.NewScanner(res.Body)
		s.Split(utils.ScanHref)

		for s.Scan() {
			t := s.Text()

			fmt.Printf("%s href: %s\n", req.Host, t)

			if strings.HasPrefix(t, "http") {
				req, _ := http.NewRequest("GET", t, nil)

				p.Add(req)
			}
		}

		res.Body.Close()
	}
}
