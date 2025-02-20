package frontend

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

func createYearView() *yearView {
	tvTitle := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvTitle.SetBorderPadding(1, 1, 3, 3)

	yearRecTable := newUpdatableTable(strings.Split(":Jan:Feb:Mar:Apr:May:Jun:Jul:Aug:Sep:Oct:Nov:Dec", ":"))
	yearRecTable.SetBorder(false)
	yearRecTable.fGetMaxPage = func() int { return 0 }

	yearGrid := tview.NewGrid().
		SetRows(3, 0).
		SetBorders(true).
		AddItem(tvTitle, 0, 0, 1, 1, 0, 0, false).
		AddItem(yearRecTable, 1, 0, 1, 1, 0, 0, true)

	yearGrid.SetBorder(true).SetTitle("Year Overview")

	return &yearView{
		Grid:       yearGrid,
		recTable:   &yearRecTable,
		yearOffset: 0,
		tvTitle:    tvTitle,
	}
}

func setYearViewKeybinds(yv *yearView) {
	yv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if isBackKey(event) {
			app.SetFocus(flex)
		} else if event.Rune() == 'H' {
			yv.changeYear(-1)
		} else if event.Rune() == 'L' {
			yv.changeYear(1)
		} else {
			return event
		}
		return nil
	})
}

func (yv *yearView) update(data []backend.DataRow) {
	// set title text to the year being viewed
	t := time.Now().AddDate(yv.yearOffset, 0, 0)
	yv.tvTitle.SetText(t.Format("2006"))

	// update table data
	yv.recTable.update(data)
}

func (yv *yearView) reset() {
	yv.changeYear(-yv.yearOffset)
	yv.update(yv.fGetData(0))
}
