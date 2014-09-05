package list

import (
	"net/url"
	"strings"
)

type Item struct {
	origin string
	refers []string
}

func NewItem() *Item {
	return &Item{
		refers: make([]string, 0),
	}
}

func NewItemFromUrl(u string) (*Item, error) {
	uri, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	i := NewItem()
	i.SetOrigin(uri)

	return i, nil
}

func (i *Item) Origin() string {
	return i.origin
}

func (i *Item) SetOrigin(u *url.URL) {
	normalize(u)

	i.origin = u.String()
}

func (i *Item) Refers() []string {
	return i.refers
}

func (i *Item) AddRefer(ref string) {
	if u, err := url.Parse(ref); err == nil {
		normalize(u)

		i.refers = append(i.refers, u.String())
	}
}

func normalize(u *url.URL) {
	u.Fragment = ""
	u.Host = strings.ToLower(u.Host)
	u.Scheme = strings.ToLower(u.Scheme)
}
