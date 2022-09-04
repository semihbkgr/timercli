package main

import "time"

type Chronometer struct {
	tickRate  time.Duration
	ticker    *time.Ticker
	startTime time.Time
	chans     []chan time.Duration
	stop      bool
}

func NewChronometer() *Chronometer {
	c := &Chronometer{
		tickRate:  defaultTickRate,
		ticker:    time.NewTicker(defaultTickRate),
		startTime: time.Now(),
		chans:     make([]chan time.Duration, 0),
		stop:      false,
	}
	go startChronometer(c)
	return c
}

func (c *Chronometer) Remaining() <-chan time.Duration {
	ch := make(chan time.Duration, 0)
	c.chans = append(c.chans, ch)
	return ch
}

func (c *Chronometer) Stop() {
	c.stop = true
}

func startChronometer(c *Chronometer) {
	for {
		select {
		case t := <-c.ticker.C:
			r := time.Duration(t.UnixNano() - c.startTime.UnixNano())
			for _, ch := range c.chans {
				select {
				case ch <- r:
				}
			}
		default:
			if c.stop {
				stopChronometer(c)
				return
			}
		}
	}
}

func stopChronometer(c *Chronometer) {
	c.ticker.Stop()
	for _, ch := range c.chans {
		close(ch)
	}
}
