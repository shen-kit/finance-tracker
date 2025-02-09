package frontend

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

// global function to update the categories table using the tableView object
var updateCategoriesTable func()

type categoryForm struct {
	form      *tview.Form
	iName     *tview.InputField
	iDesc     *tview.InputField
	iIsIncome *tview.Checkbox
	tvMsg     *tview.TextView
}

func createCategoriesTable() *tableView {

	/* ===== Helper Functions ===== */

	// returns a closure with the tableView saved
	createUpdateTableClosure := func(tv *tableView) func() {
		return func() {
			createTableHeaders(tv)
			cats := backend.GetCategories()

			for i, cat := range cats {
				id, name, isincome, desc := cat.Spread()
				tv.table.
					SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf(" %d ", id)).SetAlign(tview.AlignCenter)).
					SetCell(i+1, 1, tview.NewTableCell(" "+name+" ")).
					SetCell(i+1, 2, tview.NewTableCell(func() string {
						if isincome {
							return " Income "
						}
						return " Expenditure "
					}())).
					SetCell(i+1, 3, tview.NewTableCell(" "+desc+" "))
			}
		}
	}

	/* ===== Function Body ===== */
	table := createMyTable()
	tv := &tableView{
		table:   table,
		title:   "Categories",
		headers: strings.Split(" ID : Name : Type : Description ", ":"),
	}
	updateCategoriesTable = createUpdateTableClosure(tv)
	tv.fUpdate = updateCategoriesTable
	return tv
}

func setCategoryTableKeybinds(tv *tableView, cf categoryForm) {
	table := tv.table
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if isBackKey(event) {
			tv.hide(flex)
			return nil
		} else if event.Rune() == 'a' {
			showCategoryForm(table, cf, -1, "", "", false)
			return nil
		} else if event.Rune() == 'd' { // delete category
			row, _ := table.GetSelection()
			id := table.getCellInt(row, 0)
			backend.DeleteCategory(id)
			tv.fUpdate()
			return nil
		} else if event.Rune() == 'e' { // edit category
			row, _ := table.GetSelection()
			id := table.getCellInt(row, 0)
			name := table.getCellString(row, 1)
			isIncome := strings.EqualFold("income", table.getCellString(row, 2))
			desc := table.getCellString(row, 3)
			showCategoryForm(table, cf, id, name, desc, isIncome)
			return nil
		}
		return event
	})
}

func createCategoryForm() categoryForm {
	var form *tview.Form
	var inName, inDesc *tview.InputField
	var inIsIncome *tview.Checkbox
	var formMsg *tview.TextView

	inName = tview.NewInputField().
		SetLabel("Name").
		SetFieldWidth(20)

	inDesc = tview.NewInputField().
		SetLabel("Description").
		SetFieldWidth(40)

	inIsIncome = tview.NewCheckbox().
		SetLabel("Is Income?")

	formMsg = tview.NewTextView().
		SetSize(1, 35).
		SetDynamicColors(true).
		SetScrollable(false)

	form = tview.NewForm().
		AddFormItem(inName).
		AddFormItem(inDesc).
		AddFormItem(inIsIncome).
		AddFormItem(formMsg).
		AddButton("Save", nil).
		AddButton("Cancel", nil)

	form.SetBorder(true)

	return categoryForm{
		form: form, iName: inName, iDesc: inDesc, iIsIncome: inIsIncome, tvMsg: formMsg,
	}
}

func showCategoryForm(lastWidget tview.Primitive, cf categoryForm, id int, name, desc string, isIncome bool) {

	/* ===== Helper Functions ===== */

	setInputFieldValues := func() {
		cf.iName.SetText(name)
		cf.iDesc.SetText(desc)
		cf.iIsIncome.SetChecked(isIncome)
		cf.tvMsg.SetText("")
	}

	closeForm := func() {
		flex.RemoveItem(cf.form)
		app.SetFocus(lastWidget)
	}

	onSubmit := func() {
		cat, err := parseCatForm(cf)
		if err != nil {
			cf.tvMsg.SetText("[red]" + err.Error())
			return
		}

		if id == -1 {
			backend.InsertCategory(cat)
		} else {
			backend.UpdateCategory(id, cat)
		}
		updateCategoriesTable()
		closeForm()
	}

	/* ===== Function Body ===== */

	// set title
	if id == -1 {
		cf.form.SetTitle("Add Category")
	} else {
		cf.form.SetTitle("Edit Category Details")
	}

	setInputFieldValues()

	cf.form.SetInputCapture(formInputCapture(closeForm, onSubmit))
	cf.form.GetButton(cf.form.GetButtonIndex("Cancel")).SetSelectedFunc(closeForm)
	cf.form.GetButton(cf.form.GetButtonIndex("Save")).SetSelectedFunc(onSubmit)

	// display + focus form
	flex.AddItem(cf.form, 55, 0, true)
	cf.form.SetFocus(0)
	app.SetFocus(cf.form)
}

/* Takes input from the form and returns a Category object */
func parseCatForm(cf categoryForm) (backend.Category, error) {

	if cf.iName.GetText() == "" || cf.iDesc.GetText() == "" {
		return backend.Category{}, errors.New("All fields are required")
	}

	return backend.Category{
			Name:     cf.iName.GetText(),
			Desc:     cf.iDesc.GetText(),
			IsIncome: cf.iIsIncome.IsChecked()},
		nil
}
