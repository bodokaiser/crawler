package gerenuk

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/bodokaiser/gerenuk/parser/html"
	"github.com/bodokaiser/gerenuk/pool"
	"github.com/bodokaiser/gerenuk/store"
)

type Crawler struct {
	active bool
	result chan string
	pool   *pool.WorkerPool
	store  *store.Store
}

func NewCrawler(db string) (*Crawler, error) {
	s, err := store.Open(db)
	if err != nil {
		return nil, err
	}

	cr := &Crawler{
		store:  s,
		result: make(chan string, 30),
	}

	return cr, nil
}

func (c *Crawler) Put(url string) error {
	err := c.store.Put(url)
	if err != nil {
		return err
	}

	if c.active == false {
		c.pool = pool.NewWorkerPool()
		c.pool.SetNewFunc(c.work)

		c.active = true
	}

	return nil
}

func (c *Crawler) Get() string {
	return <-c.result
}

func (c *Crawler) work() *pool.Work {
	return &pool.Work{
		Params: []interface{}{
			c.result,
			c.store,
		},
		Func: func(params ...interface{}) (interface{}, error) {
			fmt.Printf("starting work\n")

			res := params[0].(chan string)
			store := params[1].(*store.Store)

			p, err := store.Get()
			if err != nil {
				panic(err)
			}

			err = work(res, p)
			if err != nil {
				panic(err)
			}

			fmt.Printf("ending work\n")

			return nil, err
		},
	}
}

func work(result chan<- string, p *store.Page) error {
	fmt.Printf("request to: %s\n", p.Origin())
	req, err := http.NewRequest("GET", p.Origin(), nil)
	if err != nil {
		p.Abort()

		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		p.Abort()

		return err
	}
	defer res.Body.Close()

	s := bufio.NewScanner(res.Body)
	s.Split(html.ScanHref)

	for s.Scan() {
		t := s.Text()

		switch {
		case strings.HasPrefix(t, "/"):
			req.URL.Path = t
			t = req.URL.String()

			fallthrough
		case strings.HasPrefix(t, "http"):
			result <- t

			if !p.HasRefer(t) {
				p.AddRefer(t)
			}
		}
	}

	return p.Commit()
}
