package httpd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSendEvent(t *testing.T) {
	Convey("Given a response", t, func() {
		res := httptest.NewRecorder()

		Convey("SendEvent()", func() {
			SendEvent(res, "Hello World")

			Convey("Should write as event stream", func() {
				So(res.Body.String(), ShouldEqual, "data: Hello World\n\n")
			})
		})
	})
}

func TestEventHandler(t *testing.T) {
	h := NewEventHandler(func(w http.ResponseWriter, _ *http.Request, f http.Flusher) {
		fmt.Fprintf(w, "data: message: %s\n\n", "Hello World")

		f.Flush()
	})

	Convey("Given a request", t, func() {
		req, res := NewEventRequestResponse("/")

		Convey("ServeHTTP()", func() {
			h.ServeHTTP(res, req)

			Convey("Should set response status", func() {
				So(res.Code, ShouldEqual, 200)
			})
			Convey("Should set response header", func() {
				So(res.HeaderMap.Get("Connection"), ShouldEqual, "keep-alive")
				So(res.HeaderMap.Get("Content-Type"), ShouldEqual, "text/event-stream")
				So(res.HeaderMap.Get("Cache-Control"), ShouldEqual, "no-cache")
			})
			Convey("Should set response body", func() {
				So(res.Body.String(), ShouldEqual, "data: message: Hello World\n\n")
			})
		})
	})
}

func NewEventRequestResponse(url string) (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest("GET", url, nil)

	return req, httptest.NewRecorder()
}
