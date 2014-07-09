package parser

import "fmt"

type EmailParser struct {
	Results []string
}

func (p *EmailParser) Write(b []byte) (int, error) {
	size := len(b)

	// traverse through byte range
	for position, char := range b {
		// until we find an @ char
		if char == '@' {
			var offset, length int

			// go back through the chars until we have something illegal
			for offset = position - 1; offset > 0; offset-- {
				char := b[offset]

				if !alpha(char) {
					break
				}
			}

			// go forward through the chars until we have something illegal
			for length = position + 1; length < size; length++ {
				char := b[length]

				if char != '.' && !alpha(char) {
					break
				}
			}

			fmt.Printf("\nfound email: %s\n", b[offset+1:length])

			// add extracted email to result set if not existent
			p.Results = append(p.Results, string(b[offset+1:length]))
		}
	}

	return size, nil
}

func alpha(c byte) bool {
	return 'a' <= c && 'z' >= c
}
