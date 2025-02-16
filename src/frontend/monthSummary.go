package frontend

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

func createMonthSummary(recTable *updatableTable) *monthGridView {

	tvTitle := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)
	tvTitle.SetBorderPadding(1, 1, 3, 3)

	tvSummary := tview.NewTextView()
	tvSummary.SetBorderPadding(0, 0, 3, 3)

	msGrid := tview.NewGrid().
		SetRows(3, 2, 0).
		SetBorders(true).
		AddItem(tvTitle, 0, 0, 1, 1, 0, 0, false).
		AddItem(tvSummary, 1, 0, 1, 1, 0, 0, false).
		AddItem(recTable, 2, 0, 1, 1, 0, 0, true)

	msGrid.SetBorder(true).SetTitle("Month Summary")

	return &monthGridView{
		Grid:      msGrid,
		tvTitle:   tvTitle,
		tvSummary: tvSummary,
		table:     recTable,
	}
}

func setMonthGridKeybinds(monthGrid *monthGridView, rf recordForm) {
	monthGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if isBackKey(event) {
			gotoHomepage()
		} else if event.Rune() == 'a' { // add record
			showRecordForm(monthGrid, rf, -1, "", "", "", "")
		} else if event.Rune() == 'd' { // delete record
			row, _ := monthGrid.table.GetSelection()
			id := monthGrid.table.getCellInt(row, 0)
			backend.DeleteRecord(id)
		} else if event.Rune() == 'e' { // edit record
			row, _ := monthGrid.table.GetSelection()
			id := monthGrid.table.getCellInt(row, 0)
			date := monthGrid.table.getCellString(row, 1)
			catName := monthGrid.table.getCellString(row, 2)
			desc := monthGrid.table.getCellString(row, 3)
			amt := monthGrid.table.getCellString(row, 4)
			showRecordForm(monthGrid, rf, id, date, desc, amt, catName)
		} else if event.Rune() == 'H' {
			monthGrid.changeMonth(-1)
		} else if event.Rune() == 'L' {
			monthGrid.changeMonth(1)
		} else {
			return event
		}
		return nil
	})
}

func showMonthSummary(monthView *monthGridView) {
	monthView.SetBorder(false)
	monthView.update(monthView.fGetData(monthView.getCurPage()))
	flex.AddItem(monthView, 0, 1, true)
	app.SetFocus(monthView)
}

func (mv monthGridView) update(recs []backend.DataRow) {
	t := time.Now().AddDate(0, mv.monthOffset, 0)
	_, income, expenditure := backend.GetMonthInfo(t)

	// set title text
	mv.tvTitle.SetText(fmt.Sprintf("%s %d", t.Month().String(), t.Year()))

	// set summary text
	incomeStr := fmt.Sprintf("$%.0f", income)
	expenditureStr := fmt.Sprintf("$%.0f", expenditure)
	mv.tvSummary.SetText(fmt.Sprintf("Income:      %8s\nExpenditure: %8s", incomeStr, expenditureStr))

	// update table data
	mv.table.update(recs)
}

func (mv monthGridView) reset() {
	mv.changeMonth(-mv.monthOffset)
	mv.table.SetBorder(false)
}
