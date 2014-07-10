package html

import "bytes"

var (
	URLPrefix = []byte(`/`)
)

// Extends SplitHref to only handle local urls.
func SplitURL(b []byte, eof bool) (int, []byte, error) {
	i, t, err := SplitHref(b, eof)

	if t != nil && bytes.HasPrefix(t, URLPrefix) {
		return i, t, nil
	}

	return i, nil, err
}
