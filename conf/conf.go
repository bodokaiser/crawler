package conf

import (
	"errors"
	"net/url"
)

var (
	MissingEntryArgument = errors.New("Please pass a entry url as argument.")
	InvalidEntryArgument = errors.New("Please pass a valid entry url as argument.")
)

type Conf map[string]string

func NewConf() Conf {
	return Conf{}
}

func (c Conf) Parse(a []string) error {
	e := a[len(a)-1]

	if len(e) < 8 {
		return MissingEntryArgument
	}
	if _, err := url.Parse(e); err != nil {
		return InvalidEntryArgument
	}

	c["entry"] = e

	return nil
}
