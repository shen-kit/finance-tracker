package frontend

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

type recordForm struct {
	form  *tview.Form
	iDate *tview.InputField
	iCat  *tview.DropDown
	iAmt  *tview.InputField
	iDesc *tview.TextArea
	tvMsg *tview.TextView
}

func createRecordsTable(monthGrid borderColorChanger) *updatableTable {
	table := newUpdatableTable(strings.Split("ID:Date:Category:Description:Amount", ":"), monthGrid)
	table.title = "Records"
	table.fGetMaxPage = backend.GetRecordsMaxPage
	return &table
}

func setRecTableKeybinds(t *updatableTable, rf recordForm) {
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if res := t.defaultInputCapture(event); res == nil {
			return nil
		}

		if event.Rune() == 'a' {
			showRecordForm(t, rf, -1, "", "", "", "")
		} else if event.Rune() == 'd' { // delete record
			row, _ := t.GetSelection()
			id := t.getCellInt(row, 0)
			showModal("Delete this record? (y/n)", func() {
				backend.DeleteRecord(id)
				t.update(t.fGetData(t.curPage))
				// set focus if deleted last row
				if row > t.GetRowCount()-1 {
					t.Select(max(0, row-1), 0)
				}
			}, t)
		} else if event.Rune() == 'e' { // edit record
			row, _ := t.GetSelection()
			id := t.getCellInt(row, 0)
			date := t.getCellString(row, 1)
			catName := t.getCellString(row, 2)
			desc := t.getCellString(row, 3)
			amt := t.getCellString(row, 4)
			showRecordForm(t, rf, id, date, desc, amt, catName)
		} else {
			return event
		}
		return nil
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
		AddButton("Cancel", nil).
		SetFieldBackgroundColor(tview.Styles.MoreContrastBackgroundColor).
		SetButtonBackgroundColor(tview.Styles.MoreContrastBackgroundColor)

	form.SetBorder(true).
		SetBorderColor(tview.Styles.TertiaryTextColor)

	return recordForm{
		form: form, iDate: inDate, iAmt: inAmt, iCat: inCat, iDesc: inDesc, tvMsg: formMsg,
	}
}

func showRecordForm(t updatablePrim, rf recordForm, id int, date, desc, amt, catName string) {

	/* ===== Helper Functions ===== */
	catOpt := 0
	setCategoryOptions := func() {
		cats := backend.GetCategories(0)
		catNames := make([]string, len(cats))
		for i, cat := range cats {
			catNames[i] = cat.SpreadToStrings()[1]
			if catNames[i] == catName {
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
		app.SetFocus(t)
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

		t.update(t.fGetData(t.getCurPage()))
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

	return backend.Record{Date: date, Amt: int(amt * 100), Desc: desc, CatId: catId}, nil
}
