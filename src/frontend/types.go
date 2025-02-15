package frontend

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shen-kit/finance-tracker/backend"
)

type updatable interface {
	update([]backend.DataRow)
	reset()
}

type updatablePrim interface {
	updatable
	tview.Primitive
}

type updatableTable struct {
	tview.Table
	title    string
	headers  []string
	curPage  int `default:"0"`
	maxPage  int `default:"0"`
	fGetData func() []backend.DataRow
}

func (t updatableTable) update(rows []backend.DataRow) {
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

func (t updatableTable) reset() {
	t.changePage(-t.curPage)
	t.update(t.fGetData())
}

func (t updatableTable) changePage(by int) {
	if t.curPage+by < 0 || t.curPage+by > t.maxPage {
		return
	}
	t.curPage += by
	t.update(t.fGetData())

	title := t.title
	if t.maxPage > 0 {
		title += fmt.Sprintf(" (%d/%d)", t.curPage+1, t.maxPage+1)
	}
	t.SetTitle(title)
}

func (t updatableTable) createHeaders() {
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

func (t updatableTable) getCellInt(row, col int) int {
	res, _ := strconv.ParseInt(strings.TrimSpace(t.GetCell(row, col).Text), 10, 32)
	return int(res)
}

func (t updatableTable) getCellString(row, col int) string {
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
		Table:   *t,
		headers: headers,
	}
}

type monthGridView struct {
	updatable
	*tview.Grid
	curOffset int `default:"0"`
}

func showUpdatablePrim(p updatablePrim) {
	p.reset()
	flex.AddItem(p, 0, 1, true)
	app.SetFocus(p)
}
