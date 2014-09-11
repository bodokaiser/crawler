package pipeline

import (
	"fmt"
	"sync"
)

type Stage interface {
	Process(<-chan Event) <-chan Event
}

type StageFunc func(<-chan Event) <-chan Event

func (s StageFunc) Process(in <-chan Event) <-chan Event {
	return s(in)
}

func fanIn(in <-chan Event, out *[]chan Event, mut *sync.Mutex) {
	for e := range in {
		mut.Lock()
		for _, out := range *out {
			out <- e
		}
		mut.Unlock()
	}
	for _, out := range *out {
		fmt.Printf("closing %v\n", out)
		mut.Lock()
		close(out)
		mut.Unlock()
	}
}
