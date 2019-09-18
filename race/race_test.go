package race

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

var numRequestsToMake int
var numConcurrentRequests int

func init() {
	flag.IntVar(&numRequestsToMake, "total-requests", 1000, "total # of requests to make")
	flag.IntVar(&numConcurrentRequests, "concurrent-requests", 100, "pool size, request concurrency")
}

func TestExplicitRace(t *testing.T) {
	flag.Parse()

	reqCount := SynchronizedCounter{}

	go func() {
		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//value := reqCount.Value()
			//fmt.Printf("handling request: %d\n", value)
			time.Sleep(1 * time.Nanosecond)
			reqCount.Inc() //Set(value + 1)
			fmt.Fprintln(w, "Hello, client")
		}))
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	requestsChan := make(chan int)

	var wg sync.WaitGroup
	wg.Add(numConcurrentRequests)

	// start a pool of 100 workers all making requests
	for i := 0; i < numConcurrentRequests; i++ {
		go func() {
			defer wg.Done()
			for range requestsChan {
				res, err := http.Get("http://localhost:8080/")
				if err != nil {
					t.Fatal(err)
				}
				_, err = ioutil.ReadAll(res.Body)
				res.Body.Close()
				if err != nil {
					t.Error(err)
				}
			}
		}()
	}

	go func() {
		for i := 0; i < numRequestsToMake; i++ {
			requestsChan <- i
		}
		close(requestsChan)
	}()

	wg.Wait()

	fmt.Printf("Num Requests TO Make: %d\n", numRequestsToMake)
	fmt.Printf("Num Handled: %d\n", reqCount.Value())
	if numRequestsToMake != reqCount.Value() {
		t.Errorf("expected %d requests: received %d", numRequestsToMake, reqCount.Value())
	}
}

type safeIncrementer struct {
	m sync.Mutex
	x int
}

func (si *safeIncrementer) inc() {
	si.m.Lock()
	defer si.m.Unlock()
	si.x = si.x + 1
}
func (si *safeIncrementer) value() int {
	si.m.Lock()
	defer si.m.Unlock()
	return si.x
}
func (si *safeIncrementer) set(v int) {
	si.m.Lock()
	defer si.m.Unlock()
	si.x = v
}
func TestMutexNoRace(t *testing.T) {
	var wg sync.WaitGroup
	var incrementer safeIncrementer

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			incrementer.set(incrementer.value() + 1)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("final value of x", incrementer.x)
}

type channelIncrementer struct {
	x int
}

func (ci *channelIncrementer) inc() {
	ci.x = ci.x + 1
}
func TestChannelNoRace(t *testing.T) {
	var wg sync.WaitGroup
	var incrementer channelIncrementer

	chanIncrement := make(chan bool, 1)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			chanIncrement <- true
			incrementer.inc()
			<-chanIncrement
		}()
	}
	wg.Wait()

	fmt.Println("final value of x", incrementer.x)
}
