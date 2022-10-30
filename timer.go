package main

import "time"

type Timer interface {
	Ticks() <-chan time.Duration
	Elapsed() time.Duration
	Interrupt()
	Interrupted() bool
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

func (c *Countdown) Ticks() <-chan time.Duration {
	return c.remaining
}

func (c *Countdown) Elapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - c.startTime.UnixNano())
}

func (c *Countdown) Interrupt() {
	c.interrupted = true
}

func (c *Countdown) Interrupted() bool {
	return c.interrupted
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

type Stopwatch struct {
	tickRate    time.Duration
	ticker      *time.Ticker
	startTime   time.Time
	remaining   chan time.Duration
	interrupted bool
}

func NewStopwatch() *Stopwatch {
	s := &Stopwatch{
		tickRate:    defaultTickRate,
		ticker:      time.NewTicker(defaultTickRate),
		startTime:   time.Now(),
		remaining:   make(chan time.Duration),
		interrupted: false,
	}
	go startStopwatch(s)
	return s
}

func (s *Stopwatch) Ticks() <-chan time.Duration {
	return s.remaining
}

func (s *Stopwatch) Elapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - s.startTime.UnixNano())
}

func (s *Stopwatch) Interrupt() {
	s.interrupted = true
}

func (s *Stopwatch) Interrupted() bool {
	return s.interrupted
}

func startStopwatch(s *Stopwatch) {
	for {
		select {
		case t := <-s.ticker.C:
			r := time.Duration(t.UnixNano() - s.startTime.UnixNano())
			select {
			case s.remaining <- r:
			case <-s.remaining:
				s.remaining <- r
			}
		default:
			if s.interrupted {
				interruptStopwatch(s)
				return
			}
		}
	}
}

func interruptStopwatch(s *Stopwatch) {
	close(s.remaining)
	s.ticker.Stop()
}
