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

	invSummary := createInvSummaryTable()
	setInvSummaryTableKeybinds(invSummary)

	createHomepage(recTable, catTable, invTable, invSummary, monthView, yearView)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// ctrl+D to exit, or any typical 'back' key when on option select page
		if event.Key() == tcell.KeyCtrlD ||
			(optionsList.HasFocus() && isBackKey(event)) {
			app.Stop()
			return nil
		} else if event.Key() == tcell.KeyCtrlC { // disable default behaviour (exit app)
			return tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
		} else if flex.GetItemCount() < 3 {
			switch event.Rune() {
			case 'm':
				optionsList.SetCurrentItem(0)
				app.SetFocus(monthView)
			case 'y':
				optionsList.SetCurrentItem(1)
				app.SetFocus(yearView)
			case 'r':
				optionsList.SetCurrentItem(2)
				app.SetFocus(recTable)
			case 'c':
				optionsList.SetCurrentItem(3)
				app.SetFocus(catTable)
			case 'i':
				optionsList.SetCurrentItem(4)
				app.SetFocus(invTable)
			}
		}
		return event
	})

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}

func createHomepage(recTable, catTable, invTable, invSummary *updatableTable, monthView *monthGridView, yearView *yearView) {
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
		AddItem("  Investment Summary ", "invSummary", 0, func() { app.SetFocus(invSummary) }).
		AddItem("  Quit", "quit", 0, func() { app.Stop() })

	optionsList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		clearScreen()
		switch secondaryText {
		case "month":
			showUpdatablePrim(monthView, false)
		case "year":
			showUpdatablePrim(yearView, false)
		case "records":
			showUpdatablePrim(recTable, false)
		case "categories":
			showUpdatablePrim(catTable, false)
		case "investments":
			showUpdatablePrim(invTable, false)
		case "invSummary":
			showUpdatablePrim(invSummary, false)
		}
	})

	optionsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		} else if event.Rune() == 'k' {
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		} else if event.Rune() == 'l' {
			return tcell.NewEventKey(tcell.KeyEnter, 'l', tcell.ModNone)
		} else if event.Rune() == 'g' {
			optionsList.SetCurrentItem(0)
		} else if event.Rune() == 'G' {
			optionsList.SetCurrentItem(optionsList.GetItemCount() - 1)
		}
		return event
	})

	optionsList.SetTitle("Options").
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2)

	flex.AddItem(optionsList, 30, 0, true)
	showUpdatablePrim(monthView, false)

	pages.AddPage("main", flex, true, true)
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
