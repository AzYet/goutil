package go_utils

import (
	"time"
	"github.com/Sirupsen/logrus"
)

type TickMinSec struct {
	C <-chan time.Time
	s chan int
	t chan int
}

func (t TickMinSec)Stop() {
	t.s <- 1
}
func (t TickMinSec)Tick() {
	t.t <- 1
}

func TickAtSecondPoint(interval time.Duration, delay time.Duration, Logger *logrus.Logger) TickMinSec {
	tc := make(chan time.Time)
	tms := TickMinSec{tc, make(chan int, 1), make(chan int, 1) }
	ts := time.Now()
	go func() {
		if sub := ts.Sub(ts.Truncate(interval)) - delay; sub >= 0 {
			//Logger.Printf("exceed %v, start now.", ts.Truncate(interval).Add(delay))
			tc <- ts
			time.Sleep(ts.Add(interval).Truncate(interval).Add(delay).Sub(ts))
			tc <- time.Now()
		} else {
			Logger.Printf("first task will start at %v.", ts.Truncate(interval).Add(delay).Format("2006-01-02 15:04:05"))
			time.Sleep(ts.Truncate(interval).Add(delay).Sub(ts))
			tc <- time.Now()
		}
		for {
			ts = time.Now()
			select {
			case <-time.NewTimer(ts.Truncate(interval).Add(interval).Add(delay).Sub(ts)).C:
				tc <- time.Now()
			case <-tms.t:
				tc <- time.Now()
			case <-tms.s:
				return
			}
		}
	}()
	return tms
}

