package main

import (
	"os"
	"time"
)

func main() {
	startTimer()
}

func startTimer() {
	st := time.Now().UnixMilli()
	r := newRenderer(os.Stdout)
	for {
		r.render(time.Now().UnixMilli() - st)
		time.Sleep(10 * time.Millisecond)
	}
}
