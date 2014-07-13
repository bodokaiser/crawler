package conf

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfParse(t *testing.T) {
	conf := NewConf()

	Convey("Given no arguments", t, func() {
		args := []string{"cmd"}

		Convey("Parse()", func() {
			err := conf.Parse(args)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should set Entry", func() {
				So(conf.Entry, ShouldEqual, "")
			})
			Convey("Should set Address", func() {
				So(conf.Address, ShouldEqual, "")
			})
		})
	})
	Convey("Given an entry argument", t, func() {
		args := []string{"cmd", "--entry", "http://www.google.com"}

		Convey("Parse()", func() {
			err := conf.Parse(args)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should set Entry", func() {
				So(conf.Entry, ShouldEqual, "http://www.google.com")
			})
		})
	})
	Convey("Given an address argument", t, func() {
		args := []string{"cmd", "--address", "localhost:3000"}

		Convey("Parse()", func() {
			err := conf.Parse(args)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should set Address", func() {
				So(conf.Address, ShouldEqual, "localhost:3000")
			})
		})
	})
	SkipConvey("Given an entry and address argument", t, func() {
		args := []string{"cmd", "entry", "http://foo.bar", "--address", "localhost:3000"}

		Convey("Parse()", func() {
			err := conf.Parse(args)

			Convey("Should return error", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("Should set Entry", func() {
				So(conf.Entry, ShouldEqual, "")
			})
			Convey("Should set Address", func() {
				So(conf.Address, ShouldEqual, "")
			})
		})
	})
}
