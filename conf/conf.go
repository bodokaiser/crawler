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
	s := flag.String("selector", "", "Query selector to filter elements.")
	a := flag.String("attribute", "", "Attribute to extract from elements.")

	flag.Parse()

	if len(e) == 0 {
		return errors.New(`"entry" flag is mandatory.`)
	}
	if len(*s) == 0 {
		return errors.New(`"selector" flag is mandatory.`)
	}
	if len(*a) == 0 {
		return errors.New(`"attribute" flag is mandatory.`)
	}

	c["entry"] = e
	c["selector"] = *s
	c["attribute"] = *a

	return nil
}
