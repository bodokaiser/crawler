package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEmailParserWrite(t *testing.T) {
	Convey("Given a string with no email", t, func() {
		str := `<a href="foobar">Email</a>`

		Convey("Write()", func() {
			parser, err := parseEmail(str)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should have no email in results", func() {
				So(parser.Results, ShouldBeEmpty)
			})
		})
	})
	Convey("Given a string with single email", t, func() {
		str := `<a href="me@example.org">Email</a>`

		Convey("Write()", func() {
			parser, err := parseEmail(str)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should append email to result", func() {
				So(parser.Results, ShouldResemble, []string{"me@example.org"})
			})
		})
	})
	Convey("Given a string with multiple emails", t, func() {
		str := `<a href="me@example.org">mail</a>info@example.org`

		Convey("Write()", func() {
			parser, err := parseEmail(str)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should append email to result once", func() {
				So(parser.Results, ShouldResemble, []string{"me@example.org", "info@example.org"})
			})
		})
	})
	SkipConvey("Given a string with duplicate emails", t, func() {
		str := `<a href="me@example.org">me@example.org</a>`

		Convey("Write()", func() {
			parser, err := parseEmail(str)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should append email to result once", func() {
				So(parser.Results, ShouldResemble, []string{"me@example.org"})
			})
		})
	})
}

func parseEmail(s string) (*EmailParser, error) {
	p := &EmailParser{}

	_, err := p.Write([]byte(s))

	return p, err
}
