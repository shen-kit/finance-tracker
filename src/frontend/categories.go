package frontend

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

var (
	categoriesTable *tview.Table

	catDetailsForm   *tview.Form
	catInName        *tview.InputField
	catInDescription *tview.InputField
	catInIsIncome    *tview.Checkbox
	catFormMsg       *tview.TextView

	catEditingId int
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
			catEditingId = -1
			showNewCategoryForm("", "", false)
			return nil
		} else if event.Rune() == 'd' { // delete category
			row, _ := categoriesTable.GetSelection()
			res, err := strconv.ParseInt(strings.Trim(categoriesTable.GetCell(row, 0).Text, " "), 10, 32)
			if err != nil {
				panic(err)
			}
			backend.DeleteCategory(int(res))
			updateCategoriesTable()
		} else if event.Rune() == 'e' { // edit category
			row, _ := categoriesTable.GetSelection()
			id, _ := strconv.ParseInt(strings.Trim(categoriesTable.GetCell(row, 0).Text, " "), 10, 32)
			name := strings.Trim(categoriesTable.GetCell(row, 1).Text, " ")
			desc := strings.Trim(categoriesTable.GetCell(row, 2).Text, " ")
			isIncome := strings.ToUpper(categoriesTable.GetCell(row, 3).Text) == "INCOME"

			catEditingId = int(id)
			showNewCategoryForm(name, desc, isIncome)
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

	closeForm := func() {
		flex.RemoveItem(catDetailsForm)
		app.SetFocus(categoriesTable)
	}

	parseForm := func() (backend.Category, error) {

		if catInName.GetText() == "" || catInDescription.GetText() == "" {
			return backend.Category{}, errors.New("All fields are required")
		}

		return backend.Category{
				Name:     catInName.GetText(),
				Desc:     catInDescription.GetText(),
				IsIncome: catInIsIncome.IsChecked()},
			nil
	}

	onSubmit := func() {

		cat, err := parseForm()
		if err != nil {
			catFormMsg.SetText("[red]" + err.Error())
			return
		}

		if catEditingId == -1 {
			backend.InsertCategory(cat)
		} else {
			backend.UpdateCategory(catEditingId, cat)
		}
		updateCategoriesTable()
		closeForm()
	}

	catInName = tview.NewInputField().
		SetLabel("Name").
		SetFieldWidth(20)

	catInDescription = tview.NewInputField().
		SetLabel("Description").
		SetFieldWidth(40)

	catInIsIncome = tview.NewCheckbox().
		SetLabel("Is Income?")

	catFormMsg = tview.NewTextView().
		SetSize(1, 35).
		SetDynamicColors(true).
		SetScrollable(false)

	catDetailsForm = tview.NewForm().
		AddFormItem(catInName).
		AddFormItem(catInDescription).
		AddFormItem(catInIsIncome).
		AddFormItem(catFormMsg).
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

func showNewCategoryForm(name, desc string, isincome bool) {
	catInName.SetText(name)
	catInDescription.SetText(desc)
	catInIsIncome.SetChecked(isincome)
	catFormMsg.SetText("")

	flex.AddItem(catDetailsForm, 55, 0, true)
	app.SetFocus(catDetailsForm)
}
