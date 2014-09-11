package store

import "errors"

var drivers = make(map[string]Driver)

type Store interface {
	Page() (PageStore, error)
	Close() error
}

type Driver interface {
	Open(string) (Store, error)
}

var ErrDupRow = errors.New("duplicate row")

func Open(name, url string) (Store, error) {
	d, ok := drivers[name]
	if !ok {
		panic("store: driver does not exist " + name)
	}

	return d.Open(url)
}

func Register(name string, d Driver) {
	if d == nil {
		panic("store: driver is nil")
	}
	if _, ok := drivers[name]; ok {
		panic("store: driver already registered " + name)
	}
	drivers[name] = d
}
