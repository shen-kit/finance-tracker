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
	categoriesTable *tview.Table
)

func CreateTUI() {
	setTheme()
	app = tview.NewApplication()
	pages = tview.NewPages()

	createHomepage()
	createCategoriesPage()

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
		AddItem("Add Record", "", 'a', func() { pages.SwitchToPage("add_record") }).
		AddItem("View Month Summary", "", 'm', nil).
		AddItem("View Year Summary", "", 'y', nil).
		AddItem("Records", "", 'r', nil).
		AddItem("Categories", "", 'c', func() { pages.SwitchToPage("categories"); updateCategoriesTable() }).
		AddItem("Investments", "", 'i', func() { pages.SwitchToPage("investments") }).
		AddItem("Quit", "", 'q', func() { app.Stop() })

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

	pages.AddPage("homepage", lv, true, true)
}

func createCategoriesPage() {
	categoriesTable = tview.NewTable().
		SetBorders(false).
		SetSeparator(tview.Borders.Vertical).
		SetFixed(1, 0).
		SetSelectable(true, false)

	categoriesTable.SetBorder(true).SetTitle("Categories")

	categoriesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'h' {
			pages.SwitchToPage("homepage")
			return nil
		}
		return event
	})

	updateCategoriesTable()
	pages.AddPage("categories", categoriesTable, true, false)
}

/*
Update the categories table by pulling data from the backend
*/
func updateCategoriesTable() {
	cats, err := backend.GetCategories()
	if err != nil {
		panic(err)
	}

	headers := strings.Split(" ID : Name : Type : Description ", ":")
	for i, h := range headers {
		categoriesTable.SetCell(0, i, tview.NewTableCell(h).SetSelectable(false))
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
