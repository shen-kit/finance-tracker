package backend

import (
	"database/sql"
	"fmt"
	"time"
)

type DataRow interface {
	// spread to a slice of strings, used to display as a table row
	SpreadToStrings() []string
}

type Record struct {
	Id    int
	Date  time.Time
	CatId int
	Desc  string
	Amt   int
}

func (rec Record) Spread() (int, time.Time, string, int, int) {
	return rec.Id, rec.Date, rec.Desc, rec.Amt, rec.CatId
}

func (rec Record) SpreadToStrings() []string {
	return []string{
		fmt.Sprint(rec.Id),
		rec.Date.Format("2006-01-02"),
		GetCategoryNameFromId(rec.CatId),
		rec.Desc,
		getMoneyCellString(float32(rec.Amt)/100, 2, 8),
	}
}

type Category struct {
	Id       int
	Name     string
	IsIncome bool
	Desc     string
}

func (cat Category) Spread() (int, string, bool, string) {
	return cat.Id, cat.Name, cat.IsIncome, cat.Desc
}

func (c Category) SpreadToStrings() []string {
	if c.IsIncome {
		return []string{fmt.Sprint(c.Id), c.Name, "Income", c.Desc}
	} else {
		return []string{fmt.Sprint(c.Id), c.Name, "Expenditure", c.Desc}
	}
}

type Investment struct {
	Id        int
	Date      time.Time
	Code      string
	Unitprice int
	Qty       float32
}

func (inv Investment) Spread() (id int, date time.Time, code string, unitprice int, qty float32) {
	return inv.Id, inv.Date, inv.Code, inv.Unitprice, inv.Qty
}

// returns in order: ID, date, code, unitprice, qty, total value
func (inv Investment) SpreadToStrings() []string {
	return []string{
		fmt.Sprint(inv.Id),            // id
		inv.Date.Format("2006-01-02"), // date
		inv.Code,                      // code
		getMoneyCellString(float32(inv.Unitprice)/100, 2, 8),         // unitprice
		fmt.Sprintf("%.1f", inv.Qty),                                 // qty
		getMoneyCellString(float32(inv.Unitprice)*inv.Qty/100, 2, 9), // value
	}
}

func dbRowsToInvestments(rows *sql.Rows) []DataRow {
	var investments []DataRow

	// for each row, assign column data to struct fields and append struct to slice
	for rows.Next() {
		var inv Investment
		if err := rows.Scan(&inv.Id, &inv.Date, &inv.Code, &inv.Unitprice, &inv.Qty); err != nil {
			panic(err)
		}
		investments = append(investments, inv)
	}

	// check for errors then return
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return investments
}

func dbRowsToRecords(rows *sql.Rows) []DataRow {
	var records []DataRow

	// for each row, assign column data to struct fields and append struct to slice
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.Id, &rec.Date, &rec.Desc, &rec.Amt, &rec.CatId); err != nil {
			panic(err)
		}
		records = append(records, rec)
	}

	// check for errors then return
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return records
}

func dbRowsToCategories(rows *sql.Rows) []DataRow {
	var categories []DataRow

	// for each row, assign column data to struct fields and append struct to slice
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.Id, &cat.Name, &cat.Desc, &cat.IsIncome); err != nil {
			panic(err)
		}
		categories = append(categories, cat)
	}

	// check for errors then return
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return categories
}

type CategoryYear struct {
	CatId     int
	MonthSums [12]int // sum of records for this category for each month
}

func (cy CategoryYear) SpreadToStrings() []string {
	var res = make([]string, 13, 13)
	res[0] = GetCategoryNameFromId(cy.CatId)
	for i, val := range cy.MonthSums {
		res[i+1] = getMoneyCellString(float32(val)/100, 0, 6)
	}
	return res
}

type InvSummaryRow struct {
	code     string
	qty      float32
	avgBuy   int
	curPrice float32 // float32 as retrieved from yahoo finance
}

func (isr InvSummaryRow) SpreadToStrings() []string {
	avgBuyF := float32(isr.avgBuy) / 100
	totalIn := avgBuyF * isr.qty
	curVal := isr.curPrice * isr.qty
	return []string{
		isr.code,                                            // code
		fmt.Sprintf("%.2f", isr.qty),                        // qty
		fmt.Sprintf("$%.2f", avgBuyF),                       // avg buy
		fmt.Sprintf("$%.2f", isr.curPrice),                  // cur price
		fmt.Sprintf("$%.2f", totalIn),                       // total in
		fmt.Sprintf("$%.2f", curVal),                        // current value
		fmt.Sprintf("$%.2f", curVal-totalIn),                // P/L
		fmt.Sprintf("%.2f%%", 100*(curVal-totalIn)/totalIn), // %P/L
	}
}

type FilterOpts struct {
	minCost   float32
	maxCost   float32
	startDate time.Time
	endDate   time.Time
	catIds    []int
	code      string
}

func NewFilterOpts() FilterOpts {
	/*
	  Set default options for filters, allow functions to be passed to modify these
	*/
	startDate, _ := makeDate(2000, 1, 1)
	endDate, _ := makeDate(3000, 1, 1)

	opts := &FilterOpts{
		minCost:   -10000000,
		maxCost:   10000000,
		startDate: startDate,
		endDate:   endDate,
		catIds:    []int{},
		code:      "",
	}

	return *opts
}

func (opts FilterOpts) WithMinCost(val float32) FilterOpts {
	opts.minCost = val
	return opts
}

func (opts FilterOpts) WithMaxCost(val float32) FilterOpts {
	opts.maxCost = val
	return opts
}

func (opts FilterOpts) WithStartDate(val time.Time) FilterOpts {
	opts.startDate = val
	return opts
}

func (opts FilterOpts) WithEndDate(val time.Time) FilterOpts {
	opts.endDate = val
	return opts
}

func (opts FilterOpts) WithCatId(val []int) FilterOpts {
	opts.catIds = val
	return opts
}

func (opts FilterOpts) WithCode(val string) FilterOpts {
	opts.code = val
	return opts
}
