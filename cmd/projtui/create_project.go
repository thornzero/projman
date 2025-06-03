package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/thornzero/projman/pkg/projman"
)

type createProjectModel struct {
	inputs  []textinput.Model
	focus   int
	baseDir string
	done    bool
	message string
}

func newCreateProjectModel() createProjectModel {
	base := projman.GetDefaultBaseDir()
	fields := []string{"Project ID", "Project Name", "Description", "Tags (comma-separated)"}
	inputs := make([]textinput.Model, len(fields))
	for i := range inputs {
		ti := textinput.New()
		ti.Placeholder = fields[i]
		ti.Focus()
		ti.CharLimit = 100
		ti.Width = 40
		if i != 0 {
			ti.Blur()
		}
		inputs[i] = ti
	}

	return createProjectModel{
		inputs:  inputs,
		focus:   0,
		baseDir: base,
	}
}

func (m createProjectModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m createProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return mainMenuModel{}, nil
		case "enter":
			if m.focus == len(m.inputs)-1 {
				id := projman.ValidateID(m.inputs[0].Value())
				name := m.inputs[1].Value()
				desc := m.inputs[2].Value()
				tags := m.inputs[3].Value()

				if id == "" || name == "" {
					m.message = "❌ ID and Name are required"
					return m, nil
				}

				projman.CreateProject(m.baseDir, projman.Params{
					ID:          id,
					Name:        name,
					Description: desc,
					Tags:        tags,
					Status:      "active",
				})

				m.done = true
				m.message = fmt.Sprintf("✅ Project %s created!", id)
				return m, nil
			}

			m.inputs[m.focus].Blur()
			m.focus = (m.focus + 1) % len(m.inputs)
			m.inputs[m.focus].Focus()
			return m, nil

		case "tab", "down":
			m.inputs[m.focus].Blur()
			m.focus = (m.focus + 1) % len(m.inputs)
			m.inputs[m.focus].Focus()
		case "shift+tab", "up":
			m.inputs[m.focus].Blur()
			m.focus = (m.focus - 1 + len(m.inputs)) % len(m.inputs)
			m.inputs[m.focus].Focus()
		}
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m createProjectModel) View() string {
	if m.done {
		return fmt.Sprintf("%s\n\n[esc] Back to menu", m.message)
	}

	var b strings.Builder
	b.WriteString("🆕 Create New Project\n\n")
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View() + "\n")
	}
	b.WriteString("\n[tab] to switch • [enter] to submit • [esc] cancel\n")
	if m.message != "" {
		b.WriteString("\n" + m.message + "\n")
	}
	return b.String()
}
