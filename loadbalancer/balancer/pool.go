package balancer

type workerPool []*worker

/**
Heap implementations...
*/
func (p workerPool) Len() int { return len(p) }

func (p workerPool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p *workerPool) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
	//idx := a[i].GetId()
	//a[i].SetId(a[j].GetId())
	//a[j].SetId(idx)
}

func (p *workerPool) Push(x interface{}) {
	//n := len(*p)
	item := x.(*worker)
	//item.SetId(n)
	*p = append(*p, item)
}

func (p *workerPool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	//item.SetId(-2) // for safety  // what safty? what's the danger there if id is not cleared off??
	*p = old[:n-1]
	return item
}
