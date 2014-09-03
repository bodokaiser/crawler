package html

import "bytes"

var (
	HrefPrefix = []byte(`href="`)
	HrefSuffix = []byte(`"`)
)

func ScanHref(b []byte, eof bool) (int, []byte, error) {
	i := bytes.Index(b, HrefPrefix)

	if i != -1 {
		i += len(HrefPrefix)
		n := bytes.Index(b[i:], HrefSuffix) + i

		// check if slice is in range
		if n != -1 && n > i && n < len(b) {
			return n + 1, b[i:n], nil
		}
	}

	return len(b), nil, nil
}
