package main

import "time"

type Chronometer struct {
	tickRate  time.Duration
	ticker    *time.Ticker
	startTime time.Time
	chans     []chan time.Duration
}

func NewChronometer() *Chronometer {
	c := &Chronometer{
		tickRate:  defaultTickRate,
		ticker:    time.NewTicker(defaultTickRate),
		startTime: time.Now(),
		chans:     make([]chan time.Duration, 0),
	}
	go startChronometer(c)
	return c
}

func (c *Chronometer) Remaining() <-chan time.Duration {
	ch := make(chan time.Duration, 0)
	c.chans = append(c.chans, ch)
	return ch
}

func startChronometer(c *Chronometer) {
	for {
		t := <-c.ticker.C
		r := time.Duration(t.UnixNano() - c.startTime.UnixNano())
		for _, ch := range c.chans {
			ch <- r
		}
	}
}
