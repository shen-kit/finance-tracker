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
		rightAlign(float32(rec.Amt)/100, 2, 8, "$"),
	}
}

type Category struct {
	Id       int // special: 0 = "Net Change" | -1 = "Deleted" | -2 = blank | -3 = "Total Income" | -4 = "Total Expenditure"
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
		"#" + rightAlign(float32(inv.Unitprice)/100, 2, 8, "$"),   // unitprice
		rightAlign(inv.Qty, 1, 6, ""),                             // qty
		rightAlign(float32(inv.Unitprice)*inv.Qty/100, 2, 9, "$"), // value
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
			rec.CatId = -1
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
	if cy.CatId == -2 { // blank line
		return []string{"-----------------", "------", "------", "------", "------", "------", "------", "------", "------", "------", "------", "------", "------"}
	}
	var res = make([]string, 13, 13)
	res[0] = GetCategoryNameFromId(cy.CatId)
	for i, val := range cy.MonthSums {
		res[i+1] = rightAlign(float32(val)/100, 0, 6, "$")
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
	if isr.code == "separator" {
		return []string{"------", "------", "-------------", "-------------", "----------", "-------------", "---------", "-------"}
	}

	avgBuyF := float32(isr.avgBuy) / 100

	if isr.code == "total" {
		return []string{
			"Total",
			"", "", "", // qty, avg buy, cur price
			"#" + rightAlign(avgBuyF, 2, 9, "$"),                           // total in
			"#" + rightAlign(isr.curPrice, 2, 12, "$"),                     // current value
			rightAlign(isr.curPrice-avgBuyF, 2, 9, "$"),                    // P/L
			rightAlign(100*(isr.curPrice-avgBuyF)/avgBuyF, 2, 6, "") + "%", // %P/L
		}
	}

	totalIn := avgBuyF * isr.qty
	curVal := isr.curPrice * isr.qty
	return []string{
		isr.code,                                                 // code
		rightAlign(isr.qty, 2, 6, ""),                            // qty
		"#" + rightAlign(avgBuyF, 2, 12, "$"),                    // avg buy
		"#" + rightAlign(isr.curPrice, 2, 12, "$"),               // cur price
		"#" + rightAlign(totalIn, 2, 9, "$"),                     // total in
		"#" + rightAlign(curVal, 2, 12, "$"),                     // current val
		rightAlign(curVal-totalIn, 2, 9, "$"),                    // P/L
		rightAlign(100*(curVal-totalIn)/totalIn, 2, 6, "") + "%", // %P/L
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
