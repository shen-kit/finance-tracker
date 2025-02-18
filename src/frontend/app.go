package frontend

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app         *tview.Application
	pages       *tview.Pages
	flex        *tview.Flex
	optionsList *tview.List
)

func CreateTUI() {
	setTheme()
	app = tview.NewApplication()
	pages = tview.NewPages()

	recTable := createRecordsTable()
	rf := createRecordForm()
	setRecTableKeybinds(recTable, rf)

	monthView := createMonthSummary(recTable)
	setMonthGridKeybinds(monthView, rf)

	yearView := createYearView()
	setYearViewKeybinds(yearView)

	catTable := createCategoriesView()
	cf := createCategoryForm()
	setCatTableKeybinds(catTable, cf)

	invTable := createInvestmentsTable()
	invForm := createInvestmentForm()
	setInvTableKeybinds(invTable, invForm)

	createHomepage(recTable, catTable, invTable, monthView, yearView)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// ctrl+D to exit, or any typical 'back' key when on option select page
		if event.Key() == tcell.KeyCtrlD ||
			(flex.GetItemCount() == 1 && isBackKey(event)) {
			app.Stop()
			return nil
		} else if event.Key() == tcell.KeyCtrlC { // disable default behaviour (exit app)
			return tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
		} else if event.Rune() == 'm' {
			clearScreen()
			optionsList.SetCurrentItem(0)
			showUpdatablePrim(monthView)
		} else if event.Rune() == 'y' {
			clearScreen()
			optionsList.SetCurrentItem(1)
			showUpdatablePrim(yearView)
		} else if event.Rune() == 'r' {
			clearScreen()
			optionsList.SetCurrentItem(2)
			showUpdatablePrim(recTable)
		} else if event.Rune() == 'c' {
			clearScreen()
			optionsList.SetCurrentItem(3)
			showUpdatablePrim(catTable)
		} else if event.Rune() == 'i' {
			clearScreen()
			optionsList.SetCurrentItem(4)
			showUpdatablePrim(invTable)
		}
		return event
	})

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}

func setTheme() {
	// frappe -> https://catppuccin.com/palette
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

func createHomepage(recTable, catTable, invTable *updatableTable, monthView *monthGridView, yearView *yearView) {
	flex = tview.NewFlex()

	optionsList = tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tview.Styles.ContrastBackgroundColor).
		// AddItem("  Add Record            ", "", 0, func() { showRecordForm(flex, rf, -1, "", "", "", "") }).
		AddItem("  View Month Summary", "month", 0, func() { app.SetFocus(recTable) }).
		AddItem("  View Year Summary", "year", 0, func() { app.SetFocus(yearView) }).
		AddItem("  Records", "records", 0, func() { app.SetFocus(recTable) }).
		AddItem("  Categories", "categories", 0, func() { app.SetFocus(catTable) }).
		AddItem("  Investments", "investments", 0, func() { app.SetFocus(invTable) }).
		AddItem("  Quit", "quit", 0, func() { app.Stop() })

	optionsList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		clearScreen()
		switch secondaryText {
		case "month":
			showUpdatablePrim(monthView)
		case "year":
			showUpdatablePrim(yearView)
		case "records":
			showUpdatablePrim(recTable)
		case "categories":
			showUpdatablePrim(catTable)
		case "investments":
			showUpdatablePrim(invTable)
		}
	})

	optionsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		} else if event.Rune() == 'k' {
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		} else if event.Rune() == 'l' {
			return tcell.NewEventKey(tcell.KeyEnter, 'l', tcell.ModNone)
		}
		return event
	})

	optionsList.SetTitle("Options").SetBorder(true).SetBorderPadding(1, 1, 2, 2)

	flex.AddItem(optionsList, 30, 0, true)
	showUpdatablePrim(monthView)

	pages.AddPage("main", flex, true, true)
}
