package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"time"
)

var commandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
var flagParseError = false

// flags
var theme = commandLine.String("t", "light", "theme of the renderer on the console")

//todo: interrupted and proceed signal
//todo: print elapsed time at the end
//todo: error when duration is 0
func main() {
	defer handleError()
	initTermbox()
	defer closeTermbox()
	d, err := parseDuration()
	checkErr(err)
	var t Timer
	if d != 0 {
		checkErr(validateDuration(d))
		t = NewCountdown(d)
	} else {
		t = NewChronometer()
	}
	handleCtrlCInput(func() {
		t.Interrupt()
	})
	r := NewRenderer(*theme)
	c := t.Remaining()
	for d := range c {
		err := r.Render(formatDuration(d))
		checkErr(err)
	}
}

func commandLineArgs() []string {
	if !commandLine.Parsed() {
		err := commandLine.Parse(os.Args[1:])
		if err != nil {
			flagParseError = true
			panic(err)
		}
	}
	return commandLine.Args()
}

func parseDuration() (time.Duration, error) {
	args := commandLineArgs()
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
		panic(err)
	}
}

func handleError() {
	r := recover()
	if r != nil {
		switch t := r.(type) {
		case error:
			if flagParseError {
				commandLine.Usage()
			}
			if t != flag.ErrHelp {
				_, err := fmt.Fprintln(os.Stderr, t.Error())
				if err != nil {
					panic(err)
				}
			}
		default:
			_, err := fmt.Fprintln(os.Stderr, t)
			if err != nil {
				panic(err)
			}
		}
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
