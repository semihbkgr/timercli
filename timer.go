package main

import "time"

type Timer interface {
	Remaining() <-chan time.Duration
	Interrupt()
}

const defaultTickRate = 100 * time.Millisecond

type Countdown struct {
	duration    time.Duration
	tickRate    time.Duration
	ticker      *time.Ticker
	startTime   time.Time
	remaining   chan time.Duration
	interrupted bool
}

func NewCountdown(d time.Duration) *Countdown {
	c := &Countdown{
		duration:    d,
		tickRate:    defaultTickRate,
		ticker:      time.NewTicker(defaultTickRate),
		startTime:   time.Now(),
		remaining:   make(chan time.Duration),
		interrupted: false,
	}
	go startCountdown(c)
	return c
}

func (c *Countdown) Remaining() <-chan time.Duration {
	return c.remaining
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
			select {
			case c.remaining <- r:
			case <-c.remaining:
				c.remaining <- r
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
	close(c.remaining)
	c.ticker.Stop()
}

type Chronometer struct {
	tickRate    time.Duration
	ticker      *time.Ticker
	startTime   time.Time
	remaining   chan time.Duration
	interrupted bool
}

func NewChronometer() *Chronometer {
	c := &Chronometer{
		tickRate:    defaultTickRate,
		ticker:      time.NewTicker(defaultTickRate),
		startTime:   time.Now(),
		remaining:   make(chan time.Duration),
		interrupted: false,
	}
	go startChronometer(c)
	return c
}

func (c *Chronometer) Remaining() <-chan time.Duration {
	return c.remaining
}

func (c *Chronometer) Interrupt() {
	c.interrupted = true
}

func startChronometer(c *Chronometer) {
	for {
		select {
		case t := <-c.ticker.C:
			r := time.Duration(t.UnixNano() - c.startTime.UnixNano())
			select {
			case c.remaining <- r:
			case <-c.remaining:
				c.remaining <- r
			}
		default:
			if c.interrupted {
				interruptChronometer(c)
				return
			}
		}
	}
}

func interruptChronometer(c *Chronometer) {
	close(c.remaining)
	c.ticker.Stop()
}
