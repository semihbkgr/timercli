package main

import (
	"github.com/nsf/termbox-go"
	"strings"
)

const (
	zero = `
######.
#....#.
#....#.
#....#.
######.
`
	one = `
.....#.
.....#.
.....#.
.....#.
.....#.
`
	two = `
######.
.....#.
######.
#......
######.
`
	three = `
######.
.....#.
...###.
.....#.
######.
`
	four = `
#......
#......
#...#..
######.
....#..
`
	five = `
######.
#......
######.
.....#.
######.
`
	six = `
######.
#......
######.
#....#.
######.
`
	seven = `
######.
.....#.
.....#.
.....#.
.....#.
`
	eight = `
######.
#....#.
######.
#....#.
######.
`
	nine = `
######.
#....#.
######.
.....#.
######.
`
	column = `
..
#.
..
#.
..
`
)

var numbers = [10][][]rune{
	stringToRuneMatrix(one),
	stringToRuneMatrix(two),
	stringToRuneMatrix(three),
	stringToRuneMatrix(four),
	stringToRuneMatrix(five),
	stringToRuneMatrix(six),
	stringToRuneMatrix(seven),
	stringToRuneMatrix(eight),
	stringToRuneMatrix(nine),
}

var separator = stringToRuneMatrix(column)

func stringToRuneMatrix(s string) [][]rune {
	lines := strings.Split(s, "\n")
	matrix := make([][]rune, len(lines))
	for y, line := range lines {
		matrix[y] = []rune(line)[0 : len(line)-1]
	}
	return matrix
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

func (r *Renderer) Render(i ...int) {
	termbox.Clear(r.t.bg, r.t.bg)
}
