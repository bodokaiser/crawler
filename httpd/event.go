package httpd

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

// MaxEvents defines amount of maximum allowed client which are handled by
// EventHandler.
var MaxEvents = 10

// Errors supported by EventHandler.
var (
	ErrBadStreamRequest = errors.New("request does not support streaming")
	ErrMaxEventsReached = errors.New("too many event clients connected")
)

// Type EventHandler provides an high level interface to support HTTP event
// streams for Server Sent Events.
type EventHandler struct {
	// Open event streams.
	events []chan []byte
	// Mutex to regulate channel allocation.
	*sync.Mutex
}

// Returns an new initialized EventHandler.
func NewEventHandler() *EventHandler {
	return &EventHandler{
		events: make([]chan []byte, MaxEvents),
		Mutex:  &sync.Mutex{},
	}
}

// Writes the provided data to all open event streams.
func (h *EventHandler) Write(b []byte) (int, error) {
	n := 0

	for _, c := range h.events {
		if c != nil {
			c <- b

			n += len(b)
		}
	}

	return n, nil
}

// Allocates a new channel if possible to listen for writes.
func (h *EventHandler) open() (chan []byte, error) {
	h.Lock()
	defer h.Unlock()

	for i, ch := range h.events {
		if ch == nil {
			h.events[i] = make(chan []byte)

			return h.events[i], nil
		}
	}

	return nil, ErrMaxEventsReached
}

func (h *EventHandler) free(c chan []byte) {
	h.Lock()

	for i, ch := range h.events {
		if c == ch {
			h.events[i] = nil
		}

		close(c)
	}

	h.Unlock()
}

// Implements http.Handler interface which will enable an event stream which can
// receive data written to EventHandler.
func (h *EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f, ok := w.(http.Flusher); ok {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		c, err := h.open()

		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)

			return
		}

		go (func() {
			for b := range c {
				_, err := fmt.Fprintf(w, "data: %s\n\n", string(b))

				if err != nil {
					h.free(c)
				}

				f.Flush()
			}
		})()

		return
	}

	http.Error(w, ErrBadStreamRequest.Error(), http.StatusInternalServerError)
}
