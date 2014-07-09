package parser

type EmailParser struct {
	results []string
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

			// add extracted email to result set if not existent
			p.results = append(p.results, string(b[offset+1:length]))
		}
	}

	return size, nil
}

func (p *EmailParser) Result() []string {
	return p.results
}

func alpha(c byte) bool {
	return 'a' <= c && 'z' >= c
}
