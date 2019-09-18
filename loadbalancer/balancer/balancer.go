package balancer

import (
	"container/heap"
	"fmt"
)

type Balancer struct {
	pool workerPool
	done chan *worker
}

func InitBalancer(nWorker, nRequester int) *Balancer {
	done := make(chan *worker, nWorker)
	b := &Balancer{make(workerPool, 0, nWorker), done}

	// set up workers to communicate with balancer
	for i := 0; i < nWorker; i++ {
		w := &worker{wok: make(chan Request, nRequester)}
		w.inx = i + 1
		heap.Push(&b.pool, w)

		go w.doWork(done)
	}
	return b
}

func (b *Balancer) Balance(req chan Request) {
	for {
		select {
		case request := <-req:
			b.dispatch(request)
		case w := <-b.done:
			b.completed(w)
		}
		b.Print()
	}
}

func (b *Balancer) dispatch(req Request) {
	// grab least loaded worker
	w := heap.Pop(&b.pool).(*worker)
	w.wok <- req
	w.pending++
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *worker) {
	w.pending--
	// remove id from heap
	//heap.Remove(&b.pool, w.inx)

	// put it back to the pool
	//heap.Push(&b.pool, w)
}

func (b *Balancer) Print() {
	sum := 0
	sumsq := 0
	// pring pending stats per worker
	for _, w := range b.pool {
		fmt.Printf("%d ", w.pending)
		sum += w.pending
		sumsq += w.pending * w.pending
	}
	// print avg for worker pool
	avg := float64(sum) / float64(len(b.pool))
	variance := float64(sumsq)/float64(len(b.pool)) - avg*avg
	fmt.Printf(" %.2f %.2f\n", avg, variance)
	fmt.Println("finished one round of stats display\n")
}
