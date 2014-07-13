package conf

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfParse(t *testing.T) {
	Convey("Given an url argument", t, func() {
		args := []string{"cmd", "http://www.google.com"}

		Convey("Parse()", func() {
			err := NewConf().Parse(args)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
	Convey("Given an invalid url argument", t, func() {
		args := []string{"cmd", "????"}

		Convey("Parse()", func() {
			err := NewConf().Parse(args)

			Convey("Should return error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
