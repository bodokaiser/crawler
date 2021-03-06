package event

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Error used when response writer does not support Flusher.
var ErrBadStream = errors.New("bad event stream")

// Type Stream provdes an abstraction for http responses of content type
// event stream.
type Stream struct {
	writer  io.Writer
	flusher http.Flusher
	request *http.Request
}

// Returns initialized EventStream and error depending if response supports cast
// to Flusher.
func NewStream(r *http.Request, w http.ResponseWriter) (*Stream, error) {
	f, ok := w.(http.Flusher)

	if ok {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
	} else {
		return nil, ErrBadStream
	}

	return &Stream{
		writer:  w,
		flusher: f,
		request: r,
	}, nil
}

// Emits a new event to the event stream.
// Returns error if write fails.
func (e *Stream) Emit(data string) error {
	_, err := fmt.Fprintf(e.writer, "data: %s\n\n", data)

	if err != nil {
		return err
	}

	e.flusher.Flush()

	return nil
}
