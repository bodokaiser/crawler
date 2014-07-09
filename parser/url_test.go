package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestURLParserWrite(t *testing.T) {
	Convey("Given a string with no url", t, func() {
		str := `<a href="foobar">Email</a>`

		Convey("Write()", func() {
			parser, err := parseURL(str)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Result()", func() {
				result := parser.Result()

				Convey("Should return no url", func() {
					So(result, ShouldBeEmpty)
					So(result, ShouldHaveSameTypeAs, []string{})
				})
			})
		})
	})
	Convey("Given a string with single url", t, func() {
		str := `<a href="/about.html">About</a>`

		Convey("Write()", func() {
			parser, err := parseURL(str)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Result()", func() {
				result := parser.Result()

				Convey("Should return url", func() {
					So(result, ShouldResemble, []string{"/about.html"})
				})
			})
		})
	})
	Convey("Given a string with multiple urls", t, func() {
		str := `<a href="/about.html">about</a><a href="/legal">legal</a>`

		Convey("Write()", func() {
			parser, err := parseURL(str)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Result()", func() {
				result := parser.Result()

				Convey("Should return urls", func() {
					So(result, ShouldResemble, []string{"/about.html", "/legal"})
				})
			})
		})
	})
}

func parseURL(s string) (*URLParser, error) {
	p := &URLParser{}

	_, err := p.Write([]byte(s))

	return p, err
}
