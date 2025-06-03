package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/thornzero/projman/pkg/projman"
)

type projectListModel struct {
	all       []projman.Project
	filtered  []projman.Project
	searchBar textinput.Model
	searching bool
	cursor    int
	baseDir   string
}

func newProjectListModel() projectListModel {
	base := projman.GetDefaultBaseDir()
	entries, err := os.ReadDir(base)
	if err != nil {
		return projectListModel{baseDir: base}
	}

	all := []projman.Project{}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		p, err := projman.ReadProjectFile(base, entry.Name())
		if err == nil {
			all = append(all, p)
		}
	}

	input := textinput.New()
	input.Placeholder = "Search Project ID..."
	input.CharLimit = 30
	input.Width = 30

	return projectListModel{
		all:       all,
		filtered:  all,
		searchBar: input,
		baseDir:   base}
}

func (m projectListModel) Init() tea.Cmd {
	return nil
}

func (m projectListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searching {
			switch msg.String() {
			case "esc":
				m.searching = false
				m.searchBar.Blur()
				m.filtered = m.all
				m.cursor = 0
			default:
				m.searchBar, cmd = m.searchBar.Update(msg)
				m.filtered = filterProjects(m.all, m.searchBar.Value())
				if m.cursor >= len(m.filtered) {
					m.cursor = 0
				}
				return m, cmd
			}
		}

		switch msg.String() {
		case "ctrl+f":
			m.searching = true
			m.searchBar.Focus()
			return m, textinput.Blink
		case "ctrl+c", "q", "esc", "b":
			return mainMenuModel{}, nil // back to main menu
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				PlaySound(Sounds.NavUpSound)
			}
		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
				PlaySound(Sounds.NavDownSound)
			}
		case "enter", " ":
			if len(m.filtered) > 0 {
				PlaySound(Sounds.SelectSound)
				return newProjectSubmenuModel(m.filtered[m.cursor]), nil
			}
		}
	}
	return m, nil
}

func (m projectListModel) View() string {
	var b strings.Builder
	b.WriteString("📋 Project List\n\n")

	if m.searching {
		b.WriteString(fmt.Sprintf("🔍 %s\n\n", m.searchBar.View()))
	}

	if len(m.filtered) == 0 {
		b.WriteString("📭 No projects found.\n\n[esc] Back\n")
		return b.String()
	}

	for i, p := range m.filtered {
		prefix := "  "
		if i == m.cursor {
			prefix = "👉"
		}
		line := fmt.Sprintf("%s %-10s  %-25s  %-10s", prefix, p.ID, p.Name, p.Status)
		b.WriteString(line + "\n")
	}

	b.WriteString("\n[↑/↓] Navigate • [ctrl+f] Search • [enter/Spacebar] Select • [esc/b] Back/Cancel\n")
	return b.String()
}

func filterProjects(projects []projman.Project, query string) []projman.Project {
	if query == "" {
		return projects
	}

	var filtered []projman.Project
	q := strings.ToUpper(query)
	for _, p := range projects {
		if strings.Contains(strings.ToUpper(p.ID), q) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
