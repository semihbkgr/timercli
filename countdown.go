package main

import "time"

const defaultTickRate = 100 * time.Millisecond

type Countdown struct {
	duration  time.Duration
	tickRate  time.Duration
	ticker    *time.Ticker
	startTime time.Time
	chans     []chan time.Duration
	stop      bool
}

func NewCountdown(d time.Duration) *Countdown {
	c := &Countdown{
		duration:  d,
		tickRate:  defaultTickRate,
		ticker:    time.NewTicker(defaultTickRate),
		startTime: time.Now(),
		chans:     make([]chan time.Duration, 0),
		stop:      false,
	}
	go startCountdown(c)
	return c
}

func (c *Countdown) Remaining() <-chan time.Duration {
	ch := make(chan time.Duration, 0)
	c.chans = append(c.chans, ch)
	return ch
}

func (c *Countdown) Stop() {
	c.stop = true
}

func startCountdown(c *Countdown) {
	for {
		select {
		case t := <-c.ticker.C:
			r := c.duration - time.Duration(t.UnixNano()-c.startTime.UnixNano())
			if r < 0 {
				r = 0
			}
			for _, ch := range c.chans {
				ch <- r
			}
			if r == 0 {
				stopCountdown(c)
				return
			}
		default:
			if c.stop {
				stopCountdown(c)
				return
			}
		}
	}
}

func stopCountdown(c *Countdown) {
	c.ticker.Stop()
	for _, ch := range c.chans {
		close(ch)
	}
}
