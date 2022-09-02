package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	d, err := parseDuration()
	checkErr(err)
	c := NewCountdown(d)
	r := c.Remaining()
	consume(r, func(d time.Duration) {
		fmt.Printf("/r%s", d)
	})
}

func parseDuration() (time.Duration, error) {
	if !flag.Parsed() {
		flag.Parse()
	}
	args := flag.Args()
	if len(args) == 0 {
		return 0, errors.New("missing duration arg")
	}
	d, err := time.ParseDuration(args[0])
	if err != nil {
		return 0, err
	}
	return validateDuration(d)
}

func validateDuration(d time.Duration) (time.Duration, error) {
	if d < 0 {
		return 0, errors.New("non-positive duration")
	}
	return d, nil
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
