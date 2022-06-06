package model

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"

	"github.com/algorand/node-ui/tui/internal/constants"
)

// TODO: this function could implement a type and be passed to the tab view.
func (m Model) tabView() string {
	switch activeComponent(m.Tabs.GetActiveIndex()) {
	case explorerTab:
		return m.BlockExplorer.View()
	case accountTab:
		return m.Accounts.View()
	case configTab:
		return m.Configs.View()
	case helpTab:
		return m.About.View()
	case utilitiesTab:
		return m.Utilities.View()
	}

	return "unknown tab"
}

// View is part of the tea.Model interface.
func art() string {
	// TODO: This could take a width and indent/border to line up with the bottom
	art := `
                ▒█████ 
              ▒████████▒
            ▒████████████▓▓▓▓▒
           ▒█████▒▓████████▓
          ▒█████    ██████▓
         ▒▓████     ▒█████▓
        ▒█████     ▒███████
      ▒█████▓     ▓████████▓
     ▒█████▓     ▓███████████
    ▒█████▒     ██████ ▒█████▓
   ▒█████▒     █████▓   ██████▒
  ▒█████▒     ▒█████▓    ▒██████▒`
	return indent.String(art, 3)
}

// View is part of the tea.Model interface.
func (m Model) View() string {
	// Compose the different views by joining them together in the right orientation.
	return lipgloss.JoinVertical(0,
		lipgloss.JoinHorizontal(0,
			m.Status.View(),
			art()),
		m.Tabs.View(),
		m.tabView(),
		m.Help.View(constants.Keys),
		m.Footer.View())
}
