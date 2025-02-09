package frontend

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tableView struct {
	table *MyTable
	// fetches data from backend, then displays that in the table
	fUpdate func()
	title   string
	headers []string
	curPage int8
	maxPage int8
}

func (tv tableView) hide(parent *tview.Flex) {
	parent.RemoveItem(tv.table)
	app.SetFocus(parent)
}

func isPartialDate(s string, r rune) bool {
	regex0 := regexp.MustCompile(`^\d{0,4}$`)
	regex1 := regexp.MustCompile(`^\d{4}-\d{0,2}$`)
	regex2 := regexp.MustCompile(`^\d{4}-\d{2}-\d{0,2}$`)
	return regex0.MatchString(s) || regex1.MatchString(s) || regex2.MatchString(s)
}

/*
Changes the page that a table shows
Inputs:
  - table: the tview table
  - pgRef (*int8): the variable holding the current page number
  - newPage (int8): the page number to try to switch to
  - maxPage (int8): the last page index that can be accessed
  - title (string): the title to show in the table border
*/
func changePage(tv *tableView, newPage int8) {
	if newPage < 0 || newPage > tv.maxPage {
		return
	}
	tv.curPage = newPage
	title := fmt.Sprintf("%s", tv.title)
	if tv.maxPage > 0 {
		title += fmt.Sprintf(" (%d/%d)", newPage+1, tv.maxPage+1)
	}
	tv.table.SetTitle(title)
	tv.fUpdate()
}

func createTableHeaders(tv *tableView) {
	tv.table.Clear()
	for i, h := range tv.headers {
		tv.table.SetCell(
			0, i,
			tview.NewTableCell(h).
				SetSelectable(false).
				SetStyle(tcell.StyleDefault.Bold(true)),
		)
	}
}

func showTable(parent *tview.Flex, tv *tableView) {
	changePage(tv, 0)
	parent.AddItem(tv.table, 0, 1, true)
	app.SetFocus(tv.table)
}

func stringToInt(s string) int {
	res, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	if err != nil {
		panic(err)
	}
	return int(res)
}

/*
Was the key pressed one that should cause a 'back' navigation?
For all views except forms
*/
func isBackKey(event *tcell.EventKey) bool {
	return event.Rune() == 'q' ||
		event.Rune() == 'h' ||
		event.Key() == tcell.KeyCtrlC ||
		event.Key() == tcell.KeyCtrlQ ||
		event.Key() == tcell.KeyEscape
}

// Custom Table

type MyTable struct {
	*tview.Table
}

/* Creates a table with my settings */
func createMyTable() *MyTable {
	t := tview.NewTable().
		SetBorders(false).
		SetSeparator(tview.Borders.Vertical).
		SetFixed(1, 0).
		SetSelectable(true, false)
	t.SetBorder(true).SetBorderPadding(1, 1, 2, 2)
	return &MyTable{t}
}

func (t *MyTable) getCellInt(row, col int) int {
	res, _ := strconv.ParseInt(strings.TrimSpace(t.GetCell(row, col).Text), 10, 32)
	return int(res)
}

/* Gets the value of the cell (row,col), trimming whitespace and '$' symbols */
func (t *MyTable) getCellString(row, col int) string {
	return strings.Trim(t.GetCell(row, col).Text, " $")
}

// FORM

func formInputCapture(onCancel, onSubmit func()) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyEscape {
			onCancel()
			return nil
		} else if event.Key() == tcell.KeyEnter && event.Modifiers() == tcell.ModCtrl {
			onSubmit()
		}
		return event
	}
}
