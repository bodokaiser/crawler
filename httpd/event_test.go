package httpd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/check.v1"
)

func TestEvent(t *testing.T) {
	check.Suite(&EventSuite{})
	check.TestingT(t)
}

type EventSuite struct {
	handler *EventHandler
}

func (s *EventSuite) SetUpTest(c *check.C) {
	s.handler = NewEventHandler(func(w http.ResponseWriter, _ *http.Request, f http.Flusher) {
		fmt.Fprintf(w, "data: message: %s\n\n", "Hello World")

		f.Flush()
	})
}

func (s *EventSuite) TestSendEvent(c *check.C) {
	res := httptest.NewRecorder()

	SendEvent(res, "Hello World")

	c.Check(res.Body.String(), check.Equals, "data: Hello World\n\n")
}

func (s *EventSuite) TestServeHTTP(c *check.C) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	s.handler.ServeHTTP(res, req)

	c.Check(res.Code, check.Equals, http.StatusOK)

	c.Check(res.HeaderMap.Get("Connection"), check.Equals, "keep-alive")
	c.Check(res.HeaderMap.Get("Content-Type"), check.Equals, "text/event-stream")
	c.Check(res.HeaderMap.Get("Cache-Control"), check.Equals, "no-cache")

	c.Check(res.Body.String(), check.Equals, "data: message: Hello World\n\n")
}
