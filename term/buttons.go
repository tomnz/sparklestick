package term

import (
	"time"

	"github.com/nsf/termbox-go"
)

// RegisterHandlers registers the given handler functions for each button.
func (t *Term) RegisterHandlers(a, b, c, d, e func()) {
	t.handlerA = a
	t.handlerB = b
	t.handlerC = c
	t.handlerD = d
	t.handlerE = e
}

// pollButtons pools for keyboard events, and either calls the corresponding handler
// function, or terminates execution.
// Designed to be called as a long-running goroutine.
func (t *Term) pollButtons() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc || ev.Ch == 'q' {
				t.exit <- true
			} else if ev.Ch == '1' && t.handlerE != nil {
				t.handlerE()
			} else if ev.Ch == '2' && t.handlerD != nil {
				t.handlerD()
			} else if ev.Ch == '3' && t.handlerC != nil {
				t.handlerC()
			} else if ev.Ch == '4' && t.handlerB != nil {
				t.handlerB()
			} else if ev.Ch == '5' && t.handlerA != nil {
				t.handlerA()
			}
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}
