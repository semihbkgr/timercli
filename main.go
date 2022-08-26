package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Print("\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b" + strconv(time.Now().Unix()))
	}
}
