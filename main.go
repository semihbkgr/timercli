package main

import (
	"fmt"
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
		r.render(fmt.Sprintf("%s", time.Duration(time.Now().UnixMilli()-st)*time.Millisecond))
		time.Sleep(100 * time.Millisecond)
	}
}
