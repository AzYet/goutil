package goutil

import (
	"time"
)

type Ticker struct {
	C <-chan time.Time
	I time.Duration
	D time.Duration
	s chan int
}

func (t *Ticker) Stop() {
	t.s <- -1
}
func (t *Ticker) Tick() {
	t.s <- 1
}
func (t *Ticker) Update(interval, delay time.Duration) {
	t.I = interval
	t.D = delay
}

func NewTickerWithBuffer(interval time.Duration, delay time.Duration, bufferSize int) *Ticker {
	tc := make(chan time.Time, bufferSize)
	tms := &Ticker{tc, interval, delay, make(chan int, 1)}
	ts := time.Now()
	go func() {
		if sub := ts.Sub(ts.Truncate(tms.I)) - tms.D; sub >= 0 {
			//fmt.Printf("exceed %v, start now.", ts.Truncate(tms.I).Add(tms.D))
			tc <- ts
			time.Sleep(ts.Add(tms.I).Truncate(tms.I).Add(tms.D).Sub(ts))
			tc <- time.Now()
		} else {
			time.Sleep(ts.Truncate(tms.I).Add(tms.D).Sub(ts))
			tc <- time.Now()
		}
		for {
			ts = time.Now()
			select {
			case now := <-time.After(ts.Truncate(tms.I).Add(tms.I + tms.D).Sub(ts)):
				select {
				case <-tc:
				default:
				}
				tc <- now
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

func NewTicker(interval time.Duration, delay time.Duration) *Ticker {
	return NewTickerWithBuffer(interval, delay, 0)
}
