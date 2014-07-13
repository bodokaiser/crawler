package httpd

import (
	"fmt"
	"net/http"
	"time"
)

var (
	StaticDir = http.Dir("httpd/public")
)

type Server struct {
	static http.Handler
}

func NewServer() *Server {
	return &Server{
		static: http.FileServer(StaticDir),
	}
}

func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Server Sent Events not supported!", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "data: Message: %s\n\n", "Hello")

		f.Flush()

		time.Sleep(2 * time.Second)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := http.NewServeMux()

	m.Handle("/", s.static)
	m.HandleFunc("/events", s.events)

	m.ServeHTTP(w, r)
}

func (s *Server) Listen(addr string) {
	http.ListenAndServe(addr, s)
}
