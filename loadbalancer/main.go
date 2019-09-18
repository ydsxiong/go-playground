package main

/**

|Client1|         |Load|  <-DONE-- |Worker1| processing R3 coming from Client1
|Client2| --REQ-> |Blncr| --WOK->  |Worker2| processing R8 coming from Client2
                                   |Worker3| processing R5 coming from Client1
						  <-RESP-- response of R4 to Client2


Data Flow
k Clients pack the value x in Request object and sends it to REQ channel.
Load balancer blocks on REQ channel listening to Request(s).
Load balancer chooses a worker and sends Request to one of the channels of worker WOK(i).
Worker receives Request and processes x (say calculates sin(x) lol).
Worker updates load balancer using DONE channel. LB uses this for load-balancing.
Worker writes the sin(x) value in the RESP response channel (enclosed in Request object).


Channels in play
central REQ channel (Type: Request)
n WOK channels (n sized worker pool, Type: Work)
k RESP channels (k clients, Type: Float)
n DONE channels (Type: Work)

Glueing the code
Imports and main
Adding a print

Output
Here you can see number of pending tasks per worker.
Since work is just computing a sine value, I had to reduce sleep-time at the client level before they fire next request.
0 1 2 3 4 5 6 7 8 9  avg  variance

5 6 8 8 8 8 8 8 8 8  7.50 1.05
4 6 8 8 8 8 8 8 8 8  7.40 1.64
3 6 8 8 8 8 8 8 8 8  7.30 2.41
2 6 8 8 8 8 8 8 8 8  7.20 3.36
1 6 8 8 8 8 8 8 8 8  7.10 4.49
1 5 8 8 8 8 8 8 8 8  7.00 4.80
1 5 8 8 7 8 8 8 8 8  6.90 4.69
1 5 8 8 6 8 8 8 8 8  6.80 4.76
1 4 8 8 6 8 8 8 8 8  6.70 5.21
1 4 8 8 6 8 8 8 8 7  6.60 5.04
1 4 8 7 6 8 8 8 8 7  6.50 4.85
1 4 8 7 6 8 8 8 7 7  6.40 4.64
1 4 7 7 6 8 8 8 7 7  6.30 4.41
Footnote
Although this is still a single-process LB, it makes you appreciate the flexibility of asynchronous behavior.
How channels communicate in form of a light-weight queue and offloading the tasks in form of goroutines is pretty amazing to me.
Also, all the book-keeping of acquiring/releasing a lock is hidden from programmer and all you need to focus on "sharing data using channels and not the variables" ;).
Event driven architecture is amazing concept. There's a nice writeup on event-driven architecture I read on hackernews the other day that tells you when and when not to use it: https://herbertograca.com/2017/10/05/event-driven-architecture/


*/

import (
	"github.com/ydsxiong/go-playground/loadbalancer/balancer"
)

const (
	nRequester = 10
	nWorker    = 4
)

func main() {

	work := make(chan balancer.Request)
	for i := 0; i < nRequester; i++ {
		go balancer.CreateAndRequest(work)
	}

	balancer.InitBalancer(nWorker, nRequester).Balance(work)
}
