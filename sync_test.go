package goutil

import (
	"testing"
	"fmt"
	"time"
)

func TestNewResPool(t *testing.T) {
	count := 0
	op := NewResPool(5, func() (interface{}, error) {
		s := make([]int, 4)
		s[0] = count
		count++
		return s, nil
	})
	for {
		i := op(nil)
		fmt.Println(i)
		go func() {
			time.Sleep(time.Second * 10)
			op(i)
		}()
		time.Sleep(time.Second)
	}
	fmt.Println("loop over") // this should not be seen
}

func TestNewWorkerGroupWithReturn(t *testing.T) {
	_, c, reg := NewWorkerGroupWithReturn(2)
	go func() {
		for i := 0; i < 12; i++ {
			go func(i int, f func(r interface{})) {
				time.Sleep(time.Second)
				f(i)
			}(i, reg())
		}
	}()
	//go func() {
	for i := 0; i < 12; i++ {
		r := <-c
		fmt.Println(r)
	}
	//}()
}