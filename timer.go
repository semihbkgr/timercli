package main

import "time"

type Timer interface {
	Ticks() <-chan time.Duration
	Elapsed() time.Duration
	Stop()
	Proceed()
	Interrupt()
	Interrupted() bool
}

const defaultTickRate = 100 * time.Millisecond

type Countdown struct {
	duration        time.Duration
	tickRate        time.Duration
	ticker          *time.Ticker
	startTime       time.Time
	remaining       chan time.Duration
	stopped         bool
	stoppedAt       time.Time
	stoppedDuration time.Duration
	interrupted     bool
	interruptedAt   time.Time
}

func NewCountdown(d time.Duration) *Countdown {
	c := &Countdown{
		duration:        d,
		tickRate:        defaultTickRate,
		ticker:          time.NewTicker(defaultTickRate),
		startTime:       time.Now(),
		remaining:       make(chan time.Duration),
		stopped:         false,
		stoppedDuration: 0,
		interrupted:     false,
	}
	go startCountdown(c)
	return c
}

func (c *Countdown) Ticks() <-chan time.Duration {
	return c.remaining
}

func (c *Countdown) Elapsed() time.Duration {
	if c.interrupted {
		return time.Duration(c.interruptedAt.UnixNano() - c.stoppedDuration.Nanoseconds() - c.startTime.UnixNano())
	}
	if c.stopped {
		return time.Duration(c.stoppedAt.UnixNano() - c.stoppedDuration.Nanoseconds() - c.startTime.UnixNano())
	}
	return time.Duration(time.Now().UnixNano() - c.stoppedDuration.Nanoseconds() - c.startTime.UnixNano())
}

func (c *Countdown) Stop() {
	if !c.stopped && !c.interrupted {
		c.stoppedAt = time.Now()
		c.stopped = true
	}
}

func (c *Countdown) Proceed() {
	if c.stopped && !c.interrupted {
		c.stoppedDuration += time.Now().Sub(c.stoppedAt)
		c.stopped = false
	}
}

func (c *Countdown) Interrupt() {
	if !c.interrupted {
		c.interruptedAt = time.Now()
		c.stoppedDuration += time.Now().Sub(c.stoppedAt)
		c.interrupted = true
	}
}

func (c *Countdown) Interrupted() bool {
	return c.interrupted
}

func startCountdown(c *Countdown) {
	for {
		select {
		case <-c.ticker.C:
			r := c.duration - c.Elapsed()
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
	tickRate        time.Duration
	ticker          *time.Ticker
	startTime       time.Time
	elapsed         chan time.Duration
	stopped         bool
	stoppedAt       time.Time
	stoppedDuration time.Duration
	interrupted     bool
	interruptedAt   time.Time
}

func NewStopwatch() *Stopwatch {
	s := &Stopwatch{
		tickRate:        defaultTickRate,
		ticker:          time.NewTicker(defaultTickRate),
		startTime:       time.Now(),
		elapsed:         make(chan time.Duration),
		stopped:         false,
		stoppedDuration: 0,
		interrupted:     false,
	}
	go startStopwatch(s)
	return s
}

func (s *Stopwatch) Ticks() <-chan time.Duration {
	return s.elapsed
}

func (s *Stopwatch) Elapsed() time.Duration {
	if s.interrupted {
		return time.Duration(s.interruptedAt.UnixNano() - s.stoppedDuration.Nanoseconds() - s.startTime.UnixNano())
	}
	if s.stopped {
		return time.Duration(s.stoppedAt.UnixNano() - s.stoppedDuration.Nanoseconds() - s.startTime.UnixNano())
	}
	return time.Duration(time.Now().UnixNano() - s.stoppedDuration.Nanoseconds() - s.startTime.UnixNano())
}

func (s *Stopwatch) Stop() {
	if !s.stopped && !s.interrupted {
		s.stoppedAt = time.Now()
		s.stopped = true
	}
}

func (s *Stopwatch) Proceed() {
	if s.stopped && !s.interrupted {
		s.stoppedDuration += time.Now().Sub(s.stoppedAt)
		s.stopped = false
	}
}

func (s *Stopwatch) Interrupt() {
	if !s.interrupted {
		s.interruptedAt = time.Now()
		s.stoppedDuration += time.Now().Sub(s.stoppedAt)
		s.interrupted = true
	}
}

func (s *Stopwatch) Interrupted() bool {
	return s.interrupted
}

func startStopwatch(s *Stopwatch) {
	for {
		select {
		case <-s.ticker.C:
			e := s.Elapsed()
			select {
			case s.elapsed <- e:
			case <-s.elapsed:
				s.elapsed <- e
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
	close(s.elapsed)
	s.ticker.Stop()
}
