package loadbalancer

type workerPool []*Worker

func (p workerPool) Len() int { return len(p) }

func (p workerPool) Less(i, j int) bool {
	return p[i].Pending < p[j].Pending
}

func (p *workerPool) Swap(i, j int) {
	pl := *p
	pl[i], pl[j] = pl[j], pl[i]
}

func (p *workerPool) Push(x interface{}) {
	item := x.(*Worker)
	*p = append(*p, item)
}

func (p *workerPool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	*p = old[:n-1]
	return item
}
