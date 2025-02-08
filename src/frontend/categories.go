package frontend

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

var (
	categoriesTable *tview.Table
	catDetailsForm  *tview.Form
)

func createCategoriesTable() {
	categoriesTable = tview.NewTable().
		SetBorders(false).
		SetSeparator(tview.Borders.Vertical).
		SetFixed(1, 0).            // always show header row
		SetSelectable(true, false) // rows selectable

	categoriesTable.SetTitle("Categories").SetBorder(true).SetBorderPadding(1, 1, 2, 2)

	categoriesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'h' || event.Key() == tcell.KeyCtrlC {
			flex.RemoveItem(categoriesTable)
			app.SetFocus(flex)
			return nil
		} else if event.Rune() == 'a' {
			showNewCategoryForm()
			return nil
		}
		return event
	})
}

func updateCategoriesTable() {
	categoriesTable.Clear()

	cats, err := backend.GetCategories()
	if err != nil {
		panic(err)
	}

	headers := strings.Split(" ID : Name : Type : Description ", ":")
	for i, h := range headers {
		categoriesTable.SetCell(0, i, tview.NewTableCell(h).SetSelectable(false).SetStyle(tcell.StyleDefault.Bold(true)))
	}

	for i, cat := range cats {
		id, name, isincome, desc := cat.Spread()
		categoriesTable.
			SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf(" %d ", id)).SetAlign(tview.AlignCenter).SetMaxWidth(4)).
			SetCell(i+1, 1, tview.NewTableCell(" "+name+" ").SetMaxWidth(20)).
			SetCell(i+1, 2, tview.NewTableCell(func() string {
				if isincome {
					return " Income "
				}
				return " Expenditure "
			}()).SetMaxWidth(15)).
			SetCell(i+1, 3, tview.NewTableCell(" "+desc+" "))
	}
}

func showCategoriesTable() {
	updateCategoriesTable()
	flex.AddItem(categoriesTable, 0, 1, true)
	app.SetFocus(categoriesTable)
}

func createNewCategoryForm() {

	var iName, iDesc string
	var iIncome bool

	closeForm := func() {
		flex.RemoveItem(catDetailsForm)
		app.SetFocus(categoriesTable)
	}

	onSubmit := func() {
		if iName == "" {
			catDetailsForm.SetLabelColor(tcell.ColorRed)
			return
		}

		backend.InsertCategory(backend.Category{Name: iName, Desc: iDesc, IsIncome: iIncome})
		updateCategoriesTable()
		closeForm()
	}

	catDetailsForm = tview.NewForm().
		AddInputField("Name", "", 20, nil, func(v string) { iName = v }).
		AddInputField("Description", "", 40, nil, func(v string) { iDesc = v }).
		AddCheckbox("Is Income?", false, func(b bool) { iIncome = b }).
		AddButton("Add", onSubmit).
		AddButton("Cancel", closeForm)

	catDetailsForm.SetTitle("New Category").SetBorder(true)

	catDetailsForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			closeForm()
			return nil
		}
		return event
	})
}

func showNewCategoryForm() {
	flex.AddItem(catDetailsForm, 55, 0, true)
	app.SetFocus(catDetailsForm)
}
