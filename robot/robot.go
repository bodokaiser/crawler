package robot

import (
	. "github.com/puerkitobio/goquery"

	"github.com/bodokaiser/go-crawler/conf"
)

type Robot struct {
	Origin    string
	Results   []string
	selector  string
	attribute string
}

func New(conf conf.Conf) *Robot {
	o := conf["entry"]
	s := conf["selector"]
	a := conf["attribute"]

	return &Robot{o, []string{}, s, a}
}

func (r *Robot) Open() error {
	doc, err := NewDocument(r.Origin)

	if err != nil {
		return err
	}

	doc.Find(r.selector).Each(func(i int, el *Selection) {
		attr, exists := el.Attr(r.attribute)

		if exists {
			r.Results = append(r.Results, attr)
		}
	})

	return nil
}
