package balancer

import (
	"math"
	"math/rand"
	"time"
)

type worker struct {
	inx     int // the worker id reflecting the position of last one pushed into the pool, but why???
	pending int // loading amount used by the pool ordering in heap
	wok     chan Request
}

func (w *worker) doWork(done chan *worker) {
	for {
		req := <-w.wok
		time.Sleep(time.Duration(rand.Int63n(int64(time.Minute))))
		req.Resp <- math.Sin(float64(req.Data))
		done <- w
	}
}
