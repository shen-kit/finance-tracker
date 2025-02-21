package frontend

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

var (
	app         *tview.Application
	pages       *tview.Pages
	flex        *tview.Flex
	optionsList *tview.List
	modalText   *tview.TextView
)

func CreateTUI() {
	setTheme()
	app = tview.NewApplication()
	pages = tview.NewPages()

	rf := createRecordForm()
	cf := createCategoryForm()
	invForm := createInvestmentForm()

	monthView := createMonthSummary()
	setMonthGridKeybinds(monthView, rf)

	recTable := createRecordsTable(monthView)
	setRecTableKeybinds(recTable, rf)

	monthView.table = recTable
	monthView.AddItem(recTable, 2, 0, 1, 1, 0, 0, true)

	yearView := createYearView()
	setYearViewKeybinds(yearView)

	catTable := createCategoriesView()
	setCatTableKeybinds(catTable, cf)

	invTable := createInvestmentsTable()
	setInvTableKeybinds(invTable, invForm)

	invSummary := createInvSummaryTable()
	setInvSummaryTableKeybinds(invSummary)

	createHomepage(recTable, catTable, invTable, invSummary, monthView, yearView)
	createModal()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// ctrl+D to exit, or any typical 'back' key when on option select page
		if event.Key() == tcell.KeyCtrlD ||
			(optionsList.HasFocus() && isBackKey(event)) {
			app.Stop()
			return nil
		} else if event.Key() == tcell.KeyCtrlC { // disable default behaviour (exit app)
			return tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
		} else if !modalText.HasFocus() && flex.GetItemCount() < 3 {
			switch event.Rune() {
			case 'y':
				optionsList.SetCurrentItem(0)
				focusUpdatablePrim(yearView)
			case 'm':
				optionsList.SetCurrentItem(1)
				focusUpdatablePrim(monthView)
			case 'r':
				optionsList.SetCurrentItem(2)
				focusUpdatablePrim(recTable)
			case 'c':
				optionsList.SetCurrentItem(3)
				focusUpdatablePrim(catTable)
			case 'i':
				optionsList.SetCurrentItem(4)
				focusUpdatablePrim(invTable)
			}
		}
		return event
	})

	app.SetAfterDrawFunc(func(screen tcell.Screen) {
		_, h := screen.Size()
		backend.PAGE_ROWS = h - 5
	})

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}

func createModal() {
	modal := func(p tview.Primitive, width, height int) *tview.Flex {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}
	modalText = tview.NewTextView().SetTextAlign(tview.AlignCenter)
	modalText.
		SetBackgroundColor(tview.Styles.MoreContrastBackgroundColor).
		SetBorder(true)
	pages.AddPage("modal", modal(modalText, 50, 3), true, false)
}

func showModal(s string, actionFunc func(), prev tview.Primitive) {
	modalText.SetText(s).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'n' || event.Rune() == 'N' || isBackKey(event) {
			pages.HidePage("modal")
		} else if event.Rune() == 'y' || event.Rune() == 'Y' {
			actionFunc()
		} else {
			return event
		}
		pages.HidePage("modal")
		app.SetFocus(prev)
		return nil
	})

	pages.ShowPage("modal")
	app.SetFocus(modalText)
}

func createHomepage(recTable, catTable, invTable, invSummary *updatableTable, monthView *monthGridView, yearView *yearView) {
	flex = tview.NewFlex()

	optionsList = tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tview.Styles.ContrastBackgroundColor).
		AddItem("  View Year Summary", "year", 0, func() { focusUpdatablePrim(yearView) }).
		AddItem("  View Month Summary", "month", 0, func() { focusUpdatablePrim(monthView) }).
		AddItem("  Records", "records", 0, func() { focusUpdatablePrim(recTable) }).
		AddItem("  Categories", "categories", 0, func() { focusUpdatablePrim(catTable) }).
		AddItem("  Investments", "investments", 0, func() { focusUpdatablePrim(invTable) }).
		AddItem("  Investment Summary ", "invSummary", 0, func() { focusUpdatablePrim(invSummary) }).
		AddItem("  Quit", "quit", 0, func() { app.Stop() })

	optionsList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		clearScreen()
		switch secondaryText {
		case "year":
			showUpdatablePrim(yearView)
		case "month":
			showUpdatablePrim(monthView)
		case "records":
			showUpdatablePrim(recTable)
		case "categories":
			showUpdatablePrim(catTable)
		case "investments":
			showUpdatablePrim(invTable)
		case "invSummary":
			showUpdatablePrim(invSummary)
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
		SetBorderPadding(1, 1, 2, 2).
		SetFocusFunc(func() { optionsList.SetBorderColor(tview.Styles.TertiaryTextColor) }).
		SetBlurFunc(func() { optionsList.SetBorderColor(tview.Styles.BorderColor) })

	flex.AddItem(optionsList, 30, 0, true)
	optionsList.SetCurrentItem(1) // default to month view

	pages.AddPage("main", flex, true, true)
}

func setTheme() {
	// https://catppuccin.com/palette (frappe)
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.NewRGBColor(48, 52, 70),    // Main background color for primitives.
		ContrastBackgroundColor:     tcell.NewRGBColor(129, 200, 190), // Background color for contrasting elements.
		MoreContrastBackgroundColor: tcell.NewRGBColor(115, 121, 148), // Background color for even more contrasting elements.
		BorderColor:                 tcell.NewRGBColor(148, 156, 187), // Box borders.
		TitleColor:                  tcell.NewRGBColor(198, 208, 245), // Box titles.
		PrimaryTextColor:            tcell.NewRGBColor(198, 208, 245), // Primary text.
		SecondaryTextColor:          tcell.NewRGBColor(181, 191, 226), // Secondary text (e.g. labels).
		TertiaryTextColor:           tcell.NewRGBColor(239, 159, 118), // Tertiary text (e.g. subtitles, notes).
		InverseTextColor:            tcell.NewRGBColor(48, 52, 70),    // Text on primary-colored backgrounds.
		ContrastSecondaryTextColor:  tcell.NewRGBColor(65, 69, 89),    // Secondary text on ContrastBackgroundColor-colored backgrounds.
	}
}
