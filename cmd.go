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
	c := newCountdown(d)
	c.start()
	c.wait()
}

func parseDuration() (time.Duration, error) {
	args := flag.Args()
	if len(args) == 0 {
		return 0, errors.New("missing duration arg")
	}
	return time.ParseDuration(args[0])
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
