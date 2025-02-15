package frontend

//
// import (
// 	"fmt"
// 	"time"
//
// 	"github.com/gdamore/tcell/v2"
// 	"github.com/rivo/tview"
// 	"github.com/shen-kit/finance-tracker/backend"
// )
//
// func createMonthSummary(recTable *myTable, rf recordForm) *monthView {
//
// 	tvTitle := tview.NewTextView().
// 		SetTextAlign(tview.AlignCenter)
// 	tvTitle.SetBorderPadding(1, 1, 3, 3)
//
// 	tvSummary := tview.NewTextView()
// 	tvSummary.SetBorderPadding(0, 0, 3, 3)
//
// 	msGrid := tview.NewGrid().
// 		SetRows(3, 2, 0).
// 		SetBorders(true).
// 		AddItem(tvTitle, 0, 0, 1, 1, 0, 0, false).
// 		AddItem(tvSummary, 1, 0, 1, 1, 0, 0, false).
// 		AddItem(recTable, 2, 0, 1, 1, 0, 0, true)
//
// 	msGrid.SetBorder(true).SetTitle("Month Summary")
//
// 	mv := &monthView{
// 		grid:      msGrid,
// 		tvTitle:   tvTitle,
// 		tvSummary: tvSummary,
// 	}
//
// 	msGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		if isBackKey(event) {
// 			flex.RemoveItem(msGrid)
// 			app.SetFocus(flex)
// 		} else if event.Rune() == 'a' { // add record
// 			showRecordForm(mv, rf, -1, "", "", "", "")
// 		} else if event.Rune() == 'd' { // delete record
// 			row, _ := mv.GetSelection()
// 			id := mv.getCellInt(row, 0)
// 			backend.DeleteRecord(id)
// 		} else if event.Rune() == 'e' { // edit record
// 			row, _ := mv.GetSelection()
// 			id := mv.getCellInt(row, 0)
// 			date := mv.getCellString(row, 1)
// 			catName := mv.getCellString(row, 2)
// 			desc := mv.getCellString(row, 3)
// 			amt := mv.getCellString(row, 4)
// 			showRecordForm(mv, rf, id, date, desc, amt, catName)
// 		} else if event.Rune() == 'H' {
// 			mv.changePage(-1)
// 		} else if event.Rune() == 'L' {
// 			mv.changePage(1)
// 		} else {
// 			return event
// 		}
// 		return nil
// 	})
//
// 	return mv
// }
//
// func showMonthSummary(monthView *monthView) {
// 	monthView.SetBorder(false)
// 	monthView.update()
// 	flex.AddItem(monthView.grid, 0, 1, true)
// 	app.SetFocus(monthView)
// }
//
// func (monthView monthView) update() {
//
// 	t := time.Now().AddDate(0, monthView.monthOffset, 0)
// 	recs, income, expenditure := backend.GetMonthInfo(t)
//
// 	// set text views
// 	monthView.tvTitle.SetText(fmt.Sprintf("%s %d", t.Month().String(), t.Year()))
//
// 	incomeStr := fmt.Sprintf("$%.0f", income)
// 	expenditureStr := fmt.Sprintf("$%.0f", expenditure)
// 	monthView.tvSummary.SetText(fmt.Sprintf("Income:      %8s\nExpenditure: %8s", incomeStr, expenditureStr))
//
// 	// update table
// 	monthView.createHeaders()
// 	for i, rec := range recs {
// 		id, date, desc, amt, catId := rec.Spread()
// 		catName := backend.GetCategoryNameFromId(catId)
// 		monthView.
// 			SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf(" %d ", id)).
// 				SetAlign(tview.AlignCenter)).
// 			SetCell(i+1, 1, tview.NewTableCell(date.Format(" 2006-01-02 ")).
// 				SetAlign(tview.AlignCenter)).
// 			SetCell(i+1, 2, tview.NewTableCell(" "+catName+" ")).
// 			SetCell(i+1, 3, tview.NewTableCell(" "+desc+" ").SetExpansion(1)).
// 			SetCell(i+1, 4, tview.NewTableCell(fmt.Sprintf(" $%.2f ", amt)))
// 	}
// }
