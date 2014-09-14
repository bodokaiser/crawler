package event

import (
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
	events   *Stream
	response *httptest.ResponseRecorder
}

func (s *EventSuite) SetUpTest(c *check.C) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	evt, err := NewStream(req, res)
	c.Check(err, check.IsNil)

	s.events = evt
	s.response = res
}

func (s *EventSuite) TestHeader(c *check.C) {
	r := s.response

	c.Check(r.Code, check.Equals, http.StatusOK)
	c.Check(r.HeaderMap.Get("Connection"), check.Equals, "keep-alive")
	c.Check(r.HeaderMap.Get("Content-Type"), check.Equals, "text/event-stream")
	c.Check(r.HeaderMap.Get("Cache-Control"), check.Equals, "no-cache")
}

func (s *EventSuite) TestEmit(c *check.C) {
	r := s.response
	s.events.Emit("Hello")

	c.Check(r.Body.String(), check.Equals, "data: Hello\n\n")
}
