// Package swirl implements a swirl effect.
// Roughly based on the example from the Python Scroll pHAT HD library:
// https://github.com/pimoroni/scroll-phat-hd/blob/master/examples/swirl.py
package swirl

import (
	"math"
	"time"
)

// Swirl implements the swirl effect.
type Swirl struct {
	w, h   float64 // Converting to float64 helps later
	config *Config

	// State
	values [][]float32

	// Output
	outBuf [][]byte
}

// New instantiates a new Swirl.
func New(w, h int, config *Config) *Swirl {
	values := make([][]float32, h)
	for y := range values {
		values[y] = make([]float32, w)
	}
	outBuf := make([][]byte, h)
	for y := range outBuf {
		outBuf[y] = make([]byte, w)
	}

	return &Swirl{
		w:      float64(w),
		h:      float64(h),
		config: config,
		values: values,
		outBuf: outBuf,
	}
}

// Config defines the configuration that the effect accepts.
type Config struct {
	Btn1Config int `json:"btn1Config"`
	Btn2Config int `json:"btn2Config"`
	Btn3Config int `json:"btn3Config"`
}

// DefaultConfig returns the default configuration for the effect.
func DefaultConfig() Config {
	return Config{}
}

// Render returns the result of rendering the effect with the given time constraints.
func (s *Swirl) Render(total, elapsed time.Duration) [][]byte {
	cfgA := btn1Configs[s.config.Btn1Config]
	cfgB := btn2Configs[s.config.Btn2Config]
	cfgC := btn3Configs[s.config.Btn3Config]

	timestep := math.Sin(total.Seconds()*cfgA.rate) * cfgB.amplitude

	for y, row := range s.outBuf {
		for x := range row {
			// Center the x/y vals
			xVal := float64(x) - s.w/2.0
			yVal := float64(y) - s.h/2.0

			// Distance from center
			dist := math.Sqrt(xVal*xVal + yVal*yVal)

			angle := (timestep / 10.0) + dist*cfgC.density

			sin := math.Sin(angle)
			cos := math.Cos(angle)

			xs := xVal*cos - yVal*sin
			ys := xVal*sin + yVal*cos

			r := math.Abs(xs + ys)
			r /= 8.0
			if r < 0.0 {
				r = 0.0
			} else if r > 1.0 {
				r = 1.0
			}
			r = 1.0 - r

			row[x] = byte(r * 255)
		}
	}

	return s.outBuf
}

// Button1 handles a press of button 1.
func (s *Swirl) Button1() {
	s.config.Btn1Config = (s.config.Btn1Config + 1) % len(btn1Configs)
}

// Button2 handles a press of button 2.
func (s *Swirl) Button2() {
	s.config.Btn2Config = (s.config.Btn2Config + 1) % len(btn2Configs)
}

// Button3 handles a press of button 3.
func (s *Swirl) Button3() {
	s.config.Btn3Config = (s.config.Btn3Config + 1) % len(btn3Configs)
}

type btn1 struct {
	rate float64
}

type btn2 struct {
	amplitude float64
}

type btn3 struct {
	density float64
}

var (
	btn1Configs = []*btn1{
		{
			rate: 0.05,
		},
		{
			rate: 0.08,
		},
		{
			rate: 0.02,
		},
	}

	btn2Configs = []*btn2{
		{
			amplitude: 1500.0,
		},
		{
			amplitude: 2500.0,
		},
		{
			amplitude: 800.0,
		},
	}

	btn3Configs = []*btn3{
		{
			density: 0.6,
		},
		{
			density: 0.8,
		},
		{
			density: 0.3,
		},
		{
			density: 0.1,
		},
	}
)
