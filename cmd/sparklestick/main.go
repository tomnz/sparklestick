package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/tomnz/sparklestick/runner"
	"github.com/tomnz/sparklestick/scenes"
)

var (
	configFilename = flag.String("config", "", "JSON config file (blank for in-memory using defaults). Will be created if it doesn't exist.")
)

func main() {
	flag.Parse()

	cfg, saveConfig := getConfig()
	exit := make(chan bool, 1)
	// The actual GetDevices function depends on build flags (i2c.go or term.go)
	display, buttons, pixel := GetDevices(cfg, exit)
	allScenes := scenes.GetScenes(display.Width(), display.Height(), cfg)

	rnr := runner.New(
		cfg,
		allScenes,
		scenes.Types,
		display,
		buttons,
		pixel,
		saveConfig,
		exit,
	)

	log.Println("sparklestick is running!")

	if err := rnr.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("sparklestick exited")
}

func getConfig() (*scenes.Config, func()) {
	var (
		cfg        *scenes.Config
		saveConfig func()
	)
	if *configFilename != "" {
		jsonBytes, err := ioutil.ReadFile(*configFilename)
		if os.IsNotExist(err) || len(jsonBytes) == 0 {
			cfg = scenes.DefaultConfig()
		} else if err != nil {
			log.Fatalf("Couldn't read config file: %s", err)
		} else {
			if cfg, err = scenes.ConfigFromJSON(jsonBytes); err != nil {
				log.Fatalf("Invalid config file: %s", err)
			}
			if cfg.Version != scenes.ConfigVersion {
				// TODO: Implement version upgrades
				cfg = scenes.DefaultConfig()
			}
		}

		saveConfig = func() {
			jsonBytes, err := cfg.ToJSON()
			if err != nil {
				panic(err)
			}
			if err := ioutil.WriteFile(*configFilename, jsonBytes, os.ModePerm); err != nil {
				panic(err)
			}
		}
		// Always re-save
		saveConfig()
	} else {
		cfg = scenes.DefaultConfig()
		saveConfig = func() {}
	}
	return cfg, saveConfig
}
