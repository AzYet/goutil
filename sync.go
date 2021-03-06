package goutil

import (
	"sync"
)

func NewWorkerGroup(workerNum int) (chan struct{}, *sync.WaitGroup, func() func()) {
	routine := make(chan struct{}, workerNum)
	w := new(sync.WaitGroup)
	return routine, w, func() func() {
		routine <- struct{}{}
		w.Add(1)
		return func() {
			<-routine
			w.Done()
		}
	}
}
func NewWorkerGroupWithReturn(workerNum int) (chan interface{}, func() func(res interface{})) {
	routine, c := make(chan struct{}, workerNum), make(chan interface{}, workerNum)
	return c, func() func(interface{}) {
		routine <- struct{}{}
		return func(r interface{}) {
			c <- r
			<-routine
		}
	}
}

//Create a fixed size pool of recycled resource, call the return func with nil to get resource, call with resource to return to pool,
// when pool exhausts, getting resource causes block
//make sure return resource when done, do not return nil to pool, as pool will not check for nil, unless you intent to
func NewResPool(size int, newRes func() (interface{}, error)) func(interface{}) interface{} {
	pool := make(chan interface{}, size)
	for i := 0; i < size; i++ {
		pool <- nil
	}
	getc, retc, out := make(chan interface{}), make(chan interface{}), make(chan interface{})
	go func() {
		for range getc {
			//get
			i := <-pool
			if i == nil {
				//new
				for {
					if n, e := newRes(); e == nil {
						out <- n
						break
					}
				}
			} else {
				// cycling
				out <- i
			}
		}
	}()
	go func() {
		for res := range retc {
			pool <- res
		}
	}()
	return func(in interface{}) interface{} {
		if in == nil {
			getc <- nil
			return <-out
		} else {
			retc <- in
			return nil
		}
	}
}
