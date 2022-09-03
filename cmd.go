package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

//todo: stop operation
//todo: better format
//todo: send os notification of possible
//todo: handle os signals

func main() {
	initTermbox()
	defer closeTermbox()
	d, err := parseDuration()
	checkErr(err)
	if d != 0 { // start countdown
		checkErr(validateDuration(d))
		c := NewCountdown(d)
		//r := NewRenderer("light")
		consume(c.Remaining(), func(d time.Duration) {
			//todo: format duration
		})
	} else { // start chronometer
		c := NewChronometer()
		//r := NewRenderer("light")
		consume(c.Remaining(), func(d time.Duration) {
			//todo: format duration
		})
	}
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
