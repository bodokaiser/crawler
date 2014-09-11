package pipeline

import "sync"

type Event struct {
	Error  error
	Result interface{}
}

type Pipeline struct {
	mutex  *sync.Mutex
	input  chan Event
	output *[]chan Event
}

func NewPipeline() *Pipeline {
	pl := &Pipeline{
		mutex:  new(sync.Mutex),
		input:  make(chan Event),
		output: new([]chan Event),
	}
	go fanIn(pl.input, pl.output, pl.mutex)

	return pl
}

func (p *Pipeline) Emit(e Event) {
	p.input <- e
}

func (p *Pipeline) push(ch chan Event) {
	p.mutex.Lock()
	*p.output = append(*p.output, ch)
	p.mutex.Unlock()
}

func (p *Pipeline) Pipe(s Stage) *Pipeline {
	pl := &Pipeline{
		mutex:  new(sync.Mutex),
		input:  make(chan Event),
		output: new([]chan Event),
	}
	p.push(pl.input)
	go fanIn(s.Process(pl.input), pl.output, pl.mutex)

	return pl
}

func (p *Pipeline) Close() {
	close(p.input)
}

func (p *Pipeline) Listen() <-chan Event {
	out := make(chan Event)
	p.push(out)

	return out
}
