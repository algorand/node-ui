// Package template can be used as a starting point when creating a new bubble.
package template

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents a generic bubble
type Model struct {
	width  int
	height int
}

// New constructs the Model.
func New() Model {
	return Model{}
}

// Init is part of the tea.Model interface.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is part of the tea.Model interface.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View is part of the tea.Model interface.
func (m Model) View() string {
	return "template"
}
