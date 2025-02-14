package frontend

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

type monthView struct {
	grid        *tview.Grid
	tvTitle     *tview.TextView
	tvIncome    *tview.TextView
	tableView   *tableView
	monthOffset int
}

func createMonthSummary(recTv *tableView, rf recordForm) *monthView {

	table := recTv.table

	tvTitle := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)
	tvTitle.SetBorderPadding(1, 1, 3, 3)

	tvIncome := tview.NewTextView()
	tvIncome.SetBorderPadding(0, 0, 3, 3)

	msGrid := tview.NewGrid().
		SetRows(3, 2, 0).
		SetBorders(true).
		AddItem(tvTitle, 0, 0, 1, 1, 0, 0, false).
		AddItem(tvIncome, 1, 0, 1, 1, 0, 0, false).
		AddItem(table, 2, 0, 1, 1, 0, 0, true)

	msGrid.SetBorder(true).SetTitle("Month Summary")

	mv := &monthView{
		grid:      msGrid,
		tvTitle:   tvTitle,
		tvIncome:  tvIncome,
		tableView: recTv,
	}

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
		} else if event.Rune() == 'H' {
			changeMonth(mv, -1)
		} else if event.Rune() == 'L' {
			changeMonth(mv, 1)
		} else {
			return event
		}
		return nil
	})

	return mv
}

func showMonthSummary(monthView *monthView) {
	monthView.tableView.table.SetBorder(false)
	updateMonthSummary(monthView)
	flex.AddItem(monthView.grid, 0, 1, true)
	app.SetFocus(monthView.tableView.table)
}

func updateMonthSummary(monthView *monthView) {

	t := time.Now().AddDate(0, monthView.monthOffset, 0)
	recs, income, expenditure := backend.GetMonthInfo(t)

	// set text views
	monthView.tvTitle.SetText(fmt.Sprintf("%s %d", t.Month().String(), t.Year()))

	incomeStr := fmt.Sprintf("$%.0f", income)
	expenditureStr := fmt.Sprintf("$%.0f", expenditure)
	monthView.tvIncome.SetText(fmt.Sprintf("Income:      %8s\nExpenditure: %8s", incomeStr, expenditureStr))

	// update table
	createTableHeaders(monthView.tableView)
	for i, rec := range recs {
		id, date, desc, amt, catId := rec.Spread()
		catName := backend.GetCategoryNameFromId(catId)
		monthView.tableView.table.
			SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf(" %d ", id)).
				SetAlign(tview.AlignCenter)).
			SetCell(i+1, 1, tview.NewTableCell(date.Format(" 2006-01-02 ")).
				SetAlign(tview.AlignCenter)).
			SetCell(i+1, 2, tview.NewTableCell(" "+catName+" ")).
			SetCell(i+1, 3, tview.NewTableCell(" "+desc+" ").SetExpansion(1)).
			SetCell(i+1, 4, tview.NewTableCell(fmt.Sprintf(" $%.2f ", amt)))
	}
}

func changeMonth(monthView *monthView, offset int) {
	monthView.monthOffset += offset
	updateMonthSummary(monthView)
}
