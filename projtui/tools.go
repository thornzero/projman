package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/thornzero/projman/core"
)

type toolsModel struct {
	cursor     int
	mode       string
	inputCSV   textinput.Model
	outputYAML textinput.Model
	message    string
}

var toolItems = []string{
	"🏷️ Generate Tags",
	"⬅️ Back",
}

func newToolsModel() toolsModel {
	input := textinput.New()
	input.Placeholder = "Path to input CSV"
	input.CharLimit = 128
	input.Width = 40

	out := textinput.New()
	out.Placeholder = "Output YAML file path"
	out.CharLimit = 128
	out.Width = 40

	return toolsModel{
		inputCSV:   input,
		outputYAML: out,
	}
}

func (m toolsModel) Init() tea.Cmd {
	return nil
}

func (m toolsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case "generate":
			switch msg.String() {
			case "enter":
				csv := m.inputCSV.Value()
				out := m.outputYAML.Value()
				if csv != "" && out != "" {
					err := core.GenerateTags(csv, out)
					if err != nil {
						m.message = "❌ Failed: " + err.Error()
					} else {
						m.message = "✅ Tags generated successfully."
					}
					m.mode = ""
				}
			case "esc":
				m.mode = ""
				m.message = ""
			default:
				m.inputCSV, _ = m.inputCSV.Update(msg)
				m.outputYAML, _ = m.outputYAML.Update(msg)
			}
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(toolItems)-1 {
				m.cursor++
			}
		case "enter":
			switch m.cursor {
			case 0:
				m.mode = "generate"
				m.inputCSV.Focus()
			case 1:
				return mainMenuModel{}, nil
			}
		case "esc", "q":
			return mainMenuModel{}, nil
		}
	}
	return m, nil
}

func (m toolsModel) View() string {
	if m.mode == "generate" {
		return fmt.Sprintf(
			"🛠 Generate Tags\n\nInput CSV:\n%s\n\nOutput YAML:\n%s\n\n[enter] Generate • [esc] Cancel\n\n%s",
			m.inputCSV.View(),
			m.outputYAML.View(),
			m.message,
		)
	}

	s := "🧰 Tools\n\n"
	for i, item := range toolItems {
		prefix := "  "
		if m.cursor == i {
			prefix = "👉"
		}
		s += fmt.Sprintf("%s %s\n", prefix, item)
	}
	s += "\n[↑/↓] Navigate • [Enter] Select • [Esc] Back\n"
	return s
}
