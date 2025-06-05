package tui

import (
	"fmt"
	"os/exec"
	"runtime"
)

type SoundConfig struct {
	Enabled      bool
	NavUpSound   string
	NavDownSound string
	SelectSound  string
	ErrorSound   string
	ConfirmSound string
}

var Sounds = SoundConfig{
	Enabled:      false,
	NavUpSound:   "sounds/nav_up.wav",
	NavDownSound: "sounds/nav_down.wav",
	SelectSound:  "sounds/select.wav",
	ErrorSound:   "sounds/error.wav",
	ConfirmSound: "sounds/confirm.wav",
}

func PlaySound(path string) {
	if !Sounds.Enabled {
		return
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("aplay", path)
	case "darwin":
		cmd = exec.Command("afplay", path)
	case "windows":
		cmd = exec.Command("powershell", "-c", fmt.Sprintf(`(New-Object Media.SoundPlayer '%s').PlaySync();`, path))
	default:
		return
	}

	_ = cmd.Start()
}
