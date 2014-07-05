package httpd

import (
	"fmt"
	"net/http"
	"time"
)

type EventHandle struct{}

func NewEventHandle() *EventHandle {
	return &EventHandle{}
}

func (e *EventHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "SSE not supported.", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		fmt.Fprintf(w, "data: Message: %s\n\n", "Hello")

		f.Flush()

		time.Sleep(10 * 1e9)
	}
}
