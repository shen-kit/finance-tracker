package frontend

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

type updatable interface {
	fGetData(int) []backend.DataRow
	getCurPage() int
	update([]backend.DataRow)
	reset()
}

type updatablePrim interface {
	updatable
	tview.Primitive
}

type updatableTable struct {
	*tview.Table
	title   string
	headers []string
	curPage int `default:"0"`
	maxPage int `default:"0"`
	// fGetData func(int) []backend.DataRow
}

func (t *updatableTable) fGetData(int) []backend.DataRow {
	switch t.title {
	case "Records":
		return backend.GetRecordsRecent(t.curPage)
	case "Categories":
		return backend.GetCategories(t.curPage)
	case "Investments":
		return backend.GetInvestmentsRecent(t.curPage)
	default:
		panic("fGetData encountered an unknown title: " + t.title)
	}
}

func (t *updatableTable) getCurPage() int { return t.curPage }

/* handles keys common to all tables (back, prev/next page) */
func (t *updatableTable) defaultInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if isBackKey(event) {
		gotoHomepage()
	} else if event.Rune() == 'L' { // next page
		t.changePage(1)
	} else if event.Rune() == 'H' { // previous page
		t.changePage(-1)
	} else {
		return event
	}
	return nil
}

func (t *updatableTable) update(rows []backend.DataRow) {
	t.createHeaders()

	var strTable [][]string
	for _, row := range rows {
		strTable = append(strTable, row.SpreadToStrings())
	}

	for i, strRow := range strTable {
		for j, strCell := range strRow {
			t.SetCell(i+1, j, tview.NewTableCell(" "+strCell+" "))
		}
	}
}

func (t *updatableTable) reset() {
	t.SetBorder(true)
	t.changePage(-t.curPage)
	t.update(t.fGetData(t.curPage))
}

func (t *updatableTable) changePage(by int) {
	if t.curPage+by < 0 || t.curPage+by > t.maxPage {
		return
	}
	t.curPage += by
	t.update(t.fGetData(t.curPage))

	title := t.title
	if t.maxPage > 0 {
		title += fmt.Sprintf(" (%d/%d)", t.curPage+1, t.maxPage+1)
	}
	t.SetTitle(title)
}

func (t *updatableTable) createHeaders() {
	t.Clear()
	for i, h := range t.headers {
		t.SetCell(
			0, i,
			tview.NewTableCell(" "+h+" ").
				SetSelectable(false).
				SetStyle(tcell.StyleDefault.Bold(true)),
		)
	}
}

func (t *updatableTable) getCellInt(row, col int) int {
	res, _ := strconv.ParseInt(strings.TrimSpace(t.GetCell(row, col).Text), 10, 32)
	return int(res)
}

func (t *updatableTable) getCellString(row, col int) string {
	return strings.Trim(t.GetCell(row, col).Text, " $")
}

func newUpdatableTable(headers []string) updatableTable {
	t := tview.NewTable().
		SetBorders(false).
		SetSeparator(tview.Borders.Vertical).
		SetFixed(1, 0).
		SetSelectable(true, false)
	t.SetBorder(true).SetBorderPadding(1, 1, 2, 2)
	return updatableTable{
		Table:   t,
		headers: headers,
	}
}

type monthGridView struct {
	*tview.Grid
	table       *updatableTable
	monthOffset int `default:"0"`
	tvTitle     *tview.TextView
	tvSummary   *tview.TextView
}

func (mv *monthGridView) fGetData(offset int) []backend.DataRow {
	t := time.Now().AddDate(0, mv.monthOffset, 0)
	records, _, _ := backend.GetMonthInfo(t)
	return records
}

func (mv *monthGridView) getCurPage() int { return mv.monthOffset }

func (mv *monthGridView) changeMonth(by int) {
	mv.monthOffset += by
	mv.update(mv.fGetData(mv.monthOffset))
}

func showUpdatablePrim(p updatablePrim) {
	p.reset()
	flex.AddItem(p, 0, 1, true)
	app.SetFocus(p)
}
