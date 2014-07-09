package parser

import "io"

type Parser interface {
	io.Writer
}
