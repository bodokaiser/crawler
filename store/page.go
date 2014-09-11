package store

import (
	"net/url"
	"strings"
)

type Page struct {
	Origin string
	Refers []string
}

func (p *Page) SetOrigin(url string) {
	url = normalize(url)

	if len(url) > 0 {
		p.Origin = url
	}
}

func (p *Page) HasRefer(url string) bool {
	url = normalize(url)

	for _, ref := range p.Refers {
		if ref == url {
			return true
		}
	}

	return false
}

func (p *Page) AddRefer(url string) {
	if url := normalize(url); len(url) > 0 {
		p.Refers = append(p.Refers, url)
	}
}

type PageStore interface {
	Insert(*Page) error
}

func normalize(urlStr string) string {
	uri, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	uri.Host = strings.ToLower(uri.Host)
	uri.Scheme = strings.ToLower(uri.Scheme)
	uri.Fragment = ""

	return uri.String()
}
