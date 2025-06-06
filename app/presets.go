package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Preset struct {
	Name    string   `yaml:"name"`
	Folders []string `yaml:"folders"`
}

var (
	presetsDirectory = config.BaseDir + "Config/Presets"
	presets          = []string{}
)

func GetAvailablePresets() []string {
	files, err := os.ReadDir(presetsDirectory)
	if err != nil {
		log.Printf("failed to read preset directory: %s", err)
	}
	presetFiles := []string{}
	for _, file := range files {
		presetFiles = append(presetFiles, strings.Split(file.Name(), ".")[0])
		if strings.Contains(file.Name(), ".yaml") {
			presets = append(presets, file.Name())
		}
	}
}

func LoadPreset(name string) (Preset, error) {
	var p Preset
	path := filepath.Join(presetsDirectory, name+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return p, err
	}
	err = yaml.Unmarshal(data, &p)
	return p, err
}
