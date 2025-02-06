package frontend

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

var (
	app             *tview.Application
	pages           *tview.Pages
	flex            *tview.Flex
	categoriesTable *tview.Table
	newCatForm      *tview.Form
)

func CreateTUI() {
	setTheme()
	app = tview.NewApplication()
	pages = tview.NewPages()

	createHomepage()
	createCategoriesTable()
	createNewCategoryForm()

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}

func setTheme() {
	// from https://catppuccin.com/palette

	// mocha
	// tview.Styles = tview.Theme{
	// 	PrimitiveBackgroundColor:    tcell.NewRGBColor(30, 30, 46),    // Main background color for primitives.
	// 	ContrastBackgroundColor:     tcell.NewRGBColor(148, 226, 213), // Background color for contrasting elements.
	// 	MoreContrastBackgroundColor: tcell.NewRGBColor(250, 179, 135), // Background color for even more contrasting elements.
	// 	BorderColor:                 tcell.NewRGBColor(127, 132, 156), // Box borders.
	// 	TitleColor:                  tcell.NewRGBColor(205, 214, 244), // Box titles.
	// 	PrimaryTextColor:            tcell.NewRGBColor(205, 214, 244), // Primary text.
	// 	SecondaryTextColor:          tcell.NewRGBColor(186, 194, 222), // Secondary text (e.g. labels).
	// 	TertiaryTextColor:           tcell.NewRGBColor(166, 173, 200), // Tertiary text (e.g. subtitles, notes).
	// 	InverseTextColor:            tcell.NewRGBColor(30, 30, 46),    // Text on primary-colored backgrounds.
	// 	ContrastSecondaryTextColor:  tcell.NewRGBColor(49, 50, 68),    // Secondary text on ContrastBackgroundColor-colored backgrounds.
	// }

	// frappe
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.NewRGBColor(48, 52, 70),    // Main background color for primitives.
		ContrastBackgroundColor:     tcell.NewRGBColor(129, 200, 190), // Background color for contrasting elements.
		MoreContrastBackgroundColor: tcell.NewRGBColor(239, 159, 118), // Background color for even more contrasting elements.
		BorderColor:                 tcell.NewRGBColor(148, 156, 187), // Box borders.
		TitleColor:                  tcell.NewRGBColor(198, 208, 245), // Box titles.
		PrimaryTextColor:            tcell.NewRGBColor(198, 208, 245), // Primary text.
		SecondaryTextColor:          tcell.NewRGBColor(181, 191, 226), // Secondary text (e.g. labels).
		TertiaryTextColor:           tcell.NewRGBColor(165, 173, 206), // Tertiary text (e.g. subtitles, notes).
		InverseTextColor:            tcell.NewRGBColor(48, 52, 70),    // Text on primary-colored backgrounds.
		ContrastSecondaryTextColor:  tcell.NewRGBColor(65, 69, 89),    // Secondary text on ContrastBackgroundColor-colored backgrounds.
	}
}

func createHomepage() {
	lv := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(tview.Styles.ContrastBackgroundColor).
		AddItem("  Add Record            ", "", 0, nil).
		AddItem("  View Month Summary    ", "", 0, nil).
		AddItem("  View Year Summary     ", "", 0, nil).
		AddItem("  Records               ", "", 0, nil).
		AddItem("  Categories            ", "", 0, func() { showCategoriesTable() }).
		AddItem("  Investments           ", "", 0, nil).
		AddItem("  Quit                  ", "", 0, func() { app.Stop() })

	lv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		} else if event.Rune() == 'k' {
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		} else if event.Rune() == 'l' {
			return tcell.NewEventKey(tcell.KeyEnter, 'l', tcell.ModNone)
		}
		return event
	})

	lv.SetTitle("Options").SetBorder(true).SetBorderPadding(1, 1, 2, 2)

	flex = tview.NewFlex().
		AddItem(lv, 30, 0, true)

	pages.AddPage("homepage", flex, true, true)
}

func createCategoriesTable() {
	categoriesTable = tview.NewTable().
		SetBorders(false).
		SetSeparator(tview.Borders.Vertical).
		SetFixed(1, 0).            // always show header row
		SetSelectable(true, false) // rows selectable

	categoriesTable.SetTitle("Categories").SetBorder(true).SetBorderPadding(1, 1, 2, 2)

	categoriesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'h' {
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
		flex.RemoveItem(newCatForm)
		app.SetFocus(categoriesTable)
	}

	onSubmit := func() {
		backend.InsertCategory(backend.Category{Name: iName, Desc: iDesc, IsIncome: iIncome})
		updateCategoriesTable()
		closeForm()
	}

	newCatForm = tview.NewForm().
		AddInputField("Name", "", 20, nil, func(v string) { iName = v }).
		AddInputField("Description", "", 40, nil, func(v string) { iDesc = v }).
		AddCheckbox("Is Income?", false, func(b bool) { iIncome = b }).
		AddButton("Add", onSubmit).
		AddButton("Cancel", closeForm)

	newCatForm.SetTitle("New Category").SetBorder(true)

	newCatForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'h' {
			closeForm()
			return nil
		}
		return event
	})
}

func showNewCategoryForm() {
	flex.AddItem(newCatForm, 55, 0, true)
	app.SetFocus(newCatForm)
}
