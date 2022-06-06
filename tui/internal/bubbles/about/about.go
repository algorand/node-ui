package about

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/muesli/reflow/indent"
)

// Model represents the about bubble.
type Model struct {
	heightMargin int
	viewport     viewport.Model
}

// New creates the about Model.
func New(heightMargin int, content string) Model {
	m := Model{
		heightMargin: heightMargin,
		viewport:     viewport.New(0, 0),
	}
	m.setSize(80, 20)

	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(80),
		glamour.WithEmoji(),
	)
	c, _ := r.Render(content)
	m.viewport.SetContent(indent.String(c, 7))
	return m
}

// Init is part of the tea.Model interface.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is part of the tea.Model interface.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View is part of the tea.Model interface.
func (m Model) View() string {
	return m.viewport.View()
}

func (m *Model) setSize(width, height int) {
	m.viewport.Width = width
	m.viewport.Height = height - m.heightMargin
}
