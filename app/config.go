package app

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	SoundsEnabled bool
	BaseDir       string
	NavUpSound    string
	NavDownSound  string
	SelectSound   string
	ErrorSound    string
	ConfirmSound  string
	TaggingFormat string
	TaggingStart  int
	FolderPresets []string
}

var env = map[string]string{
	"PROJMAN_BASE_DIR":               "",
	"PROJMAN_SOUND_ENABLED":          "",
	"PROJMAN_SOUND_NAV_UP":           "",
	"PROJMAN_SOUND_NAV_DOWN":         "",
	"PROJMAN_SOUND_SELECT":           "",
	"PROJMAN_SOUND_CONFIRM":          "",
	"PROJMAN_SOUND_ERROR":            "",
	"PROJMAN_TAGGING_FORMAT":         "",
	"PROJMAN_TAGGING_START":          "",
	"PROJMAN_PROJECT_FOLDER_PRESETS": "",
}
var config = Config{
	SoundsEnabled: false,
	BaseDir:       "",
	NavUpSound:    "sounds/nav_up.wav",
	NavDownSound:  "sounds/nav_down.wav",
	SelectSound:   "sounds/select.wav",
	ErrorSound:    "sounds/error.wav",
	ConfirmSound:  "sounds/confirm.wav",
	TaggingFormat: "{category}-{subcat}-{id}",
	TaggingStart:  1,
}

func SaveConfig(config Config) {
	env["PROJMAN_BASE_DIR"] = config.BaseDir
	env["PROJMAN_SOUND_ENABLED"] = strconv.FormatBool(config.SoundsEnabled)
	env["PROJMAN_SOUND_NAV_UP"] = config.NavUpSound
	env["PROJMAN_SOUND_NAV_DOWN"] = config.NavDownSound
	env["PROJMAN_SOUND_SELECT"] = config.SelectSound
	env["PROJMAN_SOUND_CONFIRM"] = config.ConfirmSound
	env["PROJMAN_SOUND_ERROR"] = config.ErrorSound
	env["PROJMAN_TAGGING_FORMAT"] = config.TaggingFormat
	env["PROJMAN_TAGGING_START"] = strconv.Itoa(config.TaggingStart)

	err := godotenv.Write(env, config.BaseDir+"Config/projman.conf")
	if err != nil {
		fmt.Printf("Failed to save config: %s", err)
	}
}

func LoadConfig() Config {
	if err := godotenv.Load(config.BaseDir + "Config/projman.conf"); err != nil {
		log.Fatalf("Loading environment variables failed: %s", err)
	}

	config.BaseDir = os.Getenv("PROJMAN_BASE_DIR")
	if config.BaseDir == "" {
		config.BaseDir = GetDefaultBaseDir()
	}

	config.SoundsEnabled = strings.ToLower(os.Getenv("PROJMAN_SOUND_ENABLED")) == "true"
	config.NavUpSound = os.Getenv("PROJMAN_SOUND_NAV_UP")
	config.NavDownSound = os.Getenv("PROJMAN_SOUND_NAV_DOWN")
	config.SelectSound = os.Getenv("PROJMAN_SOUND_SELECT")
	config.ConfirmSound = os.Getenv("PROJMAN_SOUND_CONFIRM")
	config.ErrorSound = os.Getenv("PROJMAN_SOUND_ERROR")

	config.TaggingFormat = os.Getenv("PROJMAN_TAGGING_FORMAT")
	if val := os.Getenv("PROJMAN_TAGGING_START"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			config.TaggingStart = i
		}
	}
	return config
}
