package crawler

import (
	"net/url"
	"strings"
)

type Page struct {
	Origin string
	Refers []string
}

func (p *Page) SetOrigin(url string) {
	if len(url) > 0 {
		p.Origin = normalize(url)
	}
}

func (p *Page) Has(url string) bool {
	url = normalize(url)

	for _, ref := range p.Refers {
		if ref == url {
			return true
		}
	}

	return false
}

func (p *Page) Push(urlStr string) {
	switch {
	case strings.HasPrefix(urlStr, "/"):
		ref, err := url.Parse(urlStr)
		if err != nil {
			break
		}
		org, _ := url.Parse(p.Origin)
		org.Path = ref.Path
		urlStr = org.String()

		fallthrough
	case strings.HasPrefix(urlStr, "http"):
		p.Refers = append(p.Refers, urlStr)
	}
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

type Pages []*Page

func (ps *Pages) Has(p *Page) bool {
	for _, x := range *ps {
		if x.Origin == p.Origin {
			return true
		}
	}

	return false
}

func (ps *Pages) Add(p *Page) {
	*ps = append(*ps, p)
}
