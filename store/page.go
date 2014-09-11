package store

type Page interface {
	Origin() string
	Refers() []string
}

type PageStore interface {
	Insert(Page) error
}
