package httpd

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	s.handler = NewEventHandler()
}

func (s *EventSuite) TestEmpty(c *check.C) {
	n, err := s.handler.Write([]byte("Hello"))

	c.Check(err, check.IsNil)
	c.Check(n, check.Equals, 0)
}

func (s *EventSuite) TestSingle(c *check.C) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	s.handler.ServeHTTP(res, req)

	n, err := s.handler.Write([]byte("Hello"))
	c.Check(err, check.IsNil)
	c.Check(n, check.Equals, 5)

	time.Sleep(time.Second)

	s.checkResponse(c, res)
}

func (s *EventSuite) TestMultiple(c *check.C) {
	req1, _ := http.NewRequest("GET", "/", nil)
	req2, _ := http.NewRequest("GET", "/", nil)

	res1 := httptest.NewRecorder()
	res2 := httptest.NewRecorder()

	s.handler.ServeHTTP(res1, req1)
	s.handler.ServeHTTP(res2, req2)

	n, err := s.handler.Write([]byte("Hello"))
	c.Check(err, check.IsNil)
	c.Check(n, check.Equals, 10)

	time.Sleep(time.Second)

	s.checkResponse(c, res1)
	s.checkResponse(c, res2)
}

func (s *EventSuite) checkResponse(c *check.C, r *httptest.ResponseRecorder) {
	c.Check(r.Code, check.Equals, http.StatusOK)

	c.Check(r.HeaderMap.Get("Connection"), check.Equals, "keep-alive")
	c.Check(r.HeaderMap.Get("Content-Type"), check.Equals, "text/event-stream")
	c.Check(r.HeaderMap.Get("Cache-Control"), check.Equals, "no-cache")

	c.Check(r.Body.String(), check.Equals, "data: Hello\n\n")
}
