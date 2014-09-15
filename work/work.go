package work

// Default parameter for Worker queue size.
var MaxQueue = 10

// Type Work represents a tasks which may be executed by a Worker.
// Internally you may use channels to communicate results and errors.
type Work interface {
	Do()
}

// Type Worker executes Work concurrent.
type Worker struct {
	queue chan Work
}

// Returns initialized Worker with queue size of n.
func New() *Worker {
	return &Worker{
		queue: make(chan Work, MaxQueue),
	}
}

// Pushes work to queue where a call to Run will start execution.
func (w *Worker) Do(work Work) {
	w.queue <- work
}

// Runs work in queue on n goroutines.
func (w *Worker) Run(n int) {
	for i := 0; i < n; i++ {
		go func(queue <-chan Work) {
			for w := range queue {
				w.Do()
			}
		}(w.queue)
	}
}

// Kills running workers.
// You must recreate a Worker after Kill()..
func (w *Worker) Kill() {
	close(w.queue)
}
