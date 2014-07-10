package html

import "bytes"

var (
	MailPrefix = []byte(`mailto:`)
)

// Extends SplitHref() by checking if attribute is prefixed
// with "mailto" macro else it omits the given bytes.
func SplitMail(b []byte, eof bool) (int, []byte, error) {
	i, t, err := SplitHref(b, eof)

	if t != nil && bytes.HasPrefix(t, MailPrefix) {
		off := len(MailPrefix)

		return i, t[off:], nil
	}

	return i, nil, err
}
