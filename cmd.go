package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"time"
)

var theme = flag.String("t", "", "theme")

//todo: rename stop to interrupt
//todo: stop and proceed signal
//todo: refactor timers to single struct
//todo: print elapsed time at the end
func main() {
	initTermbox()
	d, err := parseDuration()
	checkErr(err)
	if d != 0 { // start countdown
		checkErr(validateDuration(d))
		c := NewCountdown(d)
		handleCtrlCInput(func() {
			c.Stop()
		})
		r := NewRenderer(*theme)
		consume(c.Remaining(), func(d time.Duration) {
			err := r.Render(formatDuration(d))
			checkErr(err)
		})
	} else { // start chronometer
		c := NewChronometer()
		handleCtrlCInput(func() {
			c.Stop()
		})
		r := NewRenderer(*theme)
		consume(c.Remaining(), func(d time.Duration) {
			err := r.Render(formatDuration(d))
			checkErr(err)
		})
	}
	closeTermbox()
}

func parseDuration() (time.Duration, error) {
	if !flag.Parsed() {
		flag.Parse()
	}
	args := flag.Args()
	if len(args) == 0 {
		return 0, nil
	}
	return time.ParseDuration(args[0])
}

func validateDuration(d time.Duration) error {
	if d < 0 {
		return errors.New("non-positive duration")
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func consume[T any](c <-chan T, f func(T)) {
	for d := range c {
		f(d)
	}
}

func formatDuration(d time.Duration) string {
	m := int(d.Minutes())
	s := int(d.Seconds())
	f := fmt.Sprintf("%02d : %02d", m, s)
	return f
}

func handleCtrlCInput(f func()) {
	go func() {
		termbox.SetInputMode(termbox.InputEsc)
		for {
			e := termbox.PollEvent()
			switch e.Key {
			case termbox.KeyCtrlC:
				f()
				return
			}
		}
	}()
}
