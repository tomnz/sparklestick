package candle

import (
	"math"
	"math/rand"
	"time"
)

// Candle implements the candle effect.
type Candle struct {
	w, h   int
	config *Config

	// State
	values [][]float32
	stepWave,
	stepWind,
	stepTemp float64

	// Output
	outBuf [][]byte
}

// New instantiates a new Candle.
func New(w, h int, config *Config) *Candle {
	values := make([][]float32, h)
	for y := range values {
		values[y] = make([]float32, w)
	}
	outBuf := make([][]byte, h)
	for y := range outBuf {
		outBuf[y] = make([]byte, w)
	}

	return &Candle{
		w:      w,
		h:      h,
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
func (c *Candle) Render(total, elapsed time.Duration) [][]byte {
	stepAmt := elapsed.Seconds()
	c.stepWave += rand.Float64() * stepAmt
	c.stepWind += rand.Float64() * stepAmt
	c.stepTemp += rand.Float64() * stepAmt

	next := make([][]float32, c.h)
	for y := range next {
		row := make([]float32, c.w)
		copy(row, c.values[y])
		next[y] = row
	}

	cfgA := btn1Configs[c.config.Btn1Config]
	cfgB := btn2Configs[c.config.Btn2Config]
	cfgC := btn3Configs[c.config.Btn3Config]

	ignition := cfgA.baseIgnition + float32(50.0*math.Sin(c.stepTemp*50.0))

	// TODO: Generate these dynamically based on w/h?
	c.values[15][2] = ignition
	c.values[15][3] = ignition
	c.values[15][4] = ignition
	c.values[16][3] = ignition

	// Shift sin wave upwards, scale to [0.05, 0.8], with extended low value
	wind := float32((math.Sin(c.stepWind*1.5*cfgC.speed)+1.0)/2.0 - 0.2)
	if wind < 0.05 {
		wind = 0.05
	}

	for x := 0; x < c.w; x++ {
		for y := 0; y < c.h; y++ {
			wave := float32(math.Sin((float64(y)/30.0)+c.stepWave*50.0*cfgC.speed)*(float64((c.h-y))/20.0)) * cfgB.windAmt
			value := float32(0.0)
			for i := -1; i <= 1; i++ {
				for j := 0; j <= 2; j++ {
					value += c.getPixel(
						float32(x)+float32(i)*cfgA.searchWidth+wave*wind,
						float32(y+j),
					)
				}
			}
			value /= 10
			next[y][x] = value
		}
	}

	c.values = next

	for y, row := range next {
		for x, val := range row {
			c.outBuf[y][x] = c.mapValue(val)
		}
	}
	return c.outBuf
}

func (c *Candle) mapValue(value float32) byte {
	v := int(value*2.5) - 50
	if v > 255 {
		return 255
	} else if v < 0 {
		return 0
	}
	return byte(v)
}

func (c *Candle) getPixel(x, y float32) float32 {
	if x < 0 || y < 0 || x >= float32(c.w) || y >= float32(c.h) {
		return 0.0
	}

	xInt := int(x)
	yInt := int(y)
	f := x - float32(xInt)

	next := xInt + 1
	if next >= c.w {
		next = c.w - 1
	}

	return (c.values[yInt][xInt] * (1.0 - f)) + (c.values[yInt][next] * f)
}

// Button1 handles a press of button 1.
func (c *Candle) Button1() {
	c.config.Btn1Config = (c.config.Btn1Config + 1) % len(btn1Configs)
}

// Button2 handles a press of button 2.
func (c *Candle) Button2() {
	c.config.Btn2Config = (c.config.Btn2Config + 1) % len(btn2Configs)
}

// Button3 handles a press of button 3.
func (c *Candle) Button3() {
	c.config.Btn3Config = (c.config.Btn3Config + 1) % len(btn3Configs)
}

type btn1 struct {
	searchWidth,
	baseIgnition float32
}

type btn2 struct {
	windAmt float32
}

type btn3 struct {
	speed float64
}

var (
	btn1Configs = []*btn1{
		{
			searchWidth:  1.0,
			baseIgnition: 500.0,
		},
		{
			searchWidth:  0.8,
			baseIgnition: 450.0,
		},
		{
			searchWidth:  0.5,
			baseIgnition: 400.0,
		},
	}

	btn2Configs = []*btn2{
		{
			windAmt: 1.0,
		},
		{
			windAmt: 4.0,
		},
		{
			windAmt: 0.5,
		},
	}

	btn3Configs = []*btn3{
		{
			speed: 1.0,
		},
		{
			speed: 1.4,
		},
		{
			speed: 0.6,
		},
	}
)
