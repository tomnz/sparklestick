package scenes

import (
	"encoding/json"
	"time"

	"github.com/tomnz/sparklestick/scenes/internal/candle"
	"github.com/tomnz/sparklestick/scenes/internal/random"
	"github.com/tomnz/sparklestick/scenes/internal/spread"
	"github.com/tomnz/sparklestick/scenes/internal/swirl"
)

// ConfigVersion defines the configuration version for this sparklestick binary. This value
// is incremented whenever a breaking change is made to the configuration.
const ConfigVersion = 2

// Config defines the available configuration for sparklestick.
type Config struct {
	Version int `json:"version"`

	Brightness byte `json:"brightness"`
	Scene      Type `json:"scene"`

	Candle candle.Config `json:"candle"`
	Random random.Config `json:"random"`
	Spread spread.Config `json:"spread"`
	Swirl  swirl.Config  `json:"swirl"`

	MinFrameInterval time.Duration `json:"minFrameInterval"`
	ShowFPS          bool          `json:"showFps"`
}

// DefaultConfig returns the default sparklestick configuration. This is used when a config
// file is not being used, or does not exist.
func DefaultConfig() *Config {
	return &Config{
		Version:          ConfigVersion,
		Brightness:       byte(127),
		Scene:            "candle",
		Candle:           candle.DefaultConfig(),
		Random:           random.DefaultConfig(),
		Spread:           spread.DefaultConfig(),
		MinFrameInterval: time.Millisecond * 10, // Max out at 100FPS
		ShowFPS:          true,
	}
}

// ConfigFromJSON deserializes configuration from the given JSON data.
func ConfigFromJSON(jsonBytes []byte) (*Config, error) {
	c := &Config{}
	if err := json.Unmarshal(jsonBytes, c); err != nil {
		return nil, err
	}
	return c, nil
}

// ToJSON serializes the configuration as JSON data.
func (c *Config) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}
