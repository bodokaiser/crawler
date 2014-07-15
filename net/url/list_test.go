package url

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestList(t *testing.T) {
	Convey("Given an empty list", t, func() {
		list := NewList()

		Convey("Has()", func() {
			result := list.Has("http://www.google.com")

			Convey("Should return false", func() {
				So(result, ShouldBeFalse)
			})
		})
	})
	Convey("Given a list with item", t, func() {
		list := NewList()
		list.Add("http://www.example.org")

		Convey("Has()", func() {
			result := list.Has("http://www.EXAMPLE.org#hello")

			Convey("Should return true", func() {
				So(result, ShouldBeTrue)
			})
		})
	})
}
