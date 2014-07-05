package conf

import (
	"errors"
	"flag"
	"os"
)

type Conf map[string]string

func New() Conf {
	return Conf{}
}

func (c Conf) Flags() error {
	e := os.Args[len(os.Args)-1]

	a := flag.String("addr", ":3000", `"addr" for HTTP server to listen.`)

	flag.Parse()

	if len(e) == 0 {
		return errors.New(`"entry" argument is mandatory.`)
	}

	c["addr"] = *a
	c["entry"] = e

	return nil
}
