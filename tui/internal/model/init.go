package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.Status.Init(),
		m.Accounts.Init(),
		m.BlockExplorer.Init(),
		m.Configs.Init(),
		m.Tabs.Init(),
		m.About.Init(),
		m.Utilities.Init(),
	)
}
