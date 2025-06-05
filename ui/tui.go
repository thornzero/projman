package tui

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type menuOption int

const (
	optionListProjects menuOption = iota
	optionCreateProject
	optionViewProject
	optionArchiveProject
	optionTools
	optionSettings
	optionQuit
)

var menuItems = []string{
	"📋 Projects",
	"🆕 Create New Project",
	"🔍 View Project Status",
	"📦 Archive Project",
	"🧰 Tools",
	"⚙️ Settings",
	"❌ Quit",
}

type mainMenuModel struct {
	cursor int
}

func (m mainMenuModel) Init() tea.Cmd {
	return nil
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(menuItems)-1 {
				m.cursor++
			}
		case "enter", " ":
			switch menuOption(m.cursor) {
			case optionQuit:
				return m, tea.Quit

			case optionListProjects:
				return newProjectListModel(), nil

			case optionCreateProject:
				return newCreateProjectModel(), nil

			case optionViewProject:
				return newViewProjectModel(), nil

			case optionTools:
				return newToolsModel(), nil

			case optionSettings:
				return newSettingsModel(), nil

			default:
				fmt.Printf("⏳ Selected: %s (functionality not implemented yet)\n", menuItems[m.cursor])
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m mainMenuModel) View() string {
	var b strings.Builder
	b.WriteString("🛠 Projman\n\nUse ↑/↓ to navigate and [Enter] to select\n\n")
	for i, item := range menuItems {
		cursor := " "
		if m.cursor == i {
			cursor = "👉"
		}
		fmt.Fprintf(&b, "%s %s\n", cursor, item)
	}
	return b.String()
}

func Tui() {
	p := tea.NewProgram(mainMenuModel{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
