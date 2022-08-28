package main

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type renderer struct {
	w io.Writer
	l int
}

func newRenderer(w io.Writer) renderer {
	return renderer{
		w: w,
		l: 0,
	}
}

func (r *renderer) render(a ...any) {
	b := strings.Builder{}
	if r.l > 0 {
		b.WriteString(strings.Repeat("\b", r.l))
	}
	s := fmt.Sprint(a...)
	r.l = len(s)
	b.WriteString(s)
	r.w.Write([]byte(b.String()))
	time.Sleep(time.Millisecond)
}
