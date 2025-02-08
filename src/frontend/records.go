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

	recDetailsForm *tview.Form
	recInDate      *tview.InputField
	recInDesc      *tview.TextArea
	recInAmt       *tview.InputField
	recInCat       *tview.DropDown

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
			showRecordsForm("", "", "", 0)
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
			showRecordsForm("", "", "", 0)
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
	recLastPage = backend.GetRecordsPages() - 1

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
		catName := backend.GetCategoryName(catId)
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
	closeForm := func() {
		flex.RemoveItem(recDetailsForm)
		app.SetFocus(recordsTable)
	}

	onSubmit := func() {
		rec, success := parseRecForm()
		if !success {
			return
		}
		if recEditingId == -1 {
			backend.InsertRecord(rec)
		} else {
			backend.UpdateRecord(recEditingId, rec)
		}
		closeForm()
	}

	recInDate = tview.NewInputField().
		SetLabel("Date").
		SetFieldWidth(10).
		SetPlaceholder("YYYY-MM-DD").
		SetAcceptanceFunc(isPartialDate)

	recInCat = tview.NewDropDown()

	recInAmt = tview.NewInputField().
		SetLabel("Amount").
		SetFieldWidth(7).
		SetAcceptanceFunc(tview.InputFieldFloat)

	recInDesc = tview.NewTextArea().
		SetLabel("Description").
		SetSize(30, 3)

	recDetailsForm = tview.NewForm().
		AddFormItem(recInDate).
		AddFormItem(recInCat).
		AddFormItem(recInAmt).
		AddFormItem(recInDesc).
		AddButton("Cancel", closeForm).
		AddButton("Save", onSubmit)

}

func parseRecForm() (backend.Record, bool) {
	panic("unimplemented")
}

func showRecordsForm(date, desc, amt string, catId int) {

	cats, err := backend.GetCategories()
	if err != nil {
		panic(err)
	}

	catNames := make([]string, len(cats))
	catOpt := 0
	for i, cat := range cats {
		catNames[i] = cat.Name
		if cat.Id == catId {
			catOpt = i
		}
	}
	recInCat.SetOptions(catNames, nil)

	recInDate.SetText(date)
	recInDesc.SetText(desc, true)
	recInAmt.SetText(amt)
	recInCat.SetCurrentOption(catOpt)

	flex.AddItem(recDetailsForm, 55, 0, true)
	recDetailsForm.SetFocus(0)
	app.SetFocus(recDetailsForm)
}
