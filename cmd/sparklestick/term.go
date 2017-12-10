// +build !pihardware

package main

import (
	"github.com/tomnz/sparklestick/runner"
	"github.com/tomnz/sparklestick/scenes"
	"github.com/tomnz/sparklestick/term"
)

// GetDevices returns input/output devices for a local terminal. This is intended for use locally when
// the actual Pi hardware is not present (see i2c.go for hardware implementation).
// The visual output is rendered into the terminal, and the keyboard is used for button presses.
func GetDevices(config *scenes.Config, exit chan bool) (runner.Display, runner.Buttons, runner.Pixel) {
	trm := term.New(w, h, config, exit)
	return trm, trm, trm
}

const (
	// Typically these are dynamically pulled from the hardware device, so
	// we need to define them ourselves for the terminal
	w = 7
	h = 17
)
