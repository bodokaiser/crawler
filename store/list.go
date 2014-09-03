package store

import (
	"container/list"
	"net/url"
	"strings"
)

type List struct {
	list *list.List
}

func NewList() *List {
	return &List{
		list: list.New(),
	}
}

func (l *List) Add(url string) {
	if !l.Has(url) {
		l.list.PushFront(normalize(url))
	}
}

func (l *List) Has(url string) bool {
	url = normalize(url)

	for e := l.list.Front(); e != nil; e = e.Next() {
		if e.Value == url {
			return true
		}
	}

	return false
}

func normalize(s string) string {
	url, _ := url.Parse(strings.ToLower(s))

	url.Fragment = ""

	return url.String()
}
