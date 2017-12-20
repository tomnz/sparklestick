package scenes

import (
	"time"

	"github.com/tomnz/sparklestick/scenes/internal/candle"
	"github.com/tomnz/sparklestick/scenes/internal/random"
	"github.com/tomnz/sparklestick/scenes/internal/swirl"
)

// Type defines an available scene type.
type Type string

const (
	// SceneCandle is the candle scene.
	SceneCandle Type = "candle"
	// SceneRandom is the random scene.
	SceneRandom Type = "random"
	// SceneSwirl is the swirl scene.
	SceneSwirl Type = "swirl"
)

// Scene defines the expected behavior that a scene should implement.
type Scene interface {
	Render(total, elapsed time.Duration) [][]byte
	Button1()
	Button2()
	Button3()
}

// GetScenes returns a set of initialized scenes, for the given width/height and
// sparklestick configuration.
func GetScenes(w, h int, cfg *Config) map[Type]Scene {
	return map[Type]Scene{
		SceneCandle: candle.New(w, h, &cfg.Candle),
		SceneRandom: random.New(w, h, &cfg.Random),
		SceneSwirl:  swirl.New(w, h, &cfg.Swirl),
	}
}

var (
	// Types lists the available scene types. This is used for ordering the scenes.
	Types = []Type{
		SceneCandle,
		SceneRandom,
		SceneSwirl,
	}
)
