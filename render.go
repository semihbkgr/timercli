package main

import (
	"github.com/nsf/termbox-go"
	"strings"
)

const textCharHeight = 5

const (
	zero = `
######
#    #
#    #
#    #
######
`
	one = `
     #
     #
     #
     #
     #
`
	two = `
######
     #
######
#
######
`
	three = `
######
     #
   ###
     #
######
`
	four = `
#
#
#   #
######
    #
`
	five = `
######
#
######
     #
######
`
	six = `
######
#
######
#    #
######
`
	seven = `
######
     #
     #
     #
     #
`
	eight = `
######
#    #
######
#    #
######
`
	nine = `
######
#    #
######
     #
######
`
	space = `
..
..
..
..
..
`
	column = `
..
##
..
##
..
`
	cross = `
##   ##
 ## ##
  ###
 ## ##
##   ##
`
)

type text [][]bool

func (t text) width() int {
	return len((t)[0])
}

func (t text) height() int {
	return len(t)
}

func (t text) iterate(f func(x, t int, b bool)) {
	for y, line := range t {
		for x, b := range line {
			f(x, y, b)
		}
	}
}

var textChars = map[rune]text{
	'0': stringToTextChar(zero),
	'1': stringToTextChar(one),
	'2': stringToTextChar(two),
	'3': stringToTextChar(three),
	'4': stringToTextChar(four),
	'5': stringToTextChar(five),
	'6': stringToTextChar(six),
	'7': stringToTextChar(seven),
	'8': stringToTextChar(eight),
	'9': stringToTextChar(nine),
	' ': stringToTextChar(space),
	':': stringToTextChar(column),
	'x': stringToTextChar(cross),
}

func stringToTextChar(s string) text {
	lines := strings.Split(s, "\n")
	matrix := make([][]bool, textCharHeight)
	var width int
	for _, line := range lines {
		if l := len(line); width < l {
			width = l
		}
	}
	for y := 0; y < textCharHeight; y++ {
		line := lines[y+1]
		array := make([]bool, width)
		for x := 0; x < width; x++ {
			array[x] = x < len(line) && line[x] == '#'
		}
		matrix[y] = array
	}
	return matrix
}

func textChar(r rune) text {
	t, ok := textChars[r]
	if !ok {
		return textChars['x']
	}
	return t
}

func convertToText(s string) text {
	texts := make([]text, len(s))
	for i, c := range s {
		texts[i] = textChar(c)
	}
	return concatTexts(texts...)
}

func concatTexts(t ...text) text {
	ct := make(text, textCharHeight)
	var width int
	for _, c := range t {
		width += c.width()
	}
	for y := 0; y < len(ct); y++ {
		line := make([]bool, width+len(t)-1)
		x := 0
		for i, c := range t {
			for _, b := range c[y] {
				line[x] = b
				x++
			}
			if i < len(t)-1 {
				line[x] = false
				x++
			}
		}
		ct[y] = line
	}
	return ct
}

func initTermbox() {
	checkErr(termbox.Init())
	termbox.SetInputMode(termbox.InputEsc)
}

func closeTermbox() {
	termbox.Close()
}

type Theme struct {
	bg termbox.Attribute
	fg termbox.Attribute
}

var themes = map[string]Theme{
	"light": {
		bg: termbox.ColorWhite,
		fg: termbox.ColorBlack,
	},
	"dark": {
		bg: termbox.ColorBlack,
		fg: termbox.ColorWhite,
	},
}

type Renderer struct {
	t Theme
}

func NewRenderer(theme string) *Renderer {
	t, ok := themes[theme]
	if !ok {
		t = themes["light"]
	}
	return &Renderer{
		t: t,
	}
}

func (r *Renderer) Render(s string) {
	err := termbox.Clear(r.t.bg, r.t.bg)
	checkErr(err)
	w, h := termbox.Size()
	t := convertToText(s)
	tx := (w - t.width()) / 2
	ty := (h - t.height()) / 2
	t.iterate(func(x, y int, b bool) {
		if b {
			termbox.SetCell(tx+x, ty+y, ' ', r.t.fg, r.t.fg)
		}
	})
	err = termbox.Flush()
	checkErr(err)
}
