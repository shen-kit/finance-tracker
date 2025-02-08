package frontend

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

var (
	recordsTable   *tview.Table
	recCurrentPage int8
	recLastPage    int8

	newRecForm *tview.Form
	recInDate  *tview.InputField
	recInDesc  *tview.InputField
	recInAmt   *tview.InputField
	recInCatId *tview.InputField

	recEditingId int
)

func createRecordsTable() {
	recordsTable = tview.NewTable().
		SetBorders(false).
		SetSeparator(tview.Borders.Vertical).
		SetFixed(1, 0).
		SetSelectable(true, false)

	recordsTable.SetBorder(true).SetBorderPadding(1, 1, 2, 2)

	recordsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'h' || event.Key() == tcell.KeyCtrlC {
			flex.RemoveItem(recordsTable)
			app.SetFocus(flex)
			return nil
		} else if event.Rune() == 'a' {
			recEditingId = -1
			showRecordsForm()
			return nil
		} else if event.Rune() == 'd' { // delete investment
			row, _ := recordsTable.GetSelection()
			res, err := strconv.ParseInt(strings.Trim(recordsTable.GetCell(row, 0).Text, " "), 10, 32)
			if err != nil {
				panic(err)
			}
			backend.DeleteRecord(int(res))
			updateRecordsTable()
		} else if event.Rune() == 'e' { // edit investment
			row, _ := recordsTable.GetSelection()
			id, _ := strconv.ParseInt(strings.Trim(recordsTable.GetCell(row, 0).Text, " "), 10, 32)
			// date := strings.Trim(recordsTable.GetCell(row, 1).Text, " ")
			// code := strings.Trim(recordsTable.GetCell(row, 2).Text, " ")
			// unitprice := strings.Trim(recordsTable.GetCell(row, 3).Text, " $")
			// qty := strings.Trim(recordsTable.GetCell(row, 4).Text, " ")

			recEditingId = int(id)
			showRecordsForm()
		} else if event.Rune() == 'L' { // next page
			recChangePage(recCurrentPage + 1)
		} else if event.Rune() == 'H' { // previous page
			recChangePage(recCurrentPage - 1)
		}
		return event
	})
}

func recChangePage(page int8) {
	if page < 0 || page > recLastPage {
		return
	}
	recCurrentPage = page
	recordsTable.SetTitle(fmt.Sprintf("Records (%d/%d)", recCurrentPage+1, recLastPage+1))
	updateRecordsTable()
}

func updateRecordsTable() {
	recordsTable.Clear()
	recLastPage = 3 // TODO: update with backend

	recs, err := backend.GetRecordsRecent(int(recCurrentPage))
	if err != nil {
		panic(err)
	}

	headers := strings.Split(" ID : Date : Category : Description : Amount ", ":")
	for i, h := range headers {
		recordsTable.SetCell(0, i, tview.NewTableCell(h).SetSelectable(false).SetStyle(tcell.StyleDefault.Bold(true)))
	}

	for i, rec := range recs {
		id, date, desc, amt, catId := rec.Spread()
		catName := fmt.Sprintf(" category %d ", catId) // TODO: fetch from backend
		recordsTable.
			SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf(" %d ", id)).
				SetAlign(tview.AlignCenter).
				SetMaxWidth(4)).
			SetCell(i+1, 1, tview.NewTableCell(date.Format(" 2006-01-02 ")).
				SetAlign(tview.AlignCenter).
				SetMaxWidth(12)).
			SetCell(i+1, 2, tview.NewTableCell(catName)).
			SetCell(i+1, 3, tview.NewTableCell(" "+desc+" ")).
			SetCell(i+1, 4, tview.NewTableCell(fmt.Sprintf(" $%.2f ", amt)))
	}
}

func showRecordsTable() {
	recChangePage(0)
	flex.AddItem(recordsTable, 0, 1, true)
	app.SetFocus(recordsTable)
}

func createRecordForm() {

}

func showRecordsForm() {
	panic("unimplemented")
}
