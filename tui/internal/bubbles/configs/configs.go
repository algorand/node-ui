package configs

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/messages"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

// Model representing the configs page.
type Model struct {
	heightMargin int
	viewport     viewport.Model
}

// New creates a Model.
func New(heightMargin int) Model {
	m := Model{
		viewport:     viewport.New(0, 0),
		heightMargin: heightMargin,
	}
	m.setSize(80, 20)
	return m
}

// ConfigContent allows the update function to find its config content.
type ConfigContent string

func (m Model) getContent() tea.Cmd {
	return func() tea.Msg {
		return ConfigContent(messages.GetConfigs())
	}
}

// Init is part of the tea.Model interface.
func (m Model) Init() tea.Cmd {
	return m.getContent()
}

func (m *Model) setSize(width, height int) {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())

	m.viewport.Width = width
	m.viewport.Height = height - m.heightMargin - headerHeight - footerHeight
}

// Update is part of the tea.Model interface.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case ConfigContent:
		// For some reason tabs make the viewport go crazy
		//m.viewport.SetContent(string(msg))
		m.viewport.SetContent(strings.ReplaceAll(string(msg), "\t", "    "))

	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View is part of the tea.Model interface.
func (m Model) View() string {

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView()))
	return builder.String()
}

func (m Model) headerView() string {
	title := titleStyle.Render("Node configurations")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
