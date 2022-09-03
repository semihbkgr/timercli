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
	column = `

##

##

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
	for y, line := range lines {
		array := make([]bool, width)
		for x := 0; x < width; x++ {
			array[x] = x < len(line) && line[x] == '#'
		}
		matrix[y] = array
	}
	return matrix
}

func concatTexts(t ...text) text {
	ct := make(text, textCharHeight)
	var width int
	for _, c := range t {
		width += c.width()
	}
	for y := 0; y < len(ct); y++ {
		line := make([]bool, width)
		x := 0
		for _, c := range t {
			for _, b := range c[y] {
				line[x] = b
				x++
			}
		}
	}
	return ct
}

func startTermbox() {
	checkErr(termbox.Init())
}

func closeTermbox() {
	termbox.Close()
}

type Theme struct {
	bg termbox.Attribute
	fg termbox.Attribute
}

var themes = map[string]Theme{
	"light": Theme{
		bg: termbox.ColorWhite,
		fg: termbox.ColorBlue,
	},
	"dark": Theme{
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
	termbox.Clear(r.t.bg, r.t.bg)
}
