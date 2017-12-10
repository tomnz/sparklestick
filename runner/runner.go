// Package runner orchestrates the overall system behavior.
// This includes rendering scenes, handling button presses, and writing to the display.
package runner

import (
	"log"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/tomnz/sparklestick/scenes"
)

// New returns a new Runner.
func New(
	config *scenes.Config,
	scenes map[scenes.Type]scenes.Scene,
	sceneTypes []scenes.Type,
	display Display, buttons Buttons,
	pixel Pixel,
	saveConfig func(),
	exit <-chan bool,
) *Runner {
	rnr := &Runner{
		config:     config,
		scenes:     scenes,
		sceneTypes: sceneTypes,
		display:    display,
		buttons:    buttons,
		pixel:      pixel,
		saveConfig: saveConfig,
		exit:       exit,
	}
	rnr.setup()
	return rnr
}

// Runner orchestrates the overall system behavior.
type Runner struct {
	config     *scenes.Config
	scenes     map[scenes.Type]scenes.Scene
	sceneTypes []scenes.Type
	display    Display
	buttons    Buttons
	pixel      Pixel
	saveConfig func()
	exit       <-chan bool
}

// Display defines the behavior that the runner expects for the display object.
type Display interface {
	SetBuffer([][]byte)
	SetBrightness(brightness byte)
	Show() error
	Width() int
	Height() int
}

// Buttons defines the behavior that the runner expects for the buttons.
type Buttons interface {
	RegisterHandlers(a, b, c, d, e func())
}

// Pixel defines the behavior that the runner expects for the LED pixel.
type Pixel interface {
	SetColor(r, g, b byte)
	SetBrightness(brightness byte)
}

// Run is the main orchestration loop. It has several responsibilities.
//
// First, it kicks off a goroutine to send color data to the LED pixel in a rainbow pattern.
// It keeps track of total and per-frame timings. It uses these to calculate FPS, as well as trigger
// the current scene to render.
// Finally, tt sends the resulting rendered buffer to the display, and pauses if needed to limit the
// maximum FPS.
func (r *Runner) Run() error {
	go func() {
		h := 0
		for {
			color := colorful.Hsv(float64(h), 1, 0.5)
			r.pixel.SetColor(
				byte(color.R*255),
				byte(color.G*255),
				byte(color.B*255),
			)
			h = (h + 2) % 360
			time.Sleep(time.Millisecond * 100)
		}
	}()

	buffer := make([][]byte, r.display.Height())
	for y := range buffer {
		buffer[y] = make([]byte, r.display.Width())
	}

	start := time.Now()
	last := time.Now().Add(time.Millisecond * -10)

	// Only used if displaying FPS
	lastFPS := start
	nextFPS := start.Add(time.Second * 10)
	fpsFrames := uint64(0)

	for {
		select {
		case <-r.exit:
			return nil
		default:
		}

		curr := time.Now()
		total := curr.Sub(start)
		elapsed := curr.Sub(last)
		last = curr

		if r.config.ShowFPS {
			fpsFrames++
			if curr.After(nextFPS) {
				log.Printf("FPS: %2f", float64(fpsFrames)/curr.Sub(lastFPS).Seconds())
				nextFPS = curr.Add(time.Minute * 1)
				lastFPS = curr
				fpsFrames = 0
			}
		}

		output := r.scenes[r.config.Scene].Render(total, elapsed)

		for y, row := range buffer {
			for x := range row {
				row[x] = output[y][x]
			}
		}

		r.display.SetBuffer(buffer)
		if err := r.display.Show(); err != nil {
			return err
		}

		// Sleep if we want to limit framerate
		frameTime := time.Since(curr)
		if frameTime < r.config.MinFrameInterval {
			time.Sleep(r.config.MinFrameInterval - frameTime)
		} else {
			// Otherwise sleep for a moment anyway - increases the chance that
			// button presses are detected
			time.Sleep(500 * time.Nanosecond)
		}
	}
}

func (r *Runner) setup() {
	r.display.SetBrightness(r.config.Brightness)
	r.pixel.SetBrightness(r.config.Brightness)
	r.buttons.RegisterHandlers(
		// The orientation of the shim means the letters are "reversed" top-bottom
		// Bottom three buttons are for the scenes
		func() {
			r.scenes[r.config.Scene].Button3()
			r.saveConfig()
		},
		func() {
			r.scenes[r.config.Scene].Button2()
			r.saveConfig()
		},
		func() {
			r.scenes[r.config.Scene].Button1()
			r.saveConfig()
		},
		// Brightness
		func() {
			// Will overflow/reset
			r.config.Brightness += 32
			r.display.SetBrightness(r.config.Brightness)
			r.pixel.SetBrightness(r.config.Brightness)
			r.saveConfig()
		},
		// Scene selection
		func() {
			for idx, scene := range r.sceneTypes {
				if r.config.Scene == scene {
					nextIdx := (idx + 1) % len(r.sceneTypes)
					r.config.Scene = r.sceneTypes[nextIdx]
					r.saveConfig()
					return
				}
			}
			panic("current scene has no ordering")
		},
	)
}
