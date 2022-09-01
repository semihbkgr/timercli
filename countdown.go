package main

import "time"

type countdown struct {
	duration time.Duration
	tickRate time.Duration
	ticker   *time.Ticker
}

func newCountdown(d time.Duration) *countdown {
	return &countdown{
		duration: d,
		tickRate: 100 * time.Millisecond,
		ticker:   time.NewTicker(d),
	}
}

func (c *countdown) start() {
	go func() {

	}()
}

func (c *countdown) wait() {
	<-c.ticker.C
}
