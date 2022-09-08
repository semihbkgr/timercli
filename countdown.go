package main

import "time"

const defaultTickRate = 100 * time.Millisecond

type Countdown struct {
	duration    time.Duration
	tickRate    time.Duration
	ticker      *time.Ticker
	startTime   time.Time
	channels    []chan time.Duration
	interrupted bool
}

func NewCountdown(d time.Duration) *Countdown {
	c := &Countdown{
		duration:    d,
		tickRate:    defaultTickRate,
		ticker:      time.NewTicker(defaultTickRate),
		startTime:   time.Now(),
		channels:    make([]chan time.Duration, 0),
		interrupted: false,
	}
	go startCountdown(c)
	return c
}

func (c *Countdown) Remaining() <-chan time.Duration {
	ch := make(chan time.Duration, 0)
	c.channels = append(c.channels, ch)
	return ch
}

func (c *Countdown) Interrupt() {
	c.interrupted = true
}

func startCountdown(c *Countdown) {
	for {
		select {
		case t := <-c.ticker.C:
			r := c.duration - time.Duration(t.UnixNano()-c.startTime.UnixNano())
			if r < 0 {
				r = 0
			}
			for _, ch := range c.channels {
				ch <- r
			}
			if r == 0 {
				interruptCountdown(c)
				return
			}
		default:
			if c.interrupted {
				interruptCountdown(c)
				return
			}
		}
	}
}

func interruptCountdown(c *Countdown) {
	c.ticker.Stop()
	for _, ch := range c.channels {
		close(ch)
	}
}
