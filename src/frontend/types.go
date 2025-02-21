package frontend

import (
	"fmt"
	"math"
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

type borderColorChanger interface {
	SetBorderColor(tcell.Color) *tview.Box
}

type updatablePrim interface {
	updatable
	tview.Primitive
	borderColorChanger
}

func showUpdatablePrim(p updatablePrim) {
	p.reset()
	flex.AddItem(p, 0, 1, false)
}

func focusUpdatablePrim(p updatablePrim) {
	app.SetFocus(p)
}

type updatableTable struct {
	*tview.Table
	title       string
	headers     []string
	curPage     int `default:"0"`
	maxPage     int `default:"0"`
	fGetMaxPage func() int
}

func (t *updatableTable) fGetData(int) []backend.DataRow {
	switch t.title {
	case "Records":
		return backend.GetRecordsRecent(t.curPage)
	case "Categories":
		return backend.GetCategories(t.curPage)
	case "Investments":
		return backend.GetInvestmentsRecent(t.curPage)
	case "Investment Summary":
		return backend.GetInvestmentSummary()
	default:
		panic("fGetData encountered an unknown title: " + t.title)
	}
}

func (t *updatableTable) getCurPage() int { return t.curPage }

/* handles keys common to all tables (back, prev/next page) */
func (t *updatableTable) defaultInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if isBackKey(event) {
		app.SetFocus(flex)
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
	t.maxPage = t.fGetMaxPage()
	t.createHeaders()

	// create table body
	for i, row := range rows {
		for j, str := range row.SpreadToStrings() {
			newCell := tview.NewTableCell(" " + str + " ")
			if strings.Contains(str, "#") { // use '#' symbol anywhere to leave as default colour
				newCell.SetText(strings.Replace(newCell.Text, "#", "", 1))
			} else if strings.Contains(str, "$") { // set text colour red/green gradient for money cells
				r, g, b := tview.Styles.PrimaryTextColor.RGB()
				if f, err := strconv.ParseFloat(strings.Trim(str, " $"), 32); err == nil {
					fInt := int32(f)
					rNew := max(20, min(r-fInt, 240))
					gNew := max(20, min(g+fInt, 240))
					bNew := max(20, min(b-int32(math.Abs(float64(fInt))), 240))
					newCell.SetTextColor(tcell.NewRGBColor(rNew, gNew, bNew))
					newCell.SetText(strings.Replace(newCell.Text, "$-", "-$", 1))
				}
			}
			t.SetCell(i+1, j, newCell)
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

// get the text of a cell, remove padding spaces + '$' symbol
func (t *updatableTable) getCellString(row, col int) string {
	return strings.Replace(strings.TrimSpace(t.GetCell(row, col).Text), "$", "", 1)
}

func newUpdatableTable(headers []string, parent borderColorChanger) updatableTable {
	t := tview.NewTable().
		SetBorders(false).
		SetSeparator(tview.Borders.Vertical).
		SetFixed(1, 0)
	t.SetBorder(true).SetBorderPadding(1, 1, 2, 2)

	t.SetFocusFunc(func() {
		t.SetSelectable(true, false)
		t.SetBorderColor(tview.Styles.TertiaryTextColor)
		if parent != nil {
			parent.SetBorderColor(tview.Styles.TertiaryTextColor)
		}
	}).SetBlurFunc(func() {
		t.SetSelectable(false, false)
		t.SetBorderColor(tview.Styles.BorderColor)
		if parent != nil {
			parent.SetBorderColor(tview.Styles.BorderColor)
		}

	})

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

type yearView struct {
	*tview.Grid
	invTable   *updatableTable
	recTable   *updatableTable
	yearOffset int `default:"0"`
	tvTitle    *tview.TextView
}

func (yv *yearView) changeYear(by int) {
	yv.yearOffset += by
	yv.update(yv.fGetData(yv.yearOffset))
}

func (yv *yearView) fGetData(offset int) []backend.DataRow {
	year := time.Now().Year() + offset
	return backend.GetYearSummary(year)
}

func (yv *yearView) getCurPage() int { return yv.yearOffset }
