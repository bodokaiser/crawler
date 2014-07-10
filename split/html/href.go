package html

import "bytes"

var (
	HrefOpen  = []byte(`href="`)
	HrefClose = []byte(`"`)
)

// SplitHref returns the tokens inside a HTML href attribute
// and returns the calculated offset.
// It can be used with bufio.Scanner to extract all hrefs
// from a HTML website.
func SplitHref(b []byte, eof bool) (int, []byte, error) {
	// find the first occurence of
	// our href opening sequence
	i := bytes.Index(b, HrefOpen)

	// if index occurs
	if i != -1 {
		// move index to after
		// href opening sequence
		i += len(HrefOpen)

		// now find the index from the opening
		// sequence where the href gets closed
		n := bytes.Index(b[i:], HrefClose) + i

		// if closing was found and index
		// seems in range return total
		// offset with tokens inside
		if n != -1 && n > i && n < len(b) {
			return n + 1, b[i:n], nil
		}
	}

	// else just fuck the given bytes and go on
	return len(b), nil, nil
}
