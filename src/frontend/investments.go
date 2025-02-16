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

type investmentForm struct {
	form       *tview.Form
	iDate      *tview.InputField
	iCode      *tview.InputField
	iQty       *tview.InputField
	iUnitprice *tview.InputField
	tvMsg      *tview.TextView
}

func createInvestmentsTable() *updatableTable {
	table := newUpdatableTable(strings.Split("ID:Date:Code:Unitprice:Qty:Total", ":"))
	table.title = "Investments"
	// table.fGetData = backend.GetInvestmentsRecent
	return &table
}

func setInvTableKeybinds(t *updatableTable, inf investmentForm) {
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if res := t.defaultInputCapture(event); res == nil {
			return nil
		}

		if event.Rune() == 'a' {
			showInvestmentForm(t, inf, -1, "", "", "", "")
		} else if event.Rune() == 'd' { // delete investment
			row, _ := t.GetSelection()
			id := t.getCellInt(row, 0)
			backend.DeleteInvestment(id)
			t.update(t.fGetData(t.curPage))
		} else if event.Rune() == 'e' { // edit investment
			row, _ := t.GetSelection()
			id := t.getCellInt(row, 0)
			date := t.getCellString(row, 1)
			code := t.getCellString(row, 2)
			unitprice := t.getCellString(row, 3)
			qty := t.getCellString(row, 4)
			showInvestmentForm(t, inf, id, date, code, qty, unitprice)
		} else {
			return event
		}
		return nil
	})
}

func createInvestmentForm() investmentForm {

	var form *tview.Form
	var inDate, inCode, inUnitprice, inQty *tview.InputField
	var formMsg *tview.TextView

	inDate = tview.NewInputField().
		SetLabel("Date").
		SetFieldWidth(11).
		SetPlaceholder("YYYY-MM-DD").
		SetAcceptanceFunc(isPartialDate)

	inCode = tview.NewInputField().
		SetLabel("Stock Code").
		SetFieldWidth(10)

	inUnitprice = tview.NewInputField().
		SetLabel("Unit Price").
		SetFieldWidth(7).
		SetAcceptanceFunc(tview.InputFieldFloat)

	inQty = tview.NewInputField().
		SetLabel("Qty").
		SetFieldWidth(7).
		SetAcceptanceFunc(tview.InputFieldFloat)

	formMsg = tview.NewTextView().
		SetSize(1, 35).
		SetDynamicColors(true).
		SetScrollable(false)

	form = tview.NewForm().
		AddFormItem(inDate).
		AddFormItem(inCode).AddFormItem(inUnitprice).AddFormItem(inQty).
		AddFormItem(formMsg).
		AddButton("Save", nil).
		AddButton("Cancel", nil)

	form.SetBorder(true)

	return investmentForm{
		form: form, iDate: inDate, iCode: inCode, iUnitprice: inUnitprice, iQty: inQty, tvMsg: formMsg,
	}
}

func showInvestmentForm(t *updatableTable, inf investmentForm, id int, date, code, unitprice, qty string) {

	/* ===== Helper Functions ===== */

	setInputFieldValues := func() {
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}
		inf.iDate.SetText(date)
		inf.iCode.SetText(code)
		inf.iQty.SetText(qty)
		inf.iUnitprice.SetText(unitprice)
		inf.tvMsg.SetText("")
	}

	closeForm := func() {
		flex.RemoveItem(inf.form)
		app.SetFocus(t)
	}

	onSubmit := func() {
		inv, err := parseInvForm(inf)
		if err != nil {
			inf.tvMsg.SetText("[red]" + err.Error())
			return
		}

		if id == -1 {
			backend.InsertInvestment(inv)
		} else {
			backend.UpdateInvestment(id, inv)
		}

		t.update(t.fGetData(t.curPage))
		closeForm()
	}

	/* ===== Function Body ===== */

	if id == -1 {
		inf.form.SetTitle("Add Investment")
	} else {
		inf.form.SetTitle("Edit Investment Details")
	}

	setInputFieldValues()

	inf.form.SetInputCapture(formInputCapture(closeForm, onSubmit))
	inf.form.GetButton(inf.form.GetButtonIndex("Cancel")).SetSelectedFunc(closeForm)
	inf.form.GetButton(inf.form.GetButtonIndex("Save")).SetSelectedFunc(onSubmit)

	flex.AddItem(inf.form, 55, 0, true)
	inf.form.SetFocus(0)
	app.SetFocus(inf.form)
}

/* Returns (Investment, success?) */
func parseInvForm(inf investmentForm) (backend.Investment, error) {

	fail := func(msg string) (backend.Investment, error) {
		return backend.Investment{}, errors.New(msg)
	}

	for _, field := range []*tview.InputField{inf.iCode, inf.iDate, inf.iQty, inf.iUnitprice} {
		if field.GetText() == "" {
			return fail("All fields are required")
		}
	}

	code := inf.iCode.GetText()

	qty, err := strconv.ParseFloat(inf.iQty.GetText(), 32)
	if err != nil || qty == 0 {
		return fail("Quantity is invalid")
	}

	unitprice, err := strconv.ParseFloat(inf.iUnitprice.GetText(), 32)
	if err != nil || unitprice <= 0 {
		return fail("Unitprice is invalid")
	}

	date, err := time.Parse("2006-01-02", inf.iDate.GetText())
	if err != nil {
		return fail("Date must be in YYYY-MM-DD format")
	}

	return backend.Investment{
			Date: date, Code: code, Qty: float32(qty), Unitprice: int(unitprice * 100)},
		nil
}
