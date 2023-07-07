package explorer

import (
	"fmt"
	"io"
	"strconv"

	table "github.com/calyptia/go-bubble-table"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

// transactionItem is used by the list bubble.
type transactionItem struct {
	*types.SignedTxnInBlock
}

func formatAmount(txn *types.SignedTxnInBlock) string {
	switch txn.Txn.Type {
	case types.PaymentTx:
		return fmt.Sprintf("%f", txn.Txn.Amount.ToAlgos())
	case types.AssetTransferTx:
		return strconv.FormatUint(txn.Txn.AssetTransferTxnFields.AssetAmount, 10)
	}
	return "-"
}

var transactionTableHeader = []string{"  INTRA", "type", "amount", "sigtype", "fee", "has-note", "sender"}

func computeTxnRow(b transactionItem) string {
	var sigtype string
	if !(b.Sig == types.Signature{}) {
		sigtype = "ed25519"
	} else if !b.Msig.Blank() {
		sigtype = "msig"
	} else if !b.Lsig.Blank() {
		sigtype = "lsig"
	} else {
		sigtype = "inner-txn"
	}

	return fmt.Sprintf("\t%s\t%s\t%s\t%f\t%t\t%s",
		b.Txn.Type,
		formatAmount(b.SignedTxnInBlock),
		sigtype,
		b.Txn.Fee.ToAlgos(),
		len(b.Txn.Note) > 0,
		b.Txn.Sender.String(),
	)
}

func (i transactionItem) Render(w io.Writer, model table.Model, index int) {
	var cursor string
	if index == model.Cursor() {
		cursor = "> "
	} else {
		cursor = "  "
	}

	cursor = activeStyle.Render(cursor)
	intra := keyStyle.Render(strconv.FormatUint(uint64(index), 10))
	rest := computeTxnRow(i)
	if index == model.Cursor() {
		rest = activeStyle.Render(rest)
	} else {
		rest = inactiveStyle.Render(rest)
	}
	fmt.Fprintf(w, "%s%s%s\n", cursor, intra, rest)
}

func (m *Model) updateTxnTable() {
	var rows []table.Row
	for _, t := range m.transactions {
		rows = append(rows, t)
	}

	m.table.SetRows(rows)
}

func (m *Model) initTransactions() {
	t := table.New(transactionTableHeader, 0, 0)
	t.KeyMap.Up.SetKeys(append(t.KeyMap.Up.Keys(), "k")...)
	t.KeyMap.Down.SetKeys(append(t.KeyMap.Down.Keys(), "j")...)
	t.Styles.Title = m.style.StatusBoldText
	m.table = t
	m.setSize(m.width, m.height)
	m.updateTxnTable()
}
