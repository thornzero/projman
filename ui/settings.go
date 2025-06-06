package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type settingsModel struct {
	cursor  int
	toggles []bool
}

var settingsItems = []string{
	"Enable Sound Effects",
	"Back to Menu",
}

func newSettingsModel() settingsModel {
	return settingsModel{
		toggles: []bool{config.SoundsEnabled},
	}
}

func (m settingsModel) Init() tea.Cmd {
	return nil
}

func (m settingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return mainMenuModel{}, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(settingsItems)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == 0 {
				m.toggles[0] = !m.toggles[0]
				config.SoundsEnabled = m.toggles[0]
			} else if m.cursor == 1 {
				return mainMenuModel{}, nil
			}
		}
	}
	return m, nil
}

func (m settingsModel) View() string {
	var b strings.Builder
	b.WriteString("⚙️ Settings\n\n")

	for i, item := range settingsItems {
		prefix := "  "
		if m.cursor == i {
			prefix = "👉"
		}

		value := ""
		if i == 0 {
			if m.toggles[0] {
				value = "✅ On"
			} else {
				value = "❌ Off"
			}
		}

		b.WriteString(fmt.Sprintf("%s %s %s\n", prefix, item, value))
	}

	b.WriteString("\n[↑/↓] Navigate • [Enter] Toggle/Select • [Esc] Back\n")
	return b.String()
}
