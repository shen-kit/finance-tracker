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
		} else if event.Rune() == 'd' {
			row, _ := investmentsTable.GetSelection()
			res, err := strconv.ParseInt(investmentsTable.GetCell(row, 0).Text, 10, 32)
			if err != nil {
				panic(err)
			}
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
			SetCell(i+1, 3, tview.NewTableCell(fmt.Sprintf(" %.0f ", qty))).
			SetCell(i+1, 4, tview.NewTableCell(fmt.Sprintf(" $%.2f ", unitprice))).
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

	var iCode string
	var iQty, iUnitprice float32
	var dateInput *tview.InputField

	closeForm := func() {
		flex.RemoveItem(newInvForm)
		app.SetFocus(investmentsTable)
	}

	onSubmit := func() {
		iDate, err := time.Parse("2006-01-02", dateInput.GetText())
		if err != nil || iCode == "" || iQty == 0 || iUnitprice <= 0 {
			newInvForm.SetLabelColor(tcell.ColorRed)
			return
		}

		backend.InsertInvestment(backend.Investment{Date: iDate, Code: iCode, Qty: iQty, Unitprice: iUnitprice})

		updateInvestmentsTable()
		closeForm()
	}

	dateInput = tview.NewInputField().
		SetLabel("Date").
		SetFieldWidth(10).
		SetPlaceholder("YYYY-MM-DD").
		SetAcceptanceFunc(isPartialDate)

	newInvForm = tview.NewForm().
		AddFormItem(dateInput).
		AddInputField("Stock Code", "", 7, nil, func(v string) { iCode = v }).
		AddInputField("Unitprice", "", 7, tview.InputFieldFloat, func(v string) {
			res, _ := strconv.ParseFloat(v, 32)
			iUnitprice = float32(res)
		}).
		AddInputField("Qty", "", 7, tview.InputFieldFloat, func(v string) {
			res, _ := strconv.ParseFloat(v, 32)
			iQty = float32(res)
		}).
		AddButton("Add", onSubmit).
		AddButton("Cancel", closeForm)

	newInvForm.SetTitle("Investment Record").SetBorder(true)
}

func showNewInvestmentForm() {
	flex.AddItem(newInvForm, 55, 0, true)
	app.SetFocus(newInvForm)
}
