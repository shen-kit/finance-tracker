package frontend

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

func createInvSummaryTable() *updatableTable {
	table := newUpdatableTable(strings.Split("Code:Qty:Avg Buy Price:Current Price:Total In:Current Value:P/L:%P/L", ":"))
	table.title = "Investment Summary"
	table.fGetMaxPage = func() int { return 0 }
	return &table
}

func setInvSummaryTableKeybinds(t *updatableTable) {
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if res := t.defaultInputCapture(event); res == nil {
			return nil
		}
		return event
	})
}
