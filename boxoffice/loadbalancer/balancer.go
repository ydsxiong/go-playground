package loadbalancer

import "container/heap"

type Balancer struct {
	pool workerPool
	done chan *Worker
}

func InitBalancer(nWorker int) *Balancer {
	done := make(chan *worker, nWorker)
	b := &Balancer{make(workerPool, nWorker), done}
	for i := 1; i <= nWorker; i++ {
		w := &worker{work: make(chan ClientRequest)}
		w.Id = i
		heap.Push(&b.pool, w)

		go w.DoWork(done)
	}
	return b
}

func (b *Balancer) Balance(req chan ClientRequest) {
	for {
		select {
		case request := <-req:
			b.dispatch(request)
		case w := <-b.done:
			b.complete(w)
		}
	}
}

func (b *Balancer) dispatch(request ClientRequest) {
	w := heap.Pop(&b.pool).(worker)
	w.work <- request
	w.pending++
	heap.Push(&b.pool, w)

}

func (b *Balancer) complete(w *worker) {
	w.pending--
}
