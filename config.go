package crawler

import (
	"errors"
	"flag"
)

type Config struct {
	Entry string
}

func (c *Config) Parse() error {
	flag.StringVar(&c.Entry, "entry", "", "entry url")
	flag.Parse()

	return c.Check()
}

func (c *Config) Check() error {
	if c.Entry == "" {
		return errors.New("invalid entry url")
	}

	return nil
}
