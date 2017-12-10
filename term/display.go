package term

import termbox "github.com/nsf/termbox-go"

// SetBuffer overwrites the current display buffer.
func (t *Term) SetBuffer(buffer [][]byte) {
	t.buffer = buffer
}

// SetBrightness sets the current brightness.
func (t *Term) SetBrightness(brightness byte) {
	t.brightness = brightness
}

// termColor converts the RGB value to a term256 color.
func termColor(r, g, b byte) termbox.Attribute {
	rterm := (((uint16(r) * 5) + 127) / 255) * 36
	gterm := (((uint16(g) * 5) + 127) / 255) * 6
	bterm := (((uint16(b) * 5) + 127) / 255)

	return termbox.Attribute(rterm + gterm + bterm + 16 + 1)
}

// Show outputs the buffer contents to the terminal.
func (t *Term) Show() error {
	for y, row := range t.buffer {
		for x, color := range row {
			color = t.scaleVal(color)
			termbox.SetCell(x, y, '█', termColor(color, color, color), termColor(0, 0, 0))
		}
	}
	termbox.Flush()
	return nil
}

func (t *Term) scaleVal(val byte) byte {
	return byte(uint16(val) * uint16(t.brightness) / 255)
}
