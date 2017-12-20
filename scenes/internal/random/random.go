package random

import (
	"math/rand"
	"time"
)

// Random implements the random effect.
type Random struct {
	w, h      int
	config    *Config
	pixels    [][]float32
	lastPixel time.Time
	outBuf    [][]byte
}

// New instantiates a new Random.
func New(w, h int, config *Config) *Random {
	pixels := make([][]float32, h)
	for y := range pixels {
		pixels[y] = make([]float32, w)
	}
	outBuf := make([][]byte, h)
	for y := range outBuf {
		outBuf[y] = make([]byte, w)
	}

	return &Random{
		w:         w,
		h:         h,
		config:    config,
		pixels:    pixels,
		lastPixel: time.Now(),
		outBuf:    outBuf,
	}
}

// Config defines the configuration that the effect accepts.
type Config struct {
	Btn1Config int `json:"btn1Config"`
	Btn2Config int `json:"btn2Config"`
}

// DefaultConfig returns the default configuration for the effect.
func DefaultConfig() Config {
	return Config{}
}

// Render returns the result of rendering the effect with the given time constraints.
func (r *Random) Render(total, elapsed time.Duration) [][]byte {
	cfgA := btn1Configs[r.config.Btn1Config]
	cfgB := btn2Configs[r.config.Btn2Config]

	curr := time.Now()
	if curr.Sub(r.lastPixel).Seconds() > cfgB.newRate {
		x := rand.Intn(r.w)
		y := rand.Intn(r.h)
		r.pixels[y][x] = 1.0
		r.lastPixel = curr
	}

	dimMult := float32(elapsed.Seconds() * cfgA.dimRate)
	for y, row := range r.pixels {
		for x, val := range row {
			r.outBuf[y][x] = byte(int(val * 255.0))
			newVal := val - dimMult
			if newVal < 0 {
				newVal = 0
			}
			row[x] = newVal
		}
	}

	return r.outBuf
}

// Button1 handles a press of button 1.
func (r *Random) Button1() {
	r.config.Btn1Config = (r.config.Btn1Config + 1) % len(btn1Configs)
}

// Button2 handles a press of button 2.
func (r *Random) Button2() {
	r.config.Btn2Config = (r.config.Btn2Config + 1) % len(btn2Configs)
}

// Button3 handles a press of button 3.
func (r *Random) Button3() {}

type btn1 struct {
	dimRate float64
}

type btn2 struct {
	newRate float64
}

var (
	btn1Configs = []*btn1{
		{
			dimRate: 1.0,
		},
		{
			dimRate: 1.5,
		},
		{
			dimRate: 0.5,
		},
	}

	btn2Configs = []*btn2{
		{
			newRate: 0.08,
		},
		{
			newRate: 0.03,
		},
		{
			newRate: 0.12,
		},
	}
)
