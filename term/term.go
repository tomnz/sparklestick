package term

import (
	"log"

	termbox "github.com/nsf/termbox-go"
	"github.com/tomnz/sparklestick/scenes"
)

// Term implements the input/output device interfaces for local execution in the terminal.
// Visual elements are displayed in the terminal, and key presses are handled via the keyboard.
type Term struct {
	w, h       int
	config     *scenes.Config
	buffer     [][]byte
	brightness byte
	exit       chan<- bool
	pixel      termbox.Attribute

	handlerA,
	handlerB,
	handlerC,
	handlerD,
	handlerE func()
}

// New returns a new terminal device with the given configuration.
func New(w, h int, config *scenes.Config, exit chan<- bool) *Term {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.SetOutputMode(termbox.Output256)

	log.Println("Press Q to exit!")

	t := &Term{
		w:      w,
		h:      h,
		config: config,
		exit:   exit,
	}
	go t.pollButtons()
	return t
}

// Width returns the display width.
func (t *Term) Width() int {
	return t.w
}

// Height returns the display height.
func (t *Term) Height() int {
	return t.h
}
