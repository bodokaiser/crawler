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

type Config struct {
	DB  string
	Url string
}

type Crawler struct {
	active bool
	result chan string
	store  *store.Store
	pool   *pool.WorkPool
}

func NewCrawler(c Config) (*Crawler, error) {
	s, err := store.Open(c.DB)
	if err != nil {
		return nil, err
	}

	err = s.DropTables()
	if err != nil {
		return nil, err
	}
	err = s.EnsureTables()
	if err != nil {
		return nil, err
	}

	cr := &Crawler{
		store:  s,
		result: make(chan string, 30),
	}

	return cr, cr.put(c.Url)
}

func (c *Crawler) put(url string) error {
	err := c.store.Insert(url)
	if err != nil {
		return err
	}

	if c.active == false {
		c.pool = pool.NewWorkPool(pool.Config{
			New: c.work,
		})

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

			tx, err := store.Begin()
			if err != nil {
				panic(err)
			}

			err = work(res, tx)
			if err != nil {
				panic(err)
			}

			fmt.Printf("ending work\n")

			return nil, err
		},
	}
}

func work(result chan<- string, tx *store.Tx) error {
	fmt.Printf("request to: %s\n", tx.Origin())
	req, err := http.NewRequest("GET", tx.Origin(), nil)
	if err != nil {
		tx.Abort()

		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		tx.Abort()

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

			if err := tx.AddRefer(t); err != nil {
				if err != store.ErrRefExists {
					return err
				}
			}
		}
	}

	return tx.Commit()
}
