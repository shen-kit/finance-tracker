package frontend

import (
	"errors"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

type categoryForm struct {
	form      *tview.Form
	iName     *tview.InputField
	iDesc     *tview.InputField
	iIsIncome *tview.Checkbox
	tvMsg     *tview.TextView
}

func createCategoriesView() *updatableTable {
	table := newUpdatableTable(strings.Split("ID:Name:Type:Description", ":"), nil)
	table.title = "Categories"
	table.fGetMaxPage = func() int { return 0 }
	return &table
}

func setCatTableKeybinds(t *updatableTable, cf categoryForm) {
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if res := t.defaultInputCapture(event); res == nil {
			return nil
		}

		if event.Rune() == 'a' {
			showCategoryForm(t, cf, -1, "", "", false)
		} else if event.Rune() == 'd' { // delete category
			row, _ := t.GetSelection()
			id := t.getCellInt(row, 0)
			showModal("Confirm delete? (y/n)", func() {
				backend.DeleteCategory(id)
				t.update(t.fGetData(t.curPage))
				// set focus if deleted last row
				if row > t.GetRowCount()-1 {
					t.Select(max(0, row-1), 0)
				}
			}, t)
		} else if event.Rune() == 'e' { // edit category
			row, _ := t.GetSelection()
			id := t.getCellInt(row, 0)
			name := t.getCellString(row, 1)
			isIncome := strings.EqualFold("income", t.getCellString(row, 2))
			desc := t.getCellString(row, 3)
			showCategoryForm(t, cf, id, name, desc, isIncome)
		} else {
			return event
		}
		return nil
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
		AddButton("Cancel", nil).
		SetFieldBackgroundColor(tview.Styles.MoreContrastBackgroundColor).
		SetButtonBackgroundColor(tview.Styles.MoreContrastBackgroundColor)

	form.SetBorder(true).
		SetBorderColor(tview.Styles.TertiaryTextColor)

	return categoryForm{
		form: form, iName: inName, iDesc: inDesc, iIsIncome: inIsIncome, tvMsg: formMsg,
	}
}

func showCategoryForm(ct *updatableTable, cf categoryForm, id int, name, desc string, isIncome bool) {

	/* ===== Helper Functions ===== */

	setInputFieldValues := func() {
		cf.iName.SetText(name)
		cf.iDesc.SetText(desc)
		cf.iIsIncome.SetChecked(isIncome)
		cf.tvMsg.SetText("")
	}

	closeForm := func() {
		flex.RemoveItem(cf.form)
		app.SetFocus(ct)
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
		ct.update(ct.fGetData(ct.curPage))
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
