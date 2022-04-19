package footer

import (
	"fmt"
	"github.com/algorand/node-ui/messages"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/tui/internal/style"
)

type Model struct {
	width  int
	height int
	style  *style.Styles

	network messages.NetworkMsg
}

func New(s *style.Styles) Model {
	return Model{style: s}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case messages.NetworkMsg:
		m.network = msg
	}

	return m, nil
}

func (m Model) View() string {
	left := m.style.FooterLeft.Render("Algorand Node UI")
	//right := m.style.FooterRight.Render(config.GetAlgorandVersion())
	right := m.style.FooterRight.Render("TODO: Lookup Version")
	//middleText := fmt.Sprintf("%s (Gensis Hash %s)", m.network.GenesisID, m.network.GenesisHash)
	middleText := fmt.Sprintf("%s", m.network.GenesisID)

	middle := m.style.FooterMiddle.Copy().
		Width(m.width - lipgloss.Width(left) - lipgloss.Width(right)).
		Render(middleText)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		left,
		middle,
		right,
	)
}
