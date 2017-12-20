// +build pihardware

package main

import (
	"log"

	"github.com/tomnz/button-shim-go"
	"github.com/tomnz/scroll-phat-hd-go"
	"github.com/tomnz/sparklestick/runner"
	"github.com/tomnz/sparklestick/scenes"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

// GetDevices returns input/output devices associated with actual connected hardware on the I2C bus.
// Scroll pHAT HD provides the display.
// Button SHIM provides the buttons and LED pixel.
func GetDevices(*scenes.Config, chan bool) (runner.Display, runner.Buttons, runner.Pixel) {
	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open("1")
	if err != nil {
		log.Fatal(err)
	}

	phat, err := scrollphathd.NewDriver(bus, scrollphathd.WithRotation(scrollphathd.Rotation270))
	if err != nil {
		log.Fatal(err)
	}

	shim, err := buttonshim.New(bus)
	if err != nil {
		log.Fatal(err)
	}

	runnerShim := &buttonShimWrap{shim}

	return phat, runnerShim, runnerShim
}

// buttonShimWrap provides a simple wrapper function around the Button SHIM to register press handlers.
type buttonShimWrap struct {
	*buttonshim.Driver
}

func (s *buttonShimWrap) RegisterHandlers(a, b, c, d, e func()) {
	aPress := s.ButtonPressChan(buttonshim.ButtonA)
	bPress := s.ButtonPressChan(buttonshim.ButtonB)
	cPress := s.ButtonPressChan(buttonshim.ButtonC)
	dPress := s.ButtonPressChan(buttonshim.ButtonD)
	ePress := s.ButtonPressChan(buttonshim.ButtonE)
	go func() {
		for {
			// These shouldn't be coming in fast enough that we need to process in
			// parallel (i.e. spawn a new goroutine per handler call)
			select {
			case <-aPress:
				a()
			case <-bPress:
				b()
			case <-cPress:
				c()
			case <-dPress:
				d()
			case <-ePress:
				e()
			}
		}
	}()
}
