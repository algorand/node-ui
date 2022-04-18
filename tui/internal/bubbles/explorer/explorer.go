package explorer

import (
	"context"

	table "github.com/calyptia/go-bubble-table"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/messages"
	"github.com/algorand/node-ui/tui/internal/constants"
	"github.com/algorand/node-ui/tui/internal/style"
)

type state int

const (
	blockState = iota
	paysetState
	txnState
)

const initialBlocks = 25

type blocks []blockItem
type txnItems []transactionItem

type Model struct {
	state state

	width        int
	widthMargin  int
	height       int
	heightMargin int
	style        *style.Styles

	// for blocks page
	blocks blocks

	// cache for transactions page
	transactions txnItems

	// cache for txn page
	//txn transactions.SignedTxnInBlock
	txn []blocks

	table     table.Model
	txnView   viewport.Model
	requestor *messages.Requestor
}

func NewModel(styles *style.Styles, requestor *messages.Requestor, width, widthMargin, height, heightMargin int) Model {
	m := Model{
		state:        blockState,
		style:        styles,
		width:        width,
		widthMargin:  widthMargin,
		height:       height,
		heightMargin: heightMargin,
		requestor:    requestor,
	}
	m.initBlocks()
	return m
}

type BlocksMsg struct {
	blocks []blockItem
	err    error
}

func (m Model) InitBlocks() tea.Msg {
	status, err := m.requestor.Client.Status().Do(context.Background())
	if err != nil {
		return BlocksMsg{
			err: err,
		}
	}
	return m.getBlocks(status.LastRound-initialBlocks, status.LastRound)()
}

func (m *Model) getBlocks(first, last uint64) tea.Cmd {
	return func() tea.Msg {
		var result BlocksMsg
		for i := last; i >= first; i-- {
			block, err := m.requestor.Client.BlockRaw(i).Do(context.Background())
			if err != nil {
				result.err = err
				return result
			}
			result.blocks = append(result.blocks, blockItem{i, block})
		}
		return result
	}
}

func (m Model) Init() tea.Cmd {
	return m.InitBlocks
}

func (m Model) nextBlockCmd(round uint64) tea.Cmd {
	return func() tea.Msg {
		_, err := m.requestor.Client.StatusAfterBlock(round).Do(context.Background())
		if err != nil {
			return BlocksMsg{err: err}
		}
		blk, err := m.requestor.Client.BlockRaw(round).Do(context.Background())
		if err != nil {
			return BlocksMsg{err: err}
		}
		return BlocksMsg{
			blocks: []blockItem{
				{Round: round, Block: blk},
			},
		}
	}
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	verticalFrameSize := m.style.Bottom.GetVerticalFrameSize()
	m.table.SetSize(width-m.widthMargin, height-m.heightMargin-verticalFrameSize)
	m.txnView.Width = width - m.widthMargin
	m.txnView.Height = height - m.heightMargin - lipgloss.Height(m.headerView()) - lipgloss.Height(m.footerView())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var updateCmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case messages.StatusMsg:
		m.status = msg
	case tea.KeyMsg:
		// navigate into explorer views
		switch {
		case key.Matches(msg, constants.Keys.Forward):
			switch m.state {
			case blockState:
				// Select transactions.
				m.state = paysetState
				switch block := m.table.SelectedRow().(type) {
				case blockItem:
					// TODO: block to txn
					//m.transactions = make([]transactionItem, 0)
					//for _, txn := range block.Payset {
					//	t := txn
					//	m.transactions = append(m.transactions, transactionItem{&t})
					//}
					m.transactions = append(m.transactions, transactionItem{block.Block})
				}
				m.initTransactions()
			case paysetState:
				m.state = txnState
				switch txn := m.table.SelectedRow().(type) {
				case transactionItem:
					m.initTransaction(txn.TransactionInBlock)
				}
			}

		// navigate out of explorer views
		case key.Matches(msg, constants.Keys.Back):
			switch m.state {
			case paysetState:
				m.state = blockState
				//m.initBlocks()
				//return m, tea.Batch(append(cmds, m.getBlocks)...)
			case txnState:
				m.state = paysetState
			}
		}

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)

	case BlocksMsg:
		// append blocks
		backup := m.blocks
		m.blocks = msg.blocks
		for _, blk := range backup {
			m.blocks = append(m.blocks, blk)
		}
		cmds = append(cmds, m.nextBlockCmd(m.blocks[0].Round+1))
	}

	t, tableCmd := m.table.Update(msg)
	m.table = t
	cmds = append(cmds, tableCmd)

	switch m.state {
	case blockState:
		m, updateCmd = m.UpdateBlocks(msg)
		return m, tea.Batch(append(cmds, updateCmd)...)
	case paysetState:
		return m, nil
	case txnState:
		m.txnView, updateCmd = m.txnView.Update(msg)
		return m, tea.Batch(append(cmds, updateCmd)...)
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case blockState, paysetState:
		return m.style.Bottom.Render(m.table.View())
	case txnState:
		return m.viewTransaction()
	}
	return ""
}
