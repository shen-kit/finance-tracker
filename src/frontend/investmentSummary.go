package frontend

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createInvSummaryTable() *updatableTable {
	table := newUpdatableTable(strings.Split("Code:Qty:Avg Buy Price:Current Price:Total In:Current Value:P/L:%P/L", ":"), nil)
	table.title = "Investment Summary"
	table.fGetMaxPage = func() int { return 0 }
	table.SetBlurFunc(func() { table.SetBorderColor(tview.Styles.BorderColor) })
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
