package crawler

import (
	"errors"
	"flag"
)

type Config struct {
	Store struct {
		Url string
	}
	Origin string
}

func (c *Config) Parse() {
	flag.StringVar(&c.Origin, "url", "", "origin url")
	flag.StringVar(&c.Store.Url, "db", "", "database url")
	flag.Parse()
}

func (c *Config) Check() error {
	if c.Store.Url == "" {
		return errors.New("invalid database url")
	}

	return nil
}
