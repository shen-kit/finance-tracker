package frontend

import (
	"fmt"

	"github.com/rivo/tview"
)

const pageCount = 5

func CreateTUI() {
	app := tview.NewApplication()
	pages := tview.NewPages()
	for page := 0; page < pageCount; page++ {
		func(page int) {
			pages.AddPage(fmt.Sprintf("page-%d", page),
				tview.NewModal().
					SetText(fmt.Sprintf("This is page %d. Choose where to go next.", page+1)).
					AddButtons([]string{"Next", "Quit"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonIndex == 0 {
							pages.SwitchToPage(fmt.Sprintf("page-%d", (page+1)%pageCount))
						} else {
							app.Stop()
						}
					}),
				false,
				page == 0)
		}(page)
	}
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}

	// lv := tview.NewList().
	// 	AddItem("Add Record", "", 'a', nil).
	// 	AddItem("View Month Summary", "", 'm', nil).
	// 	AddItem("View Year Summary", "", 'y', nil).
	// 	AddItem("Records", "", 'r', nil).
	// 	AddItem("Categories", "", 'c', nil).
	// 	AddItem("Investments", "", 'i', nil).
	// 	AddItem("Quit", "", 'q', func() { app.Stop() })
	// lv.SetSelectedBackgroundColor(tcell.ColorDarkGreen.TrueColor())
	// lv.SetSelectedTextColor(tcell.ColorWhite.TrueColor())
}
