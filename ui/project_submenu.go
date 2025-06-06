package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/thornzero/projman/app"
)

type projectSubmenuModel struct {
	project app.Project
	choice  int
}

var submenuItems = []string{
	"View Status",
	"Archive Project",
	"Open Folder",
	"Back",
}

func newProjectSubmenuModel(p app.Project) projectSubmenuModel {
	return projectSubmenuModel{project: p}
}

func (m projectSubmenuModel) Init() tea.Cmd {
	return nil
}

func (m projectSubmenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			PlaySound(config.NavUpSound)
			if m.choice > 0 {
				m.choice--
			}
		case "down", "j":
			PlaySound(config.NavDownSound)
			if m.choice < len(submenuItems)-1 {
				m.choice++
			}
		case "enter":
			PlaySound(config.SelectSound)
			switch m.choice {
			case 0: // View Status
				return viewProjectModel{project: &m.project, done: true}, nil
			case 1: // Archive
				_ = app.ZipProjectFolder(m.project.Path, m.project.Path+".zip")
				return mainMenuModel{}, nil
			case 2: // Open Folder
				PlaySound(config.ConfirmSound)
				openCmd := "xdg-open"
				if runtime.GOOS == "darwin" {
					openCmd = "open"
				} else if runtime.GOOS == "windows" {
					openCmd = "explorer"
				}
				_ = exec.Command(openCmd, m.project.Path).Start()
				return mainMenuModel{}, nil
			case 3: // Back
				PlaySound(config.ErrorSound)
				return mainMenuModel{}, nil
			}
		case "esc", "q":
			PlaySound(config.ErrorSound)
			return mainMenuModel{}, nil
		}
	}
	return m, nil
}

func (m projectSubmenuModel) View() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Project: %s - %s\n\n", m.project.ID, m.project.Name))
	for i, item := range submenuItems {
		prefix := "  "
		if i == m.choice {
			prefix = "👉"
		}
		fmt.Fprintf(&b, "%s %s\n", prefix, item)
	}
	b.WriteString("\n[↑/↓] Navigate • [Enter] Select • [Esc] Cancel\n")
	return b.String()
}
