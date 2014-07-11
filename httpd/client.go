package httpd

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
)

type Client struct {
	request *http.Request

	worker []*clientWorker
}

func NewClient(url string) *Client {
	r, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	c := &Client{
		request: r,
	}

	return c
}

func (c *Client) Open() error {
	r, err := http.DefaultClient.Do(c.request)

	if err != nil {
		return err
	}

	for _, w := range c.worker {
		go w.Spawn(c.request, r)
	}

	s := bufio.NewScanner(r.Body)

	for s.Scan() {
		b := s.Bytes()

		for _, w := range c.worker {
			w.Stream <- b
		}
	}

	for _, w := range c.worker {
		close(w.Stream)
	}

	return nil
}

type ClientHandle func(*ClientResult)

func (c *Client) Handle(ch ClientHandle) {
	cw := &clientWorker{
		Handle: ch,
		Stream: make(chan []byte),
	}

	c.worker = append(c.worker, cw)
}

type ClientResult struct {
	Code int
	Host string
	Body io.Reader
}

type clientWorker struct {
	Handle ClientHandle
	Stream chan []byte
}

func (cw *clientWorker) Spawn(req *http.Request, res *http.Response) {
	cr := &ClientResult{
		Host: req.Host,
		Code: res.StatusCode,
	}

	for b := range cw.Stream {
		cr.Body = bytes.NewReader(b)

		cw.Handle(cr)
	}
}
