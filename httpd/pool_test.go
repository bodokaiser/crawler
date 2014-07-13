package httpd

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	request1, _ = http.NewRequest("GET", "http://www.google.de", nil)
	request2, _ = http.NewRequest("GET", "http://www.google.com", nil)
	request3, _ = http.NewRequest("GET", "http://iam.notexistent", nil)
)

func TestPool(t *testing.T) {
	Convey("Given an empty pool", t, func() {
		p := NewPool()

		Convey("Get()", func() {
			req, res, err := p.Get()

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should return no request", func() {
				So(req, ShouldBeNil)
			})
			Convey("Should return no response", func() {
				So(res, ShouldBeNil)
			})
		})
	})
	Convey("Given a pool with one request", t, func() {
		p := NewPool()
		p.Add(request1)
		p.Run()

		Convey("Get()", func() {
			req, res, err := p.Get()

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should return request", func() {
				So(req, ShouldHaveSameTypeAs, &http.Request{})
				So(req.Host, ShouldEqual, "www.google.de")
			})
			Convey("Should return response", func() {
				So(res, ShouldHaveSameTypeAs, &http.Response{})
				So(res.StatusCode, ShouldEqual, 200)
			})

			Convey("Get()", func() {
				req, res, err := p.Get()

				Convey("Should return no error", func() {
					So(err, ShouldBeNil)
				})
				Convey("Should return no request", func() {
					So(req, ShouldBeNil)
				})
				Convey("Should return no response", func() {
					So(res, ShouldBeNil)
				})
			})
		})
	})
	Convey("Given a pool with two requests", t, func() {
		p := NewPool()
		p.Add(request1)
		p.Add(request2)
		p.Run()

		Convey("Get()", func() {
			req, res, err := p.Get()

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should return request", func() {
				So(req, ShouldHaveSameTypeAs, &http.Request{})
				So(req.Host, ShouldEqual, "www.google.de")
			})
			Convey("Should return response", func() {
				So(res, ShouldHaveSameTypeAs, &http.Response{})
				So(res.StatusCode, ShouldEqual, 200)
			})

			Convey("Get()", func() {
				req, res, err := p.Get()

				Convey("Should return no error", func() {
					So(err, ShouldBeNil)
				})
				Convey("Should return request", func() {
					So(req, ShouldHaveSameTypeAs, &http.Request{})
					So(req.Host, ShouldEqual, "www.google.com")
				})
				Convey("Should return response", func() {
					So(res, ShouldHaveSameTypeAs, &http.Response{})
					So(res.StatusCode, ShouldEqual, 200)
				})
				Convey("Get()", func() {
					req, res, err := p.Get()

					Convey("Should return no error", func() {
						So(err, ShouldBeNil)
					})
					Convey("Should return no request", func() {
						So(req, ShouldBeNil)
					})
					Convey("Should return no response", func() {
						So(res, ShouldBeNil)
					})
				})
			})
		})
	})
	Convey("Given a pool with a invalid request", t, func() {
		p := NewPool()
		p.Add(request3)
		p.Run()

		Convey("Get()", func() {
			req, res, err := p.Get()

			Convey("Should return error", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("Should return no request", func() {
				So(req, ShouldHaveSameTypeAs, &http.Request{})
				So(req.Host, ShouldEqual, "iam.notexistent")
			})
			Convey("Should return no response", func() {
				So(res, ShouldBeNil)
			})
		})
	})
}
