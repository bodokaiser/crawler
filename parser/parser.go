package parser

import "io"

// New is a typed factory function.
// This makes it possible to not avoid reflection.
type New func(r io.Reader) Parser

// A Parser takes an io.Reader and scans on it for
// interesting tokens.
type Parser interface {
	// Next returns the next result or nil if EOF.
	Next() Result
}

// Result is returned by the Parser. We use a struct
// over raw types as this allows us to include meta
// data like the parser type and it gives us improved
// type safety through functions.
type Result interface {
	Value() string
}
