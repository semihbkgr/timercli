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

var theme = commandLine.String("t", "light", "theme of the renderer on the console")

func main() {
	defer handleError()
	d, ok := parseDuration()
	var t Timer
	if ok {
		t = NewCountdown(d)
	} else {
		t = NewStopwatch()
	}
	defer printElapsedTime(t)
	initTermbox()
	defer closeTermbox()
	handleCtrlCInput(func() {
		t.Interrupt()
	})
	r := NewRenderer(*theme)
	for d := range t.Ticks() {
		err := r.Render(formatDuration(d))
		checkErr(err)
	}
	if !t.Interrupted() {
		fmt.Print("\a")
	}
}

func commandLineArgs() ([]string, error) {
	if !commandLine.Parsed() {
		err := commandLine.Parse(os.Args[1:])
		if err != nil {
			flagParseError = err != flag.ErrHelp
			return nil, err
		}
	}
	return commandLine.Args(), nil
}

func parseDuration() (time.Duration, bool) {
	args, err := commandLineArgs()
	if err != nil {
		panic(err)
	}
	if l := len(args); l == 0 {
		return 0, false
	} else if l > 1 {
		panic(errors.New("too many arguments"))
	}
	d, err := time.ParseDuration(args[0])
	if err != nil {
		panic(err)
	}
	if d < 1 {
		panic(errors.New("non-positive duration"))
	}
	return d, true
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func handleError() {
	r := recover()
	status := 0
	if r != nil {
		switch t := r.(type) {
		case error:
			if !flagParseError && t != flag.ErrHelp {
				_, err := fmt.Fprintln(os.Stderr, t.Error())
				if err != nil {
					panic(err)
				}
				status = 1
			}
		default:
			_, err := fmt.Fprintln(os.Stderr, t)
			if err != nil {
				panic(err)
			}
			status = 1
		}
	}
	os.Exit(status)
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h == 0 {
		return fmt.Sprintf("%02d:%02d", m, s)
	} else {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
}

func printElapsedTime(t Timer) {
	fmt.Printf("Elapsed time: %s\n", formatDuration(t.Elapsed()))
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
