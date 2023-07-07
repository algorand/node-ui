package accounts

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/go-algorand-sdk/v2/types"

	"github.com/algorand/node-ui/messages"
	"github.com/algorand/node-ui/tui/internal/style"
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

type balance struct {
	MicroAlgos uint64
	TimeStamp  time.Time
}

type account struct {
	Balances       map[uint64]uint64
	BalanceHistory []balance
}

func makeAccount() *account {
	return &account{
		Balances: make(map[uint64]uint64),
		BalanceHistory: []balance{
			{0, time.Now()},
			{0, time.Now()},
			{0, time.Now()},
		}}
}

// Model representing the account bubble.
type Model struct {
	accounts []types.Address
	Accounts map[types.Address]*account

	Err          error
	style        *style.Styles
	viewport     viewport.Model
	heightMargin int

	requestor *messages.Requestor
}

// New creates the accounts Model.
func New(style *style.Styles, requestor *messages.Requestor, initialHeight int, heightMargin int, accounts []types.Address) Model {
	rval := Model{
		Accounts:     make(map[types.Address]*account),
		style:        style,
		viewport:     viewport.New(0, 0),
		heightMargin: heightMargin,
		requestor:    requestor,
	}
	rval.setSize(80, initialHeight)
	rval.SetAccounts(accounts)
	return rval
}

// SetAccounts updates the accounts to monitor.
func (m *Model) SetAccounts(accounts []types.Address) {
	updated := make(map[types.Address]*account)
	for _, addr := range accounts {
		if acct, ok := m.Accounts[addr]; ok {
			updated[addr] = acct
		} else {
			updated[addr] = makeAccount()
		}
	}
	m.Accounts = updated
	m.accounts = accounts
}

func (m *Model) setSize(width, height int) {
	footerHeight := lipgloss.Height(m.footerView())
	m.viewport.Width = width
	m.viewport.Height = height - m.heightMargin - footerHeight
}

// Init is part of the tea.Model interface.
func (m Model) Init() tea.Cmd {
	return m.requestor.GetAccountStatusCmd(m.accounts)
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

	case messages.AccountStatusMsg:
		cmds = append(cmds,
			tea.Tick(5*time.Second, func(time.Time) tea.Msg {
				return m.requestor.GetAccountStatusCmd(m.accounts)()
			}),
		)

		for msgAddress, msgBalances := range msg.Balances {
			acct := m.Accounts[msgAddress]

			// Don't update if the balance didn't change
			if msgBalances[0] == acct.Balances[0] {
				break
			}

			newBalance := balance{
				MicroAlgos: msgBalances[0],
				TimeStamp:  time.Now(),
			}

			// Prepend the balance
			tmpList := append([]balance{newBalance}, acct.BalanceHistory...)
			if len(tmpList) > 3 {
				tmpList = tmpList[:3]
			}
			acct.BalanceHistory = tmpList
			acct.Balances = msgBalances

			m.Accounts[msgAddress] = acct
		}

		m.viewport.SetContent(m.buildString())
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View is part of the tea.Model interface.
func (m Model) View() string {

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n%s", m.viewport.View(), m.footerView()))
	return builder.String()
}

func (m Model) buildString() string {
	builder := strings.Builder{}

	keys := make([]string, 0, len(m.Accounts))
	for k := range m.Accounts {
		keys = append(keys, k.String())
	}
	sort.Strings(keys)

	for _, account := range keys {
		accountType, _ := types.DecodeAddress(account)
		v := m.Accounts[accountType]
		builder.WriteString(fmt.Sprintf("%s %s\n", m.style.AccountBoldText.Render("Account:"), m.style.AccountYellowText.Render(account)))

		algoStr := fmt.Sprintf("         %f Algos", float64(v.Balances[0])/1000000.0)
		builder.WriteString(m.style.AccountBlueText.Render(algoStr) + "\n")
		for _, a := range v.BalanceHistory {
			if a.MicroAlgos == 0 {
				builder.WriteString("\n")
			} else {
				pastStr := fmt.Sprintf("         %f Algos @ %s\n", float64(a.MicroAlgos)/1000000, a.TimeStamp.Format("2006-01-02 15:04:05.1234"))
				builder.WriteString(pastStr)
			}
		}

	}

	return m.style.Account.Render(builder.String())
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
