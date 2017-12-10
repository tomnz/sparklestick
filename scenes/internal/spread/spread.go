package spread

// TODO: Make this actually cool. Sucks right now...

import (
	"math/rand"
	"time"
)

type Spread struct {
	w, h     int
	config   *Config
	pixels   [][]float32
	fronts   []*front
	waves    []*wave
	lastStep time.Time
	outBuf   [][]byte
}

type front struct {
	x, y     int
	hor, ver direction
}

type wave struct {
	x, y int
	dir  direction
}

type direction byte

const (
	dirUp direction = iota
	dirDown
	dirLeft
	dirRight
)

func New(w, h int, config *Config) *Spread {
	pixels := make([][]float32, h)
	for y := range pixels {
		pixels[y] = make([]float32, w)
	}
	outBuf := make([][]byte, h)
	for y := range outBuf {
		outBuf[y] = make([]byte, w)
	}

	return &Spread{
		w:        w,
		h:        h,
		config:   config,
		pixels:   pixels,
		outBuf:   outBuf,
		lastStep: time.Now(),
	}
}

type Config struct {
	Btn1Config int `json:"btn1Config"`
}

func DefaultConfig() Config {
	return Config{}
}

const (
	stepDuration = 0.1
	waveStrength = float32(0.5)
	dimAmt       = float32(0.5)
)

func (s *Spread) Render(total, elapsed time.Duration) [][]byte {
	cfgA := btn1Configs[s.config.Btn1Config]

	// Get the right number of fronts
	if nFronts := len(s.fronts); nFronts > cfgA.fronts {
		s.fronts = s.fronts[:cfgA.fronts]
	} else if nFronts < cfgA.fronts {
		for i := nFronts; i < cfgA.fronts; i++ {
			s.fronts = append(s.fronts, &front{
				x:   rand.Intn(s.w-1) + 1,
				y:   rand.Intn(s.h-1) + 1,
				ver: dirUp,
				hor: dirLeft,
			})
		}
	}

	curr := time.Now()
	steps := int(curr.Sub(s.lastStep).Seconds() / stepDuration)
	s.lastStep = s.lastStep.Add(time.Millisecond * 100 * time.Duration(steps))

	for i := 0; i < steps; i++ {
		newWaves := make([]*wave, 0, len(s.waves)+len(s.fronts)*2)
		// Advance all of the waves
		for _, w := range s.waves {
			switch w.dir {
			case dirUp:
				w.y--
				if w.y >= 0 {
					newWaves = append(newWaves, w)
				}
			case dirDown:
				w.y++
				if w.y < s.h {
					newWaves = append(newWaves, w)
				}
			case dirLeft:
				w.x--
				if w.x >= 0 {
					newWaves = append(newWaves, w)
				}
			case dirRight:
				w.x++
				if w.x < s.w {
					newWaves = append(newWaves, w)
				}
			}
		}

		// Advance all of the fronts
		for _, f := range s.fronts {
			newX, newY := f.x, f.y
			newVer, newHor := f.ver, f.hor
			// Vertical
			if f.ver == dirUp {
				newY--
				if newY == 0 {
					newVer = dirDown
				}
			} else {
				newY++
				if newY == s.h-1 {
					newVer = dirUp
				}
			}
			newWaves = append(newWaves, &wave{
				x:   f.x,
				y:   newY,
				dir: f.ver,
			})
			// Horizontal
			if f.hor == dirLeft {
				newX--
				if newX == 0 {
					newHor = dirRight
				}
			} else {
				newX++
				if newX == s.w-1 {
					newHor = dirLeft
				}
			}
			newWaves = append(newWaves, &wave{
				x:   newX,
				y:   f.y,
				dir: f.hor,
			})
			f.x = newX
			f.y = newY
			f.ver = newVer
			f.hor = newHor
		}
		s.waves = newWaves

		// Update the field
		for _, w := range s.waves {
			s.pixels[w.y][w.x] += waveStrength
			if s.pixels[w.y][w.x] > 1.0 {
				s.pixels[w.y][w.x] = 1.0
			}
		}
		for _, f := range s.fronts {
			s.pixels[f.y][f.x] += waveStrength
			if s.pixels[f.y][f.x] > 1.0 {
				s.pixels[f.y][f.x] = 1.0
			}
		}

		for y, row := range s.pixels {
			for x := range row {
				s.pixels[y][x] *= dimAmt
			}
		}
	}

	for y, row := range s.pixels {
		for x, val := range row {
			s.outBuf[y][x] = byte(int(val * 255.0))
		}
	}

	return s.outBuf
}

func (s *Spread) Button1() {
	s.config.Btn1Config = (s.config.Btn1Config + 1) % len(btn1Configs)
}

func (s *Spread) Button2() {}

func (s *Spread) Button3() {}

type btn1 struct {
	fronts int
}

type btn2 struct {
}

type btn3 struct {
}

var (
	btn1Configs = []*btn1{
		{
			fronts: 3,
		},
		{
			fronts: 6,
		},
		{
			fronts: 1,
		},
	}

	btn2Configs = []*btn2{
		{},
	}

	btn3Configs = []*btn3{
		{},
	}
)
