package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

//todo: stop operation
//todo: print
//todo: better format
//todo: send os notification of possible

func main() {
	d, err := parseDuration()
	checkErr(err)
	if d != 0 { // start countdown
		checkErr(validateDuration(d))
		r := NewCountdown(d).Remaining()
		consume(r, func(d time.Duration) {
			fmt.Printf("\r%s", d)
		})
	} else { // start chronometer
		r := NewChronometer().Remaining()
		consume(r, func(d time.Duration) {
			fmt.Printf("\r%s", d)
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
