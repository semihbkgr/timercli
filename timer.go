package main

import "time"

const defaultTickRate = 100 * time.Millisecond

type Timer struct {
	duration        time.Duration
	tickRate        time.Duration
	ticker          *time.Ticker
	startedAt       time.Time
	ticks           chan time.Duration
	stopped         bool
	stoppedAt       time.Time
	stoppedDuration time.Duration
	interrupted     bool
	interruptedAt   time.Time
}

func (t *Timer) Ticks() <-chan time.Duration {
	return t.ticks
}

func (t *Timer) Elapsed() time.Duration {
	if t.interrupted {
		return time.Duration(t.interruptedAt.UnixNano() - t.stoppedDuration.Nanoseconds() - t.startedAt.UnixNano())
	}
	if t.stopped {
		return time.Duration(t.stoppedAt.UnixNano() - t.stoppedDuration.Nanoseconds() - t.startedAt.UnixNano())
	}
	return time.Duration(time.Now().UnixNano() - t.stoppedDuration.Nanoseconds() - t.startedAt.UnixNano())
}

func (t *Timer) StartedAt() time.Time {
	return t.startedAt
}

func (t *Timer) Stop() {
	if t.Running() {
		t.stoppedAt = time.Now()
		t.stopped = true
	}
}

func (t *Timer) Proceed() {
	if t.stopped && !t.interrupted {
		t.stoppedDuration += time.Now().Sub(t.stoppedAt)
		t.stopped = false
	}
}

func (t *Timer) Running() bool {
	return !t.stopped && !t.interrupted
}

func (t *Timer) Interrupt() {
	if !t.interrupted {
		t.interruptedAt = time.Now()
		if t.stopped {
			t.stoppedDuration += time.Now().Sub(t.stoppedAt)
		}
		t.interrupted = true
	}
}

func (t *Timer) Interrupted() bool {
	return t.interrupted
}

func NewCountdown(d time.Duration) *Timer {
	c := &Timer{
		duration:        d,
		tickRate:        defaultTickRate,
		ticker:          time.NewTicker(defaultTickRate),
		startedAt:       time.Now(),
		ticks:           make(chan time.Duration),
		stopped:         false,
		stoppedDuration: 0,
		interrupted:     false,
	}
	go startCountdown(c)
	return c
}

func NewStopwatch() *Timer {
	s := &Timer{
		tickRate:        defaultTickRate,
		ticker:          time.NewTicker(defaultTickRate),
		startedAt:       time.Now(),
		ticks:           make(chan time.Duration),
		stopped:         false,
		stoppedDuration: 0,
		interrupted:     false,
	}
	go startStopwatch(s)
	return s
}

func startCountdown(t *Timer) {
	for {
		select {
		case <-t.ticker.C:
			r := t.duration - t.Elapsed()
			if r < 0 {
				r = 0
			}
			select {
			case t.ticks <- r:
			case <-t.ticks:
				t.ticks <- r
			}
			if r == 0 {
				interruptTimer(t)
				return
			}
		default:
			if t.interrupted {
				interruptTimer(t)
				return
			}
		}
	}
}

func startStopwatch(t *Timer) {
	for {
		select {
		case <-t.ticker.C:
			e := t.Elapsed()
			select {
			case t.ticks <- e:
			case <-t.ticks:
				t.ticks <- e
			}
		default:
			if t.interrupted {
				interruptTimer(t)
				return
			}
		}
	}
}

func interruptTimer(t *Timer) {
	close(t.ticks)
	t.ticker.Stop()
}
