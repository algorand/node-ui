package style

import (
	"github.com/charmbracelet/lipgloss"
)

// For now, this is in its own package so that it can be shared between
// different packages without incurring an illegal import cycle.

const (
	// TopHeight is the hard coded height of the top bubbles.
	TopHeight = 13
)

// Styles defines styles for the TUI.
type Styles struct {
	ActiveBorderColor   lipgloss.Color
	InactiveBorderColor lipgloss.Color

	// Accounts area
	Account           lipgloss.Style
	AccountBoldText   lipgloss.Style
	AccountGrayText   lipgloss.Style
	AccountBlueText   lipgloss.Style
	AccountYellowText lipgloss.Style

	// Status area
	Status         lipgloss.Style
	StatusBoldText lipgloss.Style

	// Bottom area
	Bottom          lipgloss.Style
	BottomPaginator lipgloss.Style

	BottomListTitle        lipgloss.Style
	BottomListItemSelector lipgloss.Style
	BottomListItemActive   lipgloss.Style
	BottomListItemInactive lipgloss.Style
	BottomListItemKey      lipgloss.Style

	// Footer
	Footer       lipgloss.Style
	FooterLeft   lipgloss.Style
	FooterMiddle lipgloss.Style
	FooterRight  lipgloss.Style
}

// DefaultStyles returns default styles for the TUI.
func DefaultStyles() *Styles {
	s := new(Styles)

	// used
	s.ActiveBorderColor = lipgloss.Color("62")
	//s.InactiveBorderColor = lipgloss.Color("236")
	s.InactiveBorderColor = lipgloss.Color("#ABB8C3")
	s.BottomPaginator = lipgloss.NewStyle().
		Margin(0).
		Align(lipgloss.Center)

	// Accounts
	s.Account = lipgloss.NewStyle().
		//BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.InactiveBorderColor).
		Padding(0, 1, 0, 1).
		MarginLeft(1)
	s.AccountBoldText = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0693E3"))
	s.AccountGrayText = lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))
	s.AccountBlueText = lipgloss.NewStyle().Foreground(lipgloss.Color("#0693E3"))
	s.AccountYellowText = lipgloss.NewStyle().Foreground(lipgloss.Color("#A3A322"))

	// Status
	s.Status = lipgloss.NewStyle().
		Width(64).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.InactiveBorderColor).
		Padding(0, 1, 0, 1).
		MarginLeft(1)
	s.StatusBoldText = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0693E3"))

	// Bottom box
	s.Bottom = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.InactiveBorderColor).
		Padding(1, 2, 1, 1).
		MarginLeft(1)

	s.BottomListItemInactive = lipgloss.NewStyle().
		MarginLeft(1)

	s.BottomListTitle = lipgloss.NewStyle().
		//Align(lipgloss.Center). // did not work.
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	s.BottomListItemSelector = s.BottomListItemInactive.Copy().
		Width(1).
		Foreground(lipgloss.Color("#B083EA"))

	s.BottomListItemActive = s.BottomListItemInactive.Copy().
		Bold(true)

	s.BottomListItemKey = s.BottomListItemInactive.Copy().
		Width(10).
		Foreground(lipgloss.Color("#A3A322"))

	// Inspired by lipgloss demo
	s.Footer = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	s.FooterLeft = lipgloss.NewStyle().
		Inherit(s.Footer).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#4CAF50")).
		Padding(0, 1).
		MarginRight(1)
	s.FooterMiddle = lipgloss.NewStyle().
		Inherit(s.Footer)

	s.FooterRight = lipgloss.NewStyle().Inherit(s.Footer).
		Background(lipgloss.Color("#A550DF")).
		Padding(0, 1).
		Align(lipgloss.Right)

	return s
}
