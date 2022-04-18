package about

import (
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/muesli/reflow/indent"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width        int
	height       int
	heightMargin int
	viewport     viewport.Model
}

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

func (m Model) Init() tea.Cmd {
	return nil
}

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

func (m Model) View() string {

	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("%s", m.viewport.View()))
	return builder.String()
}

func (m *Model) setSize(width, height int) {
	m.viewport.Width = width
	m.viewport.Height = height - m.heightMargin
}
