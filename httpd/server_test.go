package httpd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServerServeHTTP(t *testing.T) {
	StaticDir = "public"

	server := NewServer()

	Convey("Given a request to /events", t, func() {
		req, res := request("/events")

		Convey("ServeHTTP()", func() {
			server.ServeHTTP(res, req)

			Convey("Should respond OK", func() {
				So(res.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
	Convey("Given a request to /", t, func() {
		req, res := request("/")
		req.Header.Add("Accept", "text/html")

		Convey("ServeHTTP()", func() {
			server.ServeHTTP(res, req)

			Convey("Should respond OK", func() {
				So(res.Code, ShouldEqual, http.StatusOK)
			})
			Convey("Should respond html", func() {
				So(res.Body.String(), ShouldStartWith, "<!DOCTYPE html>")
			})
		})
	})
	Convey("Given a request to /javascripts/script.js", t, func() {
		req, res := request("/javascripts/script.js")

		Convey("ServeHTTP()", func() {
			server.ServeHTTP(res, req)

			Convey("Should respond OK", func() {
				So(res.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func request(url string) (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest("GET", url, nil)

	return req, httptest.NewRecorder()
}
