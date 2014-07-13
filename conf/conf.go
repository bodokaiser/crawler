package conf

import "flag"

type Conf struct {
	Entry   string
	Address string
}

func NewConf() *Conf {
	return &Conf{}
}

func (c *Conf) Parse(args []string) error {
	s := flag.NewFlagSet("default", flag.ContinueOnError)

	e := s.String("entry", "", "The URL of the first website to crawl.")
	a := s.String("address", "", "The address of the http server to listen.")

	if err := s.Parse(args[1:]); err != nil {
		return err
	}

	c.Entry = *e
	c.Address = *a

	return nil
}
