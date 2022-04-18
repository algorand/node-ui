package explorer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"

	"github.com/algorand/go-algorand-sdk/encoding/json"
	"github.com/algorand/go-algorand-sdk/types"
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

	middleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		b.Right = "├"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

func (m *Model) initTransaction(txn *types.SignedTxnInBlock) {
	m.txnView.YOffset = 0
	m.txnView.SetContent(indent.String(string(json.Encode(txn)), 6))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) headerView() string {
	info := middleStyle.Render(fmt.Sprintf("%3.f%%", m.txnView.ScrollPercent()*100))
	//title := titleStyle.Render(fmt.Sprintf("Txn: %s", m.txn.Txn.ID()))
	title := titleStyle.Render(fmt.Sprintf("Txn: %s", "TODO: Compute ID"))
	line := strings.Repeat("─", max(0, m.txnView.Width-lipgloss.Width(title)-lipgloss.Width(info)-1))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, "─", info, line)
}

func (m Model) footerView() string {
	line := strings.Repeat("─", max(0, m.txnView.Width))
	return lipgloss.JoinHorizontal(lipgloss.Center, line)
}

func (m Model) viewTransaction() string {
	return lipgloss.JoinVertical(0,
		m.headerView(),
		m.txnView.View(),
		m.footerView(),
	)
}
