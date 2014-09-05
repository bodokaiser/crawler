package httpd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Error used when response writer does not support Flusher.
var ErrBadStream = errors.New("bad event stream")

// Type EventStream provdes an abstraction for http responses of content type
// event stream.
type EventStream struct {
	writer  io.Writer
	flusher http.Flusher
	request *http.Request
}

// Returns initialized EventStream and error depending if response supports cast
// to Flusher.
func NewEventStream(r *http.Request, w http.ResponseWriter) (*EventStream, error) {
	f, ok := w.(http.Flusher)

	if ok {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
	} else {
		return nil, ErrBadStream
	}

	return &EventStream{
		writer:  w,
		flusher: f,
		request: r,
	}, nil
}

// Emits a new event to the event stream.
// Returns error if write fails.
func (e *EventStream) Emit(data string) error {
	_, err := fmt.Fprintf(e.writer, "data: %s\n\n", data)

	if err != nil {
		return err
	}

	e.flusher.Flush()

	return nil
}
