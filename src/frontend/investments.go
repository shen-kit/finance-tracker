package frontend

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

var (
	investmentsTable *tview.Table
	newInvForm       *tview.Form
)

func createInvestmentsTable() {
	investmentsTable = tview.NewTable().
		SetBorders(false).
		SetSeparator(tview.Borders.Vertical).
		SetFixed(1, 0).
		SetSelectable(true, false)

	investmentsTable.SetTitle("Investments").SetBorder(true).SetBorderPadding(1, 1, 2, 2)

	investmentsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'h' {
			flex.RemoveItem(investmentsTable)
			app.SetFocus(flex)
			return nil
		} else if event.Rune() == 'a' {
			showNewInvestmentForm()
			return nil
		}
		return event
	})
}

func updateInvestmentsTable() {
	investmentsTable.Clear()

	invs, err := backend.GetInvestmentsRecent(0)
	if err != nil {
		panic(err)
	}

	headers := strings.Split(" ID : Date : Code : Unitprice : Qty : Total ", ":")
	for i, h := range headers {
		investmentsTable.SetCell(0, i, tview.NewTableCell(h).SetSelectable(false).SetStyle(tcell.StyleDefault.Bold(true)))
	}

	for i, inv := range invs {
		id, date, code, qty, unitprice := inv.Spread()
		investmentsTable.
			SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf(" %d ", id)).
				SetAlign(tview.AlignCenter).
				SetMaxWidth(4)).
			SetCell(i+1, 1, tview.NewTableCell(date.Format("YYYY-MM-DD")).SetAlign(tview.AlignCenter).SetMaxWidth(12)).
			SetCell(i+1, 2, tview.NewTableCell(code)).
			SetCell(i+1, 3, tview.NewTableCell(fmt.Sprintf("%.1f", qty))).
			SetCell(i+1, 4, tview.NewTableCell(fmt.Sprintf("%.2f", unitprice))).
			SetCell(i+1, 5, tview.NewTableCell(fmt.Sprintf("%.2f", unitprice*qty)))
	}
}

func showInvestmentsTable() {
	updateInvestmentsTable()
	flex.AddItem(investmentsTable, 0, 1, true)
	app.SetFocus(investmentsTable)
}

func createNewInvestmentForm() {

}

func showNewInvestmentForm() {

}
