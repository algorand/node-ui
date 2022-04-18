package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	//activeStyle = lipgloss.NewStyle().
	//		Background(lipgloss.Color("#527772")).
	//		Foreground(lipgloss.Color("#6dd588"))
	//inactiveStyle = lipgloss.NewStyle()
	activeColorForeground = lipgloss.Color("#6dd588")
	activeColorBackground = lipgloss.Color("#527772")

	defaultWidth = 20

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│  ",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}
	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 2)
	activeTab = tab.Copy().
			Border(activeTabBorder, true).
			Background(activeColorBackground).
			Foreground(activeColorForeground)

	tabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	//docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

type Model struct {
	width int

	index int
	tabs  []string

	//tabWidth int

	ActiveStyle   lipgloss.Style
	InactiveStyle lipgloss.Style
}

func New(tabs []string) Model {
	//max := 0
	//for _, t := range tabs {
	//	if len(t) > max {
	//		max = len(t)
	//	}
	//}
	return Model{
		width: 80,
		tabs:  tabs,
		//tabWidth:      max,
		//ActiveStyle:   activeStyle,
		//InactiveStyle: inactiveStyle,
	}
}

func (m Model) Height() int {
	return 3
}
func (m *Model) SetActiveIndex(i int) {
	m.index = i
}

func (m Model) GetActiveIndex() int {
	return m.index
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) View() string {
	doc := strings.Builder{}

	// Tabs
	{
		var renderedTabs []string
		renderedTabs = append(renderedTabs, "\n"+tabGap.Render(strings.Repeat(" ", 5)))

		// Activate the correct tab
		for i, t := range m.tabs {
			if i == m.index {
				renderedTabs = append(renderedTabs, activeTab.Render(t))
			} else {
				renderedTabs = append(renderedTabs, tab.Render(t))
			}
		}

		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			renderedTabs...,
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row))))

		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		doc.WriteString(row)
	}

	return doc.String()
}
