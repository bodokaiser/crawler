package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHrefResult(t *testing.T) {
	Convey("Given byte array", t, func() {
		result := &Result{[]byte(`me@foo.org`)}

		Convey("String()", func() {
			value := result.String()

			Convey("Should return value as string", func() {
				So(value, ShouldEqual, "me@foo.org")
			})
		})
	})
}
