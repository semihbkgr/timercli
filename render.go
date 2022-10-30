package main

import (
	"github.com/nsf/termbox-go"
	"strings"
	"sync"
)

const textCharHeight = 5
const textPixelRune = '#'

const (
	zero = `
.######.
.#    #.
.#    #.
.#    #.
.######.
`
	one = `
.   ## .
.   ## .
.   ## .
.   ## .
.   ## .
`
	two = `
.######.
.     #.
.######.
.#     .
.######.
`
	three = `
.######.
.     #.
.   ###.
.     #.
.######.
`
	four = `
.#     .
.#     .
.#   # .
.######.
.    # .
`
	five = `
.######.
.#     .
.######.
.     #.
.######.
`
	six = `
.######.
.#     .
.######.
.#    #.
.######.
`
	seven = `
.######.
.    ##.
.    ##.
.    ##.
.    ##.
`
	eight = `
.######.
.#    #.
.######.
.#    #.
.######.
`
	nine = `
.######.
.#    #.
.######.
.     #.
.######.
`
	space = `
....
....
....
....
....
`
	column = `

.##.

.##.

`
	cross = `
.##   ##.
. ## ## .
.  ###  .
. ## ## .
.##   ##.
`
)

type text [][]bool

func (t text) width() int {
	return len((t)[0])
}

func (t text) height() int {
	return len(t)
}

func (t text) iterate(f func(x, y int, b bool)) {
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
	for y := 0; y < len(matrix); y++ {
		line := lines[y+1]
		array := make([]bool, width)
		for x := 0; x < width; x++ {
			array[x] = x < len(line) && line[x] == textPixelRune
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
		line := make([]bool, width)
		x := 0
		for _, c := range t {
			for _, b := range c[y] {
				line[x] = b
				x++
			}
		}
		ct[y] = line
	}
	return ct
}

func initTermbox() {
	checkErr(termbox.Init())
}

func closeTermbox() {
	termbox.Close()
}

type Theme struct {
	renderBg  termbox.Attribute
	renderFg  termbox.Attribute
	renderFgS termbox.Attribute
	titleFg   termbox.Attribute
	infoFg    termbox.Attribute
}

var themes = map[string]Theme{
	"dark": {
		renderBg:  termbox.ColorBlack,
		renderFg:  termbox.ColorLightGray,
		renderFgS: termbox.ColorLightRed,
		titleFg:   termbox.ColorLightBlue,
		infoFg:    termbox.ColorLightGreen,
	},
	"light": {
		renderBg:  termbox.ColorLightGray,
		renderFg:  termbox.ColorBlack,
		renderFgS: termbox.ColorRed,
		titleFg:   termbox.ColorBlue,
		infoFg:    termbox.ColorGreen,
	},
}

type Renderer struct {
	*sync.Mutex
	theme Theme
	title string
	info  []string
}

func NewRenderer(t string, title string, info ...string) *Renderer {
	theme, ok := themes[t]
	if !ok {
		theme = themes["dark"]
	}
	return &Renderer{
		Mutex: &sync.Mutex{},
		theme: theme,
		title: title,
		info:  info,
	}
}

func (r *Renderer) Render(s string, f bool) error {
	r.Lock()
	defer r.Unlock()
	err := termbox.Clear(r.theme.renderBg, r.theme.renderBg)
	if err != nil {
		return err
	}
	w, h := termbox.Size()
	for i := 0; i < len(r.info); i++ {
		termboxWriteString(r.info[i], 0, h-len(r.info)+i, r.theme.infoFg, r.theme.renderBg)
	}
	termboxWriteString(r.title, 0, 0, r.theme.titleFg, r.theme.renderBg)
	t := convertToText(s)
	tx := (w - t.width()) / 2
	ty := (h - t.height()) / 2
	t.iterate(func(x, y int, b bool) {
		if b {
			if f {
				termbox.SetCell(tx+x, ty+y, ' ', r.theme.renderFg, r.theme.renderFg)
			} else {
				termbox.SetCell(tx+x, ty+y, ' ', r.theme.renderFgS, r.theme.renderFgS)
			}
		} else {
			termbox.SetCell(tx+x, ty+y, ' ', r.theme.renderBg, r.theme.renderBg)
		}
	})
	termbox.SetCursor(w, h)
	return termbox.Flush()
}

func termboxWriteString(s string, x, y int, fg, bg termbox.Attribute) {
	for i := 0; i < len(s); i++ {
		termbox.SetCell(x+i, y, rune(s[i]), fg, bg)
	}
}
