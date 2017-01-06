package go_utils

import (
	"time"
	"github.com/Sirupsen/logrus"
)

type Ticker struct {
	C <-chan time.Time
	i time.Duration
	d time.Duration
	s chan int
}

func (t *Ticker)Stop() {
	t.s <- -1
}
func (t *Ticker)Tick() {
	t.s <- 1
}
func (t *Ticker)Update(interval, delay time.Duration) {
	t.i = interval
	t.d = delay
}

func NewTicker(interval time.Duration, delay time.Duration, Logger *logrus.Logger) *Ticker {
	tc := make(chan time.Time)
	tms := &Ticker{tc, interval, delay, make(chan int, 1) }
	ts := time.Now()
	go func() {
		if sub := ts.Sub(ts.Truncate(tms.i)) - tms.d; sub >= 0 {
			//Logger.Printf("exceed %v, start now.", ts.Truncate(tms.i).Add(tms.d))
			tc <- ts
			time.Sleep(ts.Add(tms.i).Truncate(tms.i).Add(tms.d).Sub(ts))
			tc <- time.Now()
		} else {
			Logger.Infof("first task will start at %v.", ts.Truncate(tms.i).Add(tms.d).Format("2006-01-02 15:04:05"))
			time.Sleep(ts.Truncate(tms.i).Add(tms.d).Sub(ts))
			tc <- time.Now()
		}
		for {
			ts = time.Now()
			select {
			case <-time.NewTimer(ts.Truncate(tms.i).Add(tms.i).Add(tms.d).Sub(ts)).C:
				tc <- time.Now()
			case i := <-tms.s:
				if i < 0 {
					return
				} else {
					select {
					case tc <- time.Now():
					default:
					}
				}

			}
		}
	}()
	return tms
}

