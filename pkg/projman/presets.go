package projman

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Preset struct {
	Name    string   `yaml:"name"`
	Folders []string `yaml:"folders"`
}

func LoadPreset(dir, name string) (Preset, error) {
	var p Preset
	path := filepath.Join(dir, name+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return p, err
	}
	err = yaml.Unmarshal(data, &p)
	return p, err
}
