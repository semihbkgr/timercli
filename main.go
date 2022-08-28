package main

import (
	"os"
	"time"
)

func main() {

	time.Sleep(time.Second)

	r := newRenderer(os.Stdout)

	for {
		r.render(time.Now().Format(time.RFC3339), "\n")
	}

}
