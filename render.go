package main

import (
	"fmt"
	"io"
)

type renderer struct {
	w io.Writer
}

func newRenderer(w io.Writer) renderer {
	return renderer{
		w: w,
	}
}

func (r *renderer) render(a ...any) {
	fmt.Fprint(r.w, append(a, "\r")...)
}
