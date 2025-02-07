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
	currentPage      int8 // currently viewed page
	lastPageIdx      int8

	newInvForm  *tview.Form
	inDate      *tview.InputField
	inCode      *tview.InputField
	inUnitprice *tview.InputField
	inQty       *tview.InputField
	editingId   int
)

func createInvestmentsTable() {
	investmentsTable = tview.NewTable().
		SetBorders(false).
		SetSeparator(tview.Borders.Vertical).
		SetFixed(1, 0).
		SetSelectable(true, false)

	investmentsTable.SetBorder(true).SetBorderPadding(1, 1, 2, 2)

	investmentsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'h' {
			flex.RemoveItem(investmentsTable)
			app.SetFocus(flex)
			return nil
		} else if event.Rune() == 'a' {
			editingId = -1
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

			editingId = int(id)
			showInvestmentForm(date, code, qty, unitprice)
		} else if event.Rune() == 'L' { // next page
			changePage(currentPage + 1)
		} else if event.Rune() == 'H' && currentPage > 0 { // previous page
			changePage(currentPage - 1)
		}
		return event
	})

	updateInvestmentsTable()
}

func changePage(page int8) {
	if page < 0 || page > lastPageIdx {
		return
	}
	currentPage = page
	investmentsTable.SetTitle(fmt.Sprintf("Investments (%d/%d)", currentPage+1, lastPageIdx+1))
	updateInvestmentsTable()
}

/* Pulls data from the backend to update the table */
func updateInvestmentsTable() {
	investmentsTable.Clear()
	lastPageIdx = backend.GetInvestmentsPages() - 1

	invs, err := backend.GetInvestmentsRecent(int(currentPage))
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
	changePage(0) // always start on Page 0
	flex.AddItem(investmentsTable, 0, 1, true)
	app.SetFocus(investmentsTable)
}

func createNewInvestmentForm() {

	closeForm := func() {
		flex.RemoveItem(newInvForm)
		app.SetFocus(investmentsTable)
	}

	onSubmit := func() {
		inv, success := parseForm()
		if !success {
			return
		}
		if editingId == -1 {
			backend.InsertInvestment(inv)
		} else {
			backend.UpdateInvestment(editingId, inv)
		}
		closeForm()
	}

	isPartialDate := func(s string, r rune) bool {
		regex0 := regexp.MustCompile(`^\d{0,4}$`)
		regex1 := regexp.MustCompile(`^\d{4}-\d{0,2}$`)
		regex2 := regexp.MustCompile(`^\d{4}-\d{2}-\d{0,2}$`)
		return regex0.MatchString(s) || regex1.MatchString(s) || regex2.MatchString(s)
	}

	inDate = tview.NewInputField().
		SetLabel("Date").
		SetFieldWidth(10).
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

	newInvForm = tview.NewForm().
		AddFormItem(inDate).
		AddFormItem(inCode).AddFormItem(inUnitprice).AddFormItem(inQty).
		AddButton("Cancel", closeForm).
		AddButton("Save", onSubmit)

	newInvForm.SetBorder(true)
}

/* Returns (Investment, success?) */
func parseForm() (backend.Investment, bool) {
	code := inCode.GetText()
	qty, err := strconv.ParseFloat(inQty.GetText(), 32)
	if err != nil {
		panic(err)
	}
	unitprice, err := strconv.ParseFloat(inUnitprice.GetText(), 32)
	if err != nil {
		panic(err)
	}
	date, err := time.Parse("2006-01-02", inDate.GetText())
	if err != nil {
		panic(err)
	}
	if code == "" || qty == 0 || unitprice <= 0 {
		newInvForm.SetLabelColor(tcell.ColorRed)
		return backend.Investment{}, false
	}

	return backend.Investment{Date: date, Code: code, Qty: float32(qty), Unitprice: float32(unitprice)}, true
}

/*
Set id=-1 if adding a new investment (other fields ignored).
Fill with details for editing an existing investment.
*/
func showInvestmentForm(date, code, unitprice, qty string) {
	if editingId == -1 {
		newInvForm.SetTitle("Add Investment")
	} else {
		newInvForm.SetTitle("Edit Investment Details")
	}
	inDate.SetText(date)
	inCode.SetText(code)
	inUnitprice.SetText(unitprice)
	inQty.SetText(qty)

	flex.AddItem(newInvForm, 55, 0, true)
	newInvForm.SetFocus(0)
	app.SetFocus(newInvForm)
}
