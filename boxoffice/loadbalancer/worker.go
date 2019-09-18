package loadbalancer

type worker struct {
	id      int
	pending int
	work    chan ClientRequest
}

func (w *worker) DoWork(done chan *worker) {
	for {
		req := <-w.work

		req.Resp <- nil
		done <- w
	}
}
