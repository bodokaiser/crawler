package utils

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestScanHref(t *testing.T) {
	Convey("Given a byte array without href", t, func() {
		data := []byte(`<h1>Hello World</h1>`)

		Convey("ScanHref()", func() {
			offset, token, err := ScanHref(data, false)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should return no token", func() {
				So(token, ShouldBeNil)
			})
			Convey("should return offset", func() {
				So(offset, ShouldEqual, len(data))
			})
		})
	})
	Convey("Given a byte array with single href", t, func() {
		data := []byte(`<a href="foobar">Email</a>`)

		Convey("ScanHref()", func() {
			offset, token, err := ScanHref(data, false)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should return token", func() {
				So(string(token), ShouldEqual, "foobar")
			})
			Convey("should return offset", func() {
				So(offset, ShouldEqual, bytes.LastIndex(data, []byte(`">`))+1)
			})
		})
	})
	Convey("Given a byte array with multiple href", t, func() {
		data := []byte(`<a href="foobar">Email</a><a href="helloworld">Hello</a>`)

		Convey("ScanHref()", func() {
			offset1, token1, err1 := ScanHref(data, false)
			offset2, token2, err2 := ScanHref(data[offset1:], false)

			Convey("Should return no error", func() {
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
			})
			Convey("Should return token", func() {
				So(string(token1), ShouldEqual, "foobar")
				So(string(token2), ShouldEqual, "helloworld")
			})
			Convey("should return offset", func() {
				So(offset1, ShouldEqual, bytes.Index(data, []byte(`">`))+1)
				So(offset2, ShouldEqual, bytes.Index(data[offset1:], []byte(`">`))+1)
			})
		})
	})
}
