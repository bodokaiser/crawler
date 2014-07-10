package html

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSplitMail(t *testing.T) {
	Convey("Given a byte array with mail", t, func() {
		data := []byte(`<a href="mailto:me@example.org">`)

		Convey("SplitMail()", func() {
			offset, token, err := SplitMail(data, false)

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should return token", func() {
				So(string(token), ShouldEqual, "me@example.org")
			})
			Convey("should return offset", func() {
				So(offset, ShouldEqual, bytes.Index(data, []byte(`">`))+1)
			})
		})
	})
}
