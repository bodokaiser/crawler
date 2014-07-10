package robot

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Robot struct {
	splitter []bufio.SplitFunc

	Results chan string
}

func NewRobot(h ...bufio.SplitFunc) *Robot {
	r := make(chan string)

	return &Robot{h, r}
}

func (r *Robot) Open(url string) error {
	res, err := http.Get(url)

	if err != nil {
		return err
	}

	go parse(res.Body, r.splitter, r.Results)

	return nil
}

func parse(r io.Reader, h []bufio.SplitFunc, o chan string) {
	channels := make([]chan []byte, len(h))

	for i, h := range h {
		channels[i] = make(chan []byte)

		go spawn(i, h, channels[i], o)
	}

	s := bufio.NewScanner(r)

	for s.Scan() {
		src := s.Bytes()
		dst := make([]byte, len(src))

		copy(dst, src)

		for _, c := range channels {
			c <- dst
		}
	}
}

func spawn(n int, h bufio.SplitFunc, i chan []byte, o chan string) {
	for b := range i {
		r := bytes.NewBuffer(b)
		s := bufio.NewScanner(r)
		s.Split(h)

		for s.Scan() {
			t := s.Text()

			o <- fmt.Sprintf("%d: %s", n, t)
		}
	}
}
