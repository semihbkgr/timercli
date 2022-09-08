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
				select {
				case ch <- r:
				default:
				}
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

type Chronometer struct {
	tickRate    time.Duration
	ticker      *time.Ticker
	startTime   time.Time
	channels    []chan time.Duration
	interrupted bool
}

func NewChronometer() *Chronometer {
	c := &Chronometer{
		tickRate:    defaultTickRate,
		ticker:      time.NewTicker(defaultTickRate),
		startTime:   time.Now(),
		channels:    make([]chan time.Duration, 0),
		interrupted: false,
	}
	go startChronometer(c)
	return c
}

func (c *Chronometer) Remaining() <-chan time.Duration {
	ch := make(chan time.Duration, 0)
	c.channels = append(c.channels, ch)
	return ch
}

func (c *Chronometer) Interrupt() {
	c.interrupted = true
}

func startChronometer(c *Chronometer) {
	for {
		select {
		case t := <-c.ticker.C:
			r := time.Duration(t.UnixNano() - c.startTime.UnixNano())
			for _, ch := range c.channels {
				select {
				case ch <- r:
				default:
				}
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
	c.ticker.Stop()
	for _, ch := range c.channels {
		close(ch)
	}
}
