package footer

import (
	"github.com/algorand/node-ui/messages"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/tui/internal/style"
)

// Model for the footer.
type Model struct {
	width  int
	height int
	style  *style.Styles

	network messages.NetworkMsg
}

// New creates the footer Model.
func New(s *style.Styles) Model {
	return Model{style: s}
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

	case messages.NetworkMsg:
		m.network = msg
	}

	return m, nil
}

// View is part of the tea.Model interface.
func (m Model) View() string {

	left := m.style.FooterLeft.Render("Algorand Node UI")
	//right := m.style.FooterRight.Render(config.GetAlgorandVersion())
	right := m.style.FooterRight.Render(m.network.NodeVersion)
	//middleText := fmt.Sprintf("%s (Gensis Hash %s)", m.network.GenesisID, m.network.GenesisHash)
	middleText := m.network.GenesisID

	middle := m.style.FooterMiddle.Copy().
		Width(m.width - lipgloss.Width(left) - lipgloss.Width(right)).
		Render(middleText)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		left,
		middle,
		right,
	)
}
