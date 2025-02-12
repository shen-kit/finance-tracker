package frontend

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

func createMonthSummary(recTv *tableView, rf recordForm) *tview.Grid {

	table := recTv.table
	msGrid := tview.NewGrid().
		SetRows(3, 3, 0).
		SetBorders(true).
		AddItem(table, 2, 0, 1, 1, 0, 0, true)

	msGrid.SetBorder(true).SetTitle("Month Summary")

	msGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if isBackKey(event) {
			flex.RemoveItem(msGrid)
			app.SetFocus(flex)
		} else if event.Rune() == 'a' { // add record
			showRecordForm(table, rf, -1, "", "", "", "")
		} else if event.Rune() == 'd' { // delete record
			row, _ := table.GetSelection()
			id := table.getCellInt(row, 0)
			backend.DeleteRecord(id)
		} else if event.Rune() == 'e' { // edit record
			row, _ := table.GetSelection()
			id := table.getCellInt(row, 0)
			date := table.getCellString(row, 1)
			catName := table.getCellString(row, 2)
			desc := table.getCellString(row, 3)
			amt := table.getCellString(row, 4)
			showRecordForm(table, rf, id, date, desc, amt, catName)
		} else {
			return event
		}
		return nil
	})

	return msGrid
}

func showMonthSummary(grid *tview.Grid, recTv *tableView) {
	updateMonthSummary(grid, recTv)
	flex.AddItem(grid, 0, 1, true)
	app.SetFocus(recTv.table)
}

func updateMonthSummary(grid *tview.Grid, recTv *tableView) {
	// update table
	createTableHeaders(recTv)
	recs := backend.GetRecordsMonth(time.Now())

	for i, rec := range recs {
		id, date, desc, amt, catId := rec.Spread()
		catName := backend.GetCategoryNameFromId(catId)
		recTv.table.
			SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf(" %d ", id)).
				SetAlign(tview.AlignCenter)).
			SetCell(i+1, 1, tview.NewTableCell(date.Format(" 2006-01-02 ")).
				SetAlign(tview.AlignCenter)).
			SetCell(i+1, 2, tview.NewTableCell(" "+catName+" ")).
			SetCell(i+1, 3, tview.NewTableCell(" "+desc+" ")).
			SetCell(i+1, 4, tview.NewTableCell(fmt.Sprintf(" $%.2f ", amt)))
	}
}
