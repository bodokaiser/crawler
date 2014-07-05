package robot

import (
	"strings"

	. "github.com/puerkitobio/goquery"

	"github.com/bodokaiser/go-crawler/conf"
)

type Robot struct {
	Origin  string
	Results []string
}

func New() *Robot {
	return &Robot{}
}

func (r *Robot) Open(c conf.Conf) error {
	r.Origin = c["entry"]

	doc, err := NewDocument(r.Origin)

	if err != nil {
		return err
	}

	r.Results = doc.Find("a").Map(func(i int, el *Selection) string {
		attr, _ := el.Attr("href")

		// put the hostname in front if we have local path
		if strings.Index(attr, "/") == 1 {
			attr = r.Origin + attr
		}

		return attr
	})

	return nil
}
