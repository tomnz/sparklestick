package term

import (
	termbox "github.com/nsf/termbox-go"
)

// SetColor sets the pixel color.
func (t *Term) SetColor(r, g, b byte) {
	termbox.SetCell(t.w+2, 1, '█', termColor(r, g, b), termColor(0, 0, 0))
}
