package model

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/algorand/go-algorand-sdk/types"

	"github.com/algorand/node-ui/messages"
	"github.com/algorand/node-ui/tui/internal/bubbles/about"
	"github.com/algorand/node-ui/tui/internal/bubbles/accounts"
	"github.com/algorand/node-ui/tui/internal/bubbles/configs"
	"github.com/algorand/node-ui/tui/internal/bubbles/explorer"
	"github.com/algorand/node-ui/tui/internal/bubbles/footer"
	"github.com/algorand/node-ui/tui/internal/bubbles/status"
	"github.com/algorand/node-ui/tui/internal/bubbles/tabs"
	"github.com/algorand/node-ui/tui/internal/style"
)

const (
	initialWidth  = 80
	initialHeight = 50
)

type activeComponent int

const (
	explorerTab activeComponent = iota
	utilitiesTab
	accountTab
	configTab
	helpTab
)

type Model struct {
	Status        tea.Model
	Accounts      tea.Model
	Tabs          tabs.Model
	BlockExplorer tea.Model
	Configs       tea.Model
	Utilities     tea.Model
	About         tea.Model
	Help          help.Model
	Footer        tea.Model

	network messages.NetworkMsg

	styles *style.Styles

	requestor *messages.Requestor

	active activeComponent
	// remember the last resize so we can re-send it when selecting a different bottom component.
	lastResize tea.WindowSizeMsg
}

func New(requestor *messages.Requestor, addresses []types.Address) Model {
	styles := style.DefaultStyles()
	tab := tabs.New([]string{"EXPLORER", "UTILITIES", "ACCOUNTS", "CONFIGURATION", "HELP"})
	// The tab content is the only flexible element.
	// This means the height must grow or shrink to fill the available
	// window height. It has access to the absolute height but needs to
	// be informed about the space used by other elements.
	tabContentMargin := style.TopHeight + tab.Height() + 2 /* +2 for footer/help */
	return Model{
		active:        explorerTab,
		styles:        styles,
		Status:        status.New(styles, requestor),
		Tabs:          tab,
		BlockExplorer: explorer.NewModel(styles, requestor, initialWidth, 0, initialHeight, tabContentMargin),
		Configs:       configs.New(tabContentMargin),
		Accounts:      accounts.NewModel(styles, requestor, initialHeight, tabContentMargin, addresses),
		Help:          help.New(),
		Footer:        footer.New(styles),
		About:         about.New(tabContentMargin, about.GetHelpContent()),
		Utilities:     about.New(tabContentMargin, about.GetUtilsContent()),
		requestor:     requestor,
	}
}
