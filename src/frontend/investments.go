package frontend

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

var (
	investmentsTable *tview.Table
	invCurrentPage   int8 // currently viewed page
	invLastPage      int8

	invDetailsForm *tview.Form
	invInDate      *tview.InputField
	invInCode      *tview.InputField
	invInUnitprice *tview.InputField
	invInQty       *tview.InputField
	invEditingId   int
)

func createInvestmentsTable() {
	investmentsTable = tview.NewTable().
		SetBorders(false).
		SetSeparator(tview.Borders.Vertical).
		SetFixed(1, 0).
		SetSelectable(true, false)

	investmentsTable.SetBorder(true).SetBorderPadding(1, 1, 2, 2)

	investmentsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'h' || event.Key() == tcell.KeyCtrlC {
			flex.RemoveItem(investmentsTable)
			app.SetFocus(flex)
			return nil
		} else if event.Rune() == 'a' {
			invEditingId = -1
			showInvestmentForm("", "", "", "")
			return nil
		} else if event.Rune() == 'd' { // delete investment
			row, _ := investmentsTable.GetSelection()
			res, err := strconv.ParseInt(strings.Trim(investmentsTable.GetCell(row, 0).Text, " "), 10, 32)
			if err != nil {
				panic(err)
			}
			backend.DeleteInvestment(int(res))
			updateInvestmentsTable()
		} else if event.Rune() == 'e' { // edit investment
			row, _ := investmentsTable.GetSelection()
			id, _ := strconv.ParseInt(strings.Trim(investmentsTable.GetCell(row, 0).Text, " "), 10, 32)
			date := strings.Trim(investmentsTable.GetCell(row, 1).Text, " ")
			code := strings.Trim(investmentsTable.GetCell(row, 2).Text, " ")
			unitprice := strings.Trim(investmentsTable.GetCell(row, 3).Text, " $")
			qty := strings.Trim(investmentsTable.GetCell(row, 4).Text, " ")

			invEditingId = int(id)
			showInvestmentForm(date, code, qty, unitprice)
		} else if event.Rune() == 'L' { // next page
			invChangePage(invCurrentPage + 1)
		} else if event.Rune() == 'H' { // previous page
			invChangePage(invCurrentPage - 1)
		}
		return event
	})

	updateInvestmentsTable()
}

func invChangePage(page int8) {
	if page < 0 || page > invLastPage {
		return
	}
	invCurrentPage = page
	investmentsTable.SetTitle(fmt.Sprintf("Investments (%d/%d)", invCurrentPage+1, invLastPage+1))
	updateInvestmentsTable()
}

/* Pulls data from the backend to update the table */
func updateInvestmentsTable() {
	investmentsTable.Clear()
	invLastPage = backend.GetInvestmentsPages() - 1

	invs, err := backend.GetInvestmentsRecent(int(invCurrentPage))
	if err != nil {
		panic(err)
	}

	headers := strings.Split(" ID : Date : Code : Unitprice : Qty : Total ", ":")
	for i, h := range headers {
		investmentsTable.SetCell(0, i, tview.NewTableCell(h).SetSelectable(false).SetStyle(tcell.StyleDefault.Bold(true)))
	}

	for i, inv := range invs {
		id, date, code, qty, unitprice := inv.Spread()
		investmentsTable.
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

func showInvestmentsTable() {
	invChangePage(0) // always start on Page 0
	flex.AddItem(investmentsTable, 0, 1, true)
	app.SetFocus(investmentsTable)
}

// INVESTMENT FORM

func createInvestmentForm() {

	closeForm := func() {
		flex.RemoveItem(invDetailsForm)
		app.SetFocus(investmentsTable)
	}

	onSubmit := func() {
		inv, success := parseInvForm()
		if !success {
			return
		}
		if invEditingId == -1 {
			backend.InsertInvestment(inv)
		} else {
			backend.UpdateInvestment(invEditingId, inv)
		}
		closeForm()
	}

	invInDate = tview.NewInputField().
		SetLabel("Date").
		SetFieldWidth(10).
		SetPlaceholder("YYYY-MM-DD").
		SetAcceptanceFunc(isPartialDate)

	invInCode = tview.NewInputField().
		SetLabel("Stock Code").
		SetFieldWidth(10)

	invInUnitprice = tview.NewInputField().
		SetLabel("Unit Price").
		SetFieldWidth(7).
		SetAcceptanceFunc(tview.InputFieldFloat)

	invInQty = tview.NewInputField().
		SetLabel("Qty").
		SetFieldWidth(7).
		SetAcceptanceFunc(tview.InputFieldFloat)

	invDetailsForm = tview.NewForm().
		AddFormItem(invInDate).
		AddFormItem(invInCode).AddFormItem(invInUnitprice).AddFormItem(invInQty).
		AddButton("Cancel", closeForm).
		AddButton("Save", onSubmit)

	invDetailsForm.SetBorder(true)

	invDetailsForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			closeForm()
			return nil
		}
		return event
	})
}

/* Returns (Investment, success?) */
func parseInvForm() (backend.Investment, bool) {
	code := invInCode.GetText()
	qty, err := strconv.ParseFloat(invInQty.GetText(), 32)
	if err != nil {
		panic(err)
	}
	unitprice, err := strconv.ParseFloat(invInUnitprice.GetText(), 32)
	if err != nil {
		panic(err)
	}
	date, err := time.Parse("2006-01-02", invInDate.GetText())
	if err != nil {
		panic(err)
	}
	if code == "" || qty == 0 || unitprice <= 0 {
		invDetailsForm.SetLabelColor(tcell.ColorRed)
		return backend.Investment{}, false
	}

	return backend.Investment{Date: date, Code: code, Qty: float32(qty), Unitprice: float32(unitprice)}, true
}

func showInvestmentForm(date, code, unitprice, qty string) {
	if invEditingId == -1 {
		invDetailsForm.SetTitle("Add Investment")
	} else {
		invDetailsForm.SetTitle("Edit Investment Details")
	}
	invInDate.SetText(date)
	invInCode.SetText(code)
	invInUnitprice.SetText(unitprice)
	invInQty.SetText(qty)

	flex.AddItem(invDetailsForm, 55, 0, true)
	invDetailsForm.SetFocus(0)
	app.SetFocus(invDetailsForm)
}
