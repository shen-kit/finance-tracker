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

// global function to update the records table using the tableView object
var updateRecordsTable func()

type recordForm struct {
	form  *tview.Form
	iDate *tview.InputField
	iCat  *tview.DropDown
	iAmt  *tview.InputField
	iDesc *tview.TextArea
	tvMsg *tview.TextView
}

func createRecordsTable() *tableView {

	/* ===== Helper Functions ===== */

	// returns a closure with the tableView saved
	createUpdateTableClosure := func(tv *tableView) func() {
		return func() {
			createTableHeaders(tv)
			tv.maxPage = backend.GetRecordsPages() - 1
			recs := backend.GetRecordsRecent(tv.curPage)

			for i, rec := range recs {
				id, date, desc, amt, catId := rec.Spread()
				catName := backend.GetCategoryNameFromId(catId)
				tv.table.
					SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf(" %d ", id)).
						SetAlign(tview.AlignCenter)).
					SetCell(i+1, 1, tview.NewTableCell(date.Format(" 2006-01-02 ")).
						SetAlign(tview.AlignCenter)).
					SetCell(i+1, 2, tview.NewTableCell(" "+catName+" ")).
					SetCell(i+1, 3, tview.NewTableCell(" "+desc+" ")).
					SetCell(i+1, 4, tview.NewTableCell(fmt.Sprintf(" $%.2f ", amt)))
			}
		}
	}

	/* ===== Function Body ===== */
	table := createMyTable()
	tv := &tableView{
		table:   table,
		title:   "Records",
		headers: strings.Split(" ID : Date : Category : Description : Amount ", ":"),
	}
	updateRecordsTable = createUpdateTableClosure(tv)
	tv.fUpdate = updateRecordsTable
	return tv
}

func setRecordTableKeybinds(tv *tableView, rf recordForm) {
	table := tv.table
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if isBackKey(event) {
			tv.hide(flex)
			return nil
		} else if event.Rune() == 'a' {
			showRecordForm(table, rf, -1, "", "", "", "")
			return nil
		} else if event.Rune() == 'd' { // delete record
			row, _ := table.GetSelection()
			id := table.getCellInt(row, 0)
			backend.DeleteRecord(id)
			tv.fUpdate()
			return nil
		} else if event.Rune() == 'e' { // edit record
			row, _ := table.GetSelection()
			id := table.getCellInt(row, 0)
			date := table.getCellString(row, 1)
			catName := table.getCellString(row, 2)
			desc := table.getCellString(row, 3)
			amt := table.getCellString(row, 4)
			showRecordForm(table, rf, id, date, desc, amt, catName)
			return nil
		} else if event.Rune() == 'L' { // next page
			changePage(tv, tv.curPage+1)
			return nil
		} else if event.Rune() == 'H' { // previous page
			changePage(tv, tv.curPage-1)
			return nil
		}
		return event
	})
}

func createRecordForm() recordForm {
	var form *tview.Form
	var inDate, inAmt *tview.InputField
	var inDesc *tview.TextArea
	var inCat *tview.DropDown
	var formMsg *tview.TextView

	inDate = tview.NewInputField().
		SetLabel("Date").
		SetFieldWidth(11).
		SetPlaceholder("YYYY-MM-DD").
		SetAcceptanceFunc(isPartialDate)

	inCat = tview.NewDropDown().
		SetLabel("Category")

	inCat.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' || event.Key() == tcell.KeyCtrlN {
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		} else if event.Rune() == 'k' || event.Key() == tcell.KeyCtrlP {
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		}
		return event
	})

	inAmt = tview.NewInputField().
		SetLabel("Amount").
		SetFieldWidth(7).
		SetAcceptanceFunc(tview.InputFieldFloat)

	inDesc = tview.NewTextArea().
		SetLabel("Description").
		SetSize(4, 35)

	formMsg = tview.NewTextView().
		SetSize(1, 35).
		SetDynamicColors(true).
		SetScrollable(false)

	form = tview.NewForm().
		AddFormItem(inDate).
		AddFormItem(inCat).
		AddFormItem(inAmt).
		AddFormItem(inDesc).
		AddFormItem(formMsg).
		AddButton("Save", nil).
		AddButton("Cancel", nil)

	form.SetBorder(true)

	return recordForm{
		form: form, iDate: inDate, iAmt: inAmt, iCat: inCat, iDesc: inDesc, tvMsg: formMsg,
	}
}

func showRecordForm(lastWidget tview.Primitive, rf recordForm, id int, date, desc, amt, catName string) {

	/* ===== Helper Functions ===== */
	catOpt := 0
	setCategoryOptions := func() {
		cats := backend.GetCategories()
		catNames := make([]string, len(cats))
		for i, cat := range cats {
			catNames[i] = cat.Name
			if cat.Name == catName {
				catOpt = i
			}
		}
		rf.iCat.SetOptions(catNames, nil)
	}

	setInputFieldValues := func() {
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}
		rf.iDate.SetText(date)
		rf.iDesc.SetText(desc, true)
		rf.iAmt.SetText(amt)
		rf.iCat.SetCurrentOption(catOpt)
		rf.tvMsg.SetText("")
	}

	closeForm := func() {
		flex.RemoveItem(rf.form)
		app.SetFocus(lastWidget)
	}

	onSubmit := func() {
		rec, err := parseRecForm(rf)
		if err != nil {
			rf.tvMsg.SetText("[red]" + err.Error())
			return
		}

		if id == -1 {
			backend.InsertRecord(rec)
		} else {
			backend.UpdateRecord(id, rec)
		}

		updateRecordsTable()
		closeForm()
	}

	/* ===== Function Body ===== */

	// set title
	if id == -1 {
		rf.form.SetTitle("Add Record")
	} else {
		rf.form.SetTitle("Edit Record Details")
	}

	setCategoryOptions()
	setInputFieldValues()

	rf.form.SetInputCapture(formInputCapture(closeForm, onSubmit))
	rf.form.GetButton(rf.form.GetButtonIndex("Cancel")).SetSelectedFunc(closeForm)
	rf.form.GetButton(rf.form.GetButtonIndex("Save")).SetSelectedFunc(onSubmit)

	// display + focus form
	flex.AddItem(rf.form, 55, 0, true)
	rf.form.SetFocus(0)
	app.SetFocus(rf.form)
}

/* Takes input from the form and returns a Record object */
func parseRecForm(rf recordForm) (backend.Record, error) {

	fail := func(msg string) (backend.Record, error) {
		return backend.Record{}, errors.New(msg)
	}

	if rf.iDate.GetText() == "" || rf.iAmt.GetText() == "" || rf.iDesc.GetText() == "" {
		return fail("All fields are required")
	}

	date, err := time.Parse("2006-01-02", rf.iDate.GetText())
	if err != nil {
		return fail("Date musy be in YYYY-MM-DD format")
	}

	_, cname := rf.iCat.GetCurrentOption()
	if cname == "" {
		return fail("Please choose a category")
	}
	catId := backend.GetCategoryIdFromName(cname)

	desc := rf.iDesc.GetText()

	amt, err := strconv.ParseFloat(rf.iAmt.GetText(), 32)
	if err != nil || amt == 0 {
		return fail("Invalid amount entered")
	}

	return backend.Record{Date: date, Amt: float32(amt), Desc: desc, CatId: catId}, nil
}
