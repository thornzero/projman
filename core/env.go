package core

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	EnvLoaded    = false
	EnvBaseDir   = ""
	EnvTagFormat = ""
	EnvTagStart  = 1
)

func LoadEnvConfig() {
	_ = godotenv.Load()

	EnvLoaded = true

	EnvBaseDir = os.Getenv("PROJMAN_BASE_DIR")
	if EnvBaseDir == "" {
		EnvBaseDir = GetDefaultBaseDir()
	}

	Sounds.Enabled = strings.ToLower(os.Getenv("PROJMAN_SOUND_ENABLED")) == "true"
	Sounds.NavSound = os.Getenv("PROJMAN_SOUND_NAV")
	Sounds.SelectSound = os.Getenv("PROJMAN_SOUND_SELECT")
	Sounds.ConfirmSound = os.Getenv("PROJMAN_SOUND_CONFIRM")
	Sounds.ErrorSound = os.Getenv("PROJMAN_SOUND_ERROR")

	EnvTagFormat = os.Getenv("PROJMAN_TAG_FORMAT")
	if val := os.Getenv("PROJMAN_TAG_START"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			EnvTagStart = i
		}
	}
}
