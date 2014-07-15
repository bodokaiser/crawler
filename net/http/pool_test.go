package httpd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPool(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Handle))

	request1, _ := http.NewRequest("GET", server.URL, nil)
	request2, _ := http.NewRequest("GET", server.URL, nil)
	request3, _ := http.NewRequest("GET", "error.dns", nil)

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
				So(req.URL.String(), ShouldEqual, server.URL)
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
				So(req.URL.String(), ShouldEqual, server.URL)
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
					So(req.URL.String(), ShouldEqual, server.URL)
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
			})
			Convey("Should return no response", func() {
				So(res, ShouldBeNil)
			})
		})
	})
}

func Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
