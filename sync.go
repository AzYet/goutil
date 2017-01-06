package go_utils

import "sync"

func NewWorkerGroup(workerNum int) (*sync.WaitGroup, func() func()) {
	routine := make(chan int, workerNum)
	w := new(sync.WaitGroup)
	return w, func() func() {
		w.Add(1)
		routine <- 1
		return func() {
			<-routine
			w.Done()
		}
	}
}
