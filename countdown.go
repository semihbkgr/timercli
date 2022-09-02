package main

import "time"

const defaultTickRate = 100 * time.Millisecond

type Countdown struct {
	duration  time.Duration
	tickRate  time.Duration
	ticker    *time.Ticker
	startTime time.Time
	chans     []chan time.Duration
}

func NewCountdown(d time.Duration) *Countdown {
	c := &Countdown{
		duration:  d,
		tickRate:  defaultTickRate,
		ticker:    time.NewTicker(defaultTickRate),
		startTime: time.Now(),
		chans:     make([]chan time.Duration, 0),
	}
	go startCountdown(c)
	return c
}

func (c *Countdown) Remaining() <-chan time.Duration {
	ch := make(chan time.Duration, 0)
	c.chans = append(c.chans, ch)
	return ch
}

func startCountdown(c *Countdown) {
	for {
		t := <-c.ticker.C
		r := c.duration - time.Duration(t.UnixNano()-c.startTime.UnixNano())
		if r < 0 {
			r = 0
		}
		for _, ch := range c.chans {
			ch <- r
		}
		if r == 0 {
			c.ticker.Stop()
			for _, ch := range c.chans {
				close(ch)
			}
			return
		}
	}
}
