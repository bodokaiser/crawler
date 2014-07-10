package html

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSplitURL(t *testing.T) {
	Convey("Given a byte array with url", t, func() {
		data := []byte(`<a href="/foobar">Email</a>`)

		Convey("SplitURL()", func() {
			offset, token, err := SplitURL(data, false)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should return token", func() {
				So(string(token), ShouldEqual, "/foobar")
			})
			Convey("should return offset", func() {
				So(offset, ShouldEqual, bytes.Index(data, []byte(`">`))+1)
			})
		})
	})
}
