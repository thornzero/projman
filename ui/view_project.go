package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/thornzero/projman/app"
)

type viewProjectModel struct {
	input   textinput.Model
	baseDir string
	project *app.Project
	errMsg  string
	done    bool
}

func newViewProjectModel() viewProjectModel {
	ti := textinput.New()
	ti.Placeholder = "Enter Project ID"
	ti.CharLimit = 32
	ti.Width = 30
	ti.Focus()

	return viewProjectModel{
		input:   ti,
		baseDir: app.GetDefaultBaseDir(),
	}
}

func (m viewProjectModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m viewProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return mainMenuModel{}, nil
		case "enter":
			id := app.ValidateID(m.input.Value())
			if id == "" {
				m.errMsg = "❌ Invalid ID"
				return m, nil
			}

			p, err := app.ReadProjectFile(m.baseDir, id)
			if err != nil {
				m.errMsg = fmt.Sprintf("❌ %v", err)
				return m, nil
			}

			m.project = &p
			m.done = true
			return m, nil
		}

	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m viewProjectModel) View() string {
	if m.done && m.project != nil {
		p := m.project
		var b strings.Builder
		b.WriteString("🔍 Project Info\n\n")
		b.WriteString(fmt.Sprintf("ID:          %s\n", p.ID))
		b.WriteString(fmt.Sprintf("Name:        %s\n", p.Name))
		b.WriteString(fmt.Sprintf("Description: %s\n", p.Description))
		b.WriteString(fmt.Sprintf("Status:      %s\n", p.Status))
		b.WriteString(fmt.Sprintf("Tags:        %s\n", strings.Join(p.Tags, ", ")))
		b.WriteString(fmt.Sprintf("Created At:  %s\n", p.CreatedAt))
		b.WriteString(fmt.Sprintf("Path:        %s\n", p.Path))
		b.WriteString("\n[esc] Back to menu")
		return b.String()
	}

	return fmt.Sprintf("🔍 View Project Status\n\n%s\n\n[enter] Lookup • [esc] Cancel\n%s", m.input.View(), m.errMsg)
}
