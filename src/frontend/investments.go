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

// global function to update the investments table using the tableView object
var updateInvestmentsTable func()

type investmentForm struct {
	form       *tview.Form
	iDate      *tview.InputField
	iCode      *tview.InputField
	iQty       *tview.InputField
	iUnitprice *tview.InputField
	tvMsg      *tview.TextView
}

func createInvestmentsTable() *tableView {

	/* ===== Helper Functions ===== */

	// returns a closure with the tableView saved
	createUpdateTableClosure := func(tv *tableView) func() {
		return func() {
			createTableHeaders(tv)
			tv.maxPage = backend.GetInvestmentsPages() - 1
			invs := backend.GetInvestmentsRecent(tv.curPage)

			for i, inv := range invs {
				id, date, code, qty, unitprice := inv.Spread()
				tv.table.
					SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf(" %d ", id)).
						SetAlign(tview.AlignCenter).
						SetMaxWidth(4)).
					SetCell(i+1, 1, tview.NewTableCell(date.Format(" 2006-01-02 ")).SetAlign(tview.AlignCenter).SetMaxWidth(12)).
					SetCell(i+1, 2, tview.NewTableCell(" "+code+" ")).
					SetCell(i+1, 3, tview.NewTableCell(fmt.Sprintf(" $%.2f ", unitprice))).
					SetCell(i+1, 4, tview.NewTableCell(fmt.Sprintf(" %.0f ", qty))).
					SetCell(i+1, 5, tview.NewTableCell(fmt.Sprintf(" $%.2f ", unitprice*qty)))
			}
		}
	}

	/* ===== Function Body ===== */

	table := createMyTable()
	tv := &tableView{
		table:   table,
		title:   "Investments",
		headers: strings.Split(" ID : Date : Code : Unitprice : Qty : Total ", ":"),
	}
	updateInvestmentsTable = createUpdateTableClosure(tv)
	tv.fUpdate = updateInvestmentsTable
	return tv
}

func setInvestmentTableKeybinds(tv *tableView, inf investmentForm) {
	table := tv.table
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'h' || event.Key() == tcell.KeyCtrlC {
			tv.hide(flex)
			return nil
		} else if event.Rune() == 'a' {
			showInvestmentForm(table, inf, -1, "", "", "", "")
			return nil
		} else if event.Rune() == 'd' { // delete investment
			row, _ := table.GetSelection()
			id := table.getCellInt(row, 0)
			backend.DeleteInvestment(id)
			tv.fUpdate()
			return nil
		} else if event.Rune() == 'e' { // edit investment
			row, _ := table.GetSelection()
			id := table.getCellInt(row, 0)
			date := table.getCellString(row, 1)
			code := table.getCellString(row, 2)
			unitprice := table.getCellString(row, 3)
			qty := table.getCellString(row, 4)
			showInvestmentForm(table, inf, id, date, code, qty, unitprice)
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

func showInvestmentForm(lastWidget tview.Primitive, inf investmentForm, id int, date, code, unitprice, qty string) {

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
		app.SetFocus(lastWidget)
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

		updateInvestmentsTable()
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
			Date: date, Code: code, Qty: float32(qty), Unitprice: float32(unitprice)},
		nil
}
