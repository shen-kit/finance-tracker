package frontend

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

var (
	investmentsTable *tview.Table
	newInvForm       *tview.Form
	editingId        int

	dateInput      *tview.InputField
	codeInput      *tview.InputField
	unitpriceInput *tview.InputField
	qtyInput       *tview.InputField
)

func createInvestmentsTable() {
	investmentsTable = tview.NewTable().
		SetBorders(false).
		SetSeparator(tview.Borders.Vertical).
		SetFixed(1, 0).
		SetSelectable(true, false)

	investmentsTable.SetTitle("Investments").SetBorder(true).SetBorderPadding(1, 1, 2, 2)

	investmentsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'h' {
			flex.RemoveItem(investmentsTable)
			app.SetFocus(flex)
			return nil
		} else if event.Rune() == 'a' {
			showNewInvestmentForm()
			return nil
		} else if event.Rune() == 'd' { // delete currently hovered investment
			row, _ := investmentsTable.GetSelection()
			res, err := strconv.ParseInt(strings.Trim(investmentsTable.GetCell(row, 0).Text, " "), 10, 32)
			if err != nil {
				panic(err)
			}
			backend.DeleteInvestment(int(res))
			updateInvestmentsTable()
		} else if event.Rune() == 'e' { // edit currently hovered investment
			row, _ := investmentsTable.GetSelection()
			id, _ := strconv.ParseInt(strings.Trim(investmentsTable.GetCell(row, 0).Text, " "), 10, 32)
			date := strings.Trim(investmentsTable.GetCell(row, 1).Text, " ")
			code := strings.Trim(investmentsTable.GetCell(row, 2).Text, " ")
			unitprice := strings.Trim(investmentsTable.GetCell(row, 3).Text, " $")
			qty := strings.Trim(investmentsTable.GetCell(row, 4).Text, " ")

			showEditInvestmentForm(int(id), date, code, qty, unitprice)
		}
		return event
	})
}

func updateInvestmentsTable() {
	investmentsTable.Clear()

	invs, err := backend.GetInvestmentsRecent(0)
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
	updateInvestmentsTable()
	flex.AddItem(investmentsTable, 0, 1, true)
	app.SetFocus(investmentsTable)
}

func createNewInvestmentForm() {

	isPartialDate := func(s string, r rune) bool {
		regex0 := regexp.MustCompile(`^\d{0,4}$`)
		regex1 := regexp.MustCompile(`^\d{4}-\d{0,2}$`)
		regex2 := regexp.MustCompile(`^\d{4}-\d{2}-\d{0,2}$`)
		return regex0.MatchString(s) || regex1.MatchString(s) || regex2.MatchString(s)
	}

	closeForm := func() {
		flex.RemoveItem(newInvForm)
		app.SetFocus(investmentsTable)
	}

	onSubmit := func() {

		code := codeInput.GetText()
		qty, err := strconv.ParseFloat(qtyInput.GetText(), 32)
		if err != nil {
			panic(err)
		}
		unitprice, err := strconv.ParseFloat(unitpriceInput.GetText(), 32)
		if err != nil {
			panic(err)
		}

		date, err := time.Parse("2006-01-02", dateInput.GetText())
		if err != nil || code == "" || qty == 0 || unitprice <= 0 {
			newInvForm.SetLabelColor(tcell.ColorRed)
			return
		}

		inv := backend.Investment{Date: date, Code: code, Qty: float32(qty), Unitprice: float32(unitprice)}
		if editingId == -1 {
			backend.InsertInvestment(inv)
		} else {
			backend.UpdateInvestment(editingId, inv)
			editingId = -1
		}

		updateInvestmentsTable()
		closeForm()
	}

	dateInput = tview.NewInputField().
		SetLabel("Date").
		SetFieldWidth(10).
		SetPlaceholder("YYYY-MM-DD").
		SetAcceptanceFunc(isPartialDate)

	codeInput = tview.NewInputField().
		SetLabel("Stock Code").
		SetFieldWidth(10)

	unitpriceInput = tview.NewInputField().
		SetLabel("Unit Price").
		SetFieldWidth(7)

	qtyInput = tview.NewInputField().
		SetLabel("Qty").
		SetFieldWidth(7)

	newInvForm = tview.NewForm().
		AddFormItem(dateInput).
		AddFormItem(codeInput).AddFormItem(unitpriceInput).AddFormItem(qtyInput).
		AddButton("Save", onSubmit).
		AddButton("Cancel", closeForm)

	newInvForm.SetTitle("Investment Record").SetBorder(true)
}

func showNewInvestmentForm() {
	dateInput.SetText("")
	codeInput.SetText("")
	unitpriceInput.SetText("")
	qtyInput.SetText("")
	flex.AddItem(newInvForm, 55, 0, true)
	app.SetFocus(newInvForm)
}

func showEditInvestmentForm(id int, date, code, qty, unitprice string) {
	editingId = id

	dateInput.SetText(date)
	codeInput.SetText(code)
	unitpriceInput.SetText(unitprice)
	qtyInput.SetText(qty)

	flex.AddItem(newInvForm, 55, 0, true)
	app.SetFocus(newInvForm)
}
