package http

import (
	"fmt"
	"net/http"
)

// EventHandler is a http.Handler which sends server side events
// to all connected clients.
type EventHandler struct {
	// active defines how much clients are connected at the moment.
	active int32
	// streams contains the channels for all connected clients.
	streams []chan string
	// OnConnect defines a handler which is called when a new client
	// conntects.
	Listener EventListener
}

type EventListener func(http.ResponseWriter, *http.Request, http.Flusher)

// Returns a new EventHandler.
func NewEventHandler(l EventListener) *EventHandler {
	return &EventHandler{
		Listener: l,
	}
}

func (h *EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f, ok := w.(http.Flusher); ok {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		h.Listener(w, r, f)

		return
	}

	http.Error(w, "Server Sent Events not supported!", http.StatusInternalServerError)
}

func SendEvent(w http.ResponseWriter, m string) {
	fmt.Fprintf(w, "data: %s\n\n", m)

	w.(http.Flusher).Flush()
}
