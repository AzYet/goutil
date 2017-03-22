package go_utils

import (
	"testing"
	"fmt"
	"time"
)

func TestNewResPool(t *testing.T) {
	count := 0
	get, ret := NewResPool(5, func() interface{} {
		s := make([]int, 4)
		s[0] = count
		count++
		return s
	})
	for {
		i := get()
		fmt.Println(i)
		go func() {
			time.Sleep(time.Second * 10)
			ret(i)
		}()
		time.Sleep(time.Second)
	}
	fmt.Println("loop over") // this should not be seen
}
