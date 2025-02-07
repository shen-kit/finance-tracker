package frontend

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app   *tview.Application
	pages *tview.Pages
	flex  *tview.Flex
)

func CreateTUI() {
	setTheme()
	app = tview.NewApplication()
	pages = tview.NewPages()

	createHomepage()
	createCategoriesTable()
	createNewCategoryForm()
	createInvestmentsTable()
	createNewInvestmentForm()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			app.Stop()
		case tcell.KeyCtrlC: // remap ctrl+c to escape (<C-c> exits app by default)
			return tcell.NewEventKey(tcell.KeyEscape, ' ', tcell.ModNone)
		}
		return event
	})

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}

func setTheme() {
	// from https://catppuccin.com/palette

	// // mocha
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
		AddItem("  Categories            ", "", 0, showCategoriesTable).
		AddItem("  Investments           ", "", 0, showInvestmentsTable).
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
