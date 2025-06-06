package ui

import (
	"fmt"
	"os/exec"
	"runtime"
)

func PlaySound(path string) {
	if !config.SoundsEnabled {
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
