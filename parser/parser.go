package parser

// Parser takes an io.Reader and scans on it for
// interesting tokens.
type Parser interface {
	// Next returns the next result or nil if EOF.
	Next() *Result
}

// Result will be returned by a parser.
type Result struct {
	Value []byte
}

// Returns string representation of result.
func (r *Result) String() string {
	return string(r.Value)
}
