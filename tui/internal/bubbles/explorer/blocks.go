package explorer

import (
	"fmt"
	"io"
	"strconv"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
	table "github.com/calyptia/go-bubble-table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BlockItem is used by the list bubble.
type BlockItem struct {
	Round uint64
	Block models.BlockResponse
}

// Hacked these in to workaround missing style options in table model
var inactiveStyle = lipgloss.NewStyle()
var activeStyle = inactiveStyle.Copy().Foreground(lipgloss.Color("#B083EA")).Bold(true)
var keyStyle = inactiveStyle.Copy().Width(10).Foreground(lipgloss.Color("#A3A322")).Bold(true)

var blockTableHeader = []string{"  ROUND", "Txns", "Pay", "[Sum λ]", "Axfer", "Acfg", "Afrz", "[Unique]", "Appl", "[Unique]", "Proposer"}

func proposer(cert *map[string]interface{}) string {
	if cert == nil {
		return "<unknown>"
	}

	// There must be a better way than this...
	prop, ok := (*cert)["prop"]
	if ok {
		switch p := prop.(type) {
		case map[interface{}]interface{}:
			for k, v := range p {
				switch c := k.(type) {
				case string:
					if c == "oprop" {
						switch c2 := v.(type) {
						case []byte:
							var addr types.Address
							copy(addr[:], c2)
							return addr.String()
						}
					}
				}
			}
		}
	}

	return "<unknown>"
}

func computeBlockRow(b BlockItem) string {
	block := b.Block.Block

	typeCount := make(map[types.TxType]uint)
	var paymentsTotal uint64
	assets := make(map[uint64]struct{})
	apps := make(map[uint64]struct{})
	for _, tx := range block.Payset {
		typeCount[tx.Txn.Type]++

		switch tx.Txn.Type {

		case types.PaymentTx:
			paymentsTotal += uint64(tx.Txn.PaymentTxnFields.Amount)
		case types.ApplicationCallTx:
			id := uint64(tx.Txn.ApplicationCallTxnFields.ApplicationID)
			if id == 0 {
				id = tx.ApplyData.ApplicationID
			}
			if id == 0 {
				break
			}
			if _, ok := apps[id]; !ok {
				apps[id] = struct{}{}
			}
		case types.AssetTransferTx:
			fallthrough
		case types.AssetFreezeTx:
			fallthrough
		case types.AssetConfigTx:
			id := uint64(tx.Txn.AssetTransferTxnFields.XferAsset)
			if id == 0 {
				id = tx.ApplyData.ConfigAsset
			}
			if id == 0 {
				id = uint64(tx.Txn.AssetConfigTxnFields.ConfigAsset)
			}
			if id == 0 {
				id = uint64(tx.Txn.AssetFreezeTxnFields.FreezeAsset)
			}
			if id == 0 {
				break
			}
			if _, ok := assets[id]; !ok {
				assets[id] = struct{}{}
			}
		}
	}

	return fmt.Sprintf("\t%d\t%d\t%f\t%d\t%d\t%d\t%d\t%d\t%d\t%s",
		len(b.Block.Block.Payset),
		typeCount[types.PaymentTx],
		float64(paymentsTotal)/float64(10000),
		typeCount[types.AssetTransferTx],
		typeCount[types.AssetConfigTx],
		typeCount[types.AssetFreezeTx],
		len(assets),
		typeCount[types.ApplicationCallTx],
		len(apps),
		proposer(b.Block.Cert))
}

// Render implements the Row interface to display a row of data.
func (i BlockItem) Render(w io.Writer, model table.Model, index int) {
	var cursor string
	if index == model.Cursor() {
		cursor = "> "
	} else {
		cursor = "  "
	}

	cursor = activeStyle.Render(cursor)
	//round := keyStyle.Render(strconv.FormatUint(uint64(i.Block.Round()), 10))
	round := keyStyle.Render(strconv.FormatUint(i.Round, 10))
	rest := computeBlockRow(i)
	if index == model.Cursor() {
		rest = activeStyle.Render(rest)
	} else {
		rest = inactiveStyle.Render(rest)
	}
	fmt.Fprintf(w, "%s%s%s\n", cursor, round, rest)
}

func (m *Model) updateBlockTable() {
	if len(m.blocks) <= 0 {
		return
	}

	var rows []table.Row
	for _, b := range m.blocks {
		rows = append(rows, b)
	}

	m.table.SetRows(rows)
}

func (m *Model) initBlocks() {
	t := table.New(blockTableHeader, 0, 0)
	t.KeyMap.Up.SetKeys(append(t.KeyMap.Up.Keys(), "k")...)
	t.KeyMap.Down.SetKeys(append(t.KeyMap.Down.Keys(), "j")...)
	t.Styles.Title = m.style.StatusBoldText
	m.table = t
	m.setSize(m.width, m.height)
	m.updateBlockTable()
}

// updateBlocks mimics the tea.Model update function.
func (m Model) updateBlocks(msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case BlocksMsg:
		if m.state == blockState {
			m.updateBlockTable()
		}
	}

	return m, nil
}
