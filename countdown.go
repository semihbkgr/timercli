package main

import "time"

const defaultTickRate = 100 * time.Millisecond

type countdown struct {
	duration  time.Duration
	tickRate  time.Duration
	ticker    *time.Ticker
	startTime time.Time
	chans     []chan time.Duration
}

func newCountdown(d time.Duration) *countdown {
	c := &countdown{
		duration:  d,
		tickRate:  defaultTickRate,
		ticker:    time.NewTicker(defaultTickRate),
		startTime: time.Now(),
		chans:     make([]chan time.Duration, 0),
	}
	go startCountdown(c)
	return c
}

func (c *countdown) remaining() chan<- time.Duration {
	ch := make(chan time.Duration, 0)
	c.chans = append(c.chans, ch)
	return ch
}

func startCountdown(c *countdown) {
	for {
		t, ok := <-c.ticker.C
		if !ok {
			for _, ch := range c.chans {
				close(ch)
			}
			return
		}
		r := c.duration - time.Duration(t.UnixNano()-c.startTime.UnixNano())
		if r < 0 {
			c.ticker.Stop()
			continue
		}
		for _, ch := range c.chans {
			ch <- r
		}
	}
}
