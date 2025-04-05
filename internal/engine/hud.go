package engine

import (
	"fmt"
	"strings"
	"unicode"
)

const fontStart = '!'
const fontEnd = '_'

var font []*Patch

func InitHud() {
	fontLen := fontEnd - fontStart + 1
	font = make([]*Patch, fontLen)
	for i := range fontLen {
		font[i] = LumpByName[Patch](fmt.Sprintf("STCFN%03d", fontStart+i))
	}
}

// hudStringHeight returns the pixel height of a string with one or more lines.
func hudStringHeight(s string) int {
	height := int(font[0].Height)
	return height * (1 + strings.Count(s, "\n"))
}

// hudStringWidth returns the pixel width of a string. Newlines are not expected and are counted as
// a space.
func hudStringWidth(s string) int {
	w := 0
	for _, ch := range s {
		fc := unicode.ToUpper(rune(ch)) - fontStart
		if fc < fontStart || fc > fontEnd {
			w += 4
		} else {
			w += int(font[fc].Width)
		}
	}
	return w
}

// HudDrawText draws text at the specified x and y coordinates. If the text contains newlines each
// line starts at the x coordinate and line height is 12 pixels.
func HudDrawText(x int, y int, line string) {
	cx := x
	cy := y
	for _, ch := range line {
		if ch == '\n' {
			cx = x
			cy += int(font[0].Height)
			continue
		}

		fc := int(unicode.ToUpper(ch) - fontStart)
		if fc < 0 || fc > len(font) {
			cx += 4
			continue
		}

		w := int(font[fc].Width)
		if cx+w > ScreenWidth {
			continue
		}

		DrawPatch(cx, cy, 0, font[fc])
		cx += int(font[fc].Width)
	}
}

func HudDrawCenteredText(msg string) {
	y := 100 - hudStringHeight(msg)/2
	for _, line := range strings.Split(msg, "\n") {
		x := 160 - hudStringWidth(line)/2
		HudDrawText(x, y, line)
		y += int(font[0].Height)
	}
}
