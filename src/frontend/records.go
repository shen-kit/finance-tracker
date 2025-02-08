package frontend

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	recEditingId   int
	recFormMsg     *tview.TextView
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
			showRecordsForm(-1, "", "", "", "")
			return nil
		} else if event.Rune() == 'd' { // delete record
			row, _ := recordsTable.GetSelection()
			res, err := strconv.ParseInt(strings.TrimSpace(recordsTable.GetCell(row, 0).Text), 10, 32)
			if err != nil {
				panic(err)
			}
			backend.DeleteRecord(int(res))
			updateRecordsTable()
		} else if event.Rune() == 'e' { // edit record
			row, _ := recordsTable.GetSelection()
			id, _ := strconv.ParseInt(strings.TrimSpace(recordsTable.GetCell(row, 0).Text), 10, 32)
			date := strings.TrimSpace(recordsTable.GetCell(row, 1).Text)
			catName := strings.TrimSpace(recordsTable.GetCell(row, 2).Text)
			desc := strings.TrimSpace(recordsTable.GetCell(row, 3).Text)
			amt := strings.Trim(recordsTable.GetCell(row, 4).Text, " $")

			showRecordsForm(int(id), date, desc, amt, catName)
			return nil
		} else if event.Rune() == 'L' { // next page
			recChangePage(recCurrentPage + 1)
			return nil
		} else if event.Rune() == 'H' { // previous page
			recChangePage(recCurrentPage - 1)
			return nil
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
		catName := backend.GetCategoryNameFromId(catId)
		recordsTable.
			SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf(" %d ", id)).
				SetAlign(tview.AlignCenter).
				SetMaxWidth(4)).
			SetCell(i+1, 1, tview.NewTableCell(date.Format(" 2006-01-02 ")).
				SetAlign(tview.AlignCenter).
				SetMaxWidth(12)).
			SetCell(i+1, 2, tview.NewTableCell(" "+catName+" ")).
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
		if flex.GetItemCount() == 1 { // if Add Record directly from homepage
			app.SetFocus(flex)
		} else {
			app.SetFocus(recordsTable)
		}
	}

	onSubmit := func() {
		rec, err := parseRecForm()
		if err != nil {
			recFormMsg.SetText("[red]" + err.Error())
			return
		}

		if recEditingId == -1 {
			backend.InsertRecord(rec)
		} else {
			backend.UpdateRecord(recEditingId, rec)
		}
		updateRecordsTable()
		closeForm()
	}

	recInDate = tview.NewInputField().
		SetLabel("Date").
		SetFieldWidth(11).
		SetPlaceholder("YYYY-MM-DD").
		SetAcceptanceFunc(isPartialDate)

	recInCat = tview.NewDropDown().
		SetLabel("Category")

	recInCat.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' || event.Key() == tcell.KeyCtrlN {
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		} else if event.Rune() == 'k' || event.Key() == tcell.KeyCtrlP {
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		}
		return event
	})

	recInAmt = tview.NewInputField().
		SetLabel("Amount").
		SetFieldWidth(7).
		SetAcceptanceFunc(tview.InputFieldFloat)

	recInDesc = tview.NewTextArea().
		SetLabel("Description").
		SetSize(4, 35)

	recFormMsg = tview.NewTextView().
		SetSize(1, 35).
		SetDynamicColors(true).
		SetScrollable(false)

	recDetailsForm = tview.NewForm().
		AddFormItem(recInDate).
		AddFormItem(recInCat).
		AddFormItem(recInAmt).
		AddFormItem(recInDesc).
		AddFormItem(recFormMsg).
		AddButton("Save", onSubmit).
		AddButton("Cancel", closeForm)

	recDetailsForm.SetBorder(true)

	recDetailsForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			closeForm()
			return nil
		}
		return event
	})

}

func parseRecForm() (backend.Record, error) {

	fail := func(msg string) (backend.Record, error) {
		return backend.Record{}, errors.New(msg)
	}

	if recInDate.GetText() == "" || recInAmt.GetText() == "" || recInDesc.GetText() == "" {
		return fail("All fields are required")
	}

	date, err := time.Parse("2006-01-02", recInDate.GetText())
	if err != nil {
		return fail("Date musy be in YYYY-MM-DD format")
	}

	_, cname := recInCat.GetCurrentOption()
	if cname == "" {
		return fail("Please choose a category")
	}
	catId := backend.GetCategoryIdFromName(cname)

	desc := recInDesc.GetText()

	amt, err := strconv.ParseFloat(recInAmt.GetText(), 32)
	if err != nil || amt == 0 {
		return fail("Invalid amount entered")
	}

	return backend.Record{Date: date, Amt: float32(amt), Desc: desc, CatId: catId}, nil
}

func showRecordsForm(id int, date, desc, amt, catName string) {

	recEditingId = id
	if id == -1 {
		recDetailsForm.SetTitle("Add Record")
	} else {
		recDetailsForm.SetTitle("Edit Record Details")
	}

	cats, err := backend.GetCategories()
	if err != nil {
		panic(err)
	}

	catNames := make([]string, len(cats))
	catOpt := 0
	for i, cat := range cats {
		catNames[i] = cat.Name
		if cat.Name == catName {
			catOpt = i
		}
	}
	recInCat.SetOptions(catNames, nil)

	recInDate.SetText(date)
	recInDesc.SetText(desc, true)
	recInAmt.SetText(amt)
	recInCat.SetCurrentOption(catOpt)
	recFormMsg.SetText("")

	flex.AddItem(recDetailsForm, 55, 0, true)
	recDetailsForm.SetFocus(0)
	app.SetFocus(recDetailsForm)
}
