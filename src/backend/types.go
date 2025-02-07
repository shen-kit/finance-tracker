package backend

import (
	"database/sql"
	"time"

	"github.com/shen-kit/finance-tracker/helper"
)

type Record struct {
	Id    int
	Date  time.Time
	Desc  string
	Amt   float32
	CatId int
}

func (rec Record) spread() (int, time.Time, string, float32, int) {
	return rec.Id, rec.Date, rec.Desc, rec.Amt, rec.CatId
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

type Investment struct {
	Id        int
	Date      time.Time
	Code      string
	Qty       float32
	Unitprice float32
}

func (inv Investment) Spread() (int, time.Time, string, float32, float32) {
	return inv.Id, inv.Date, inv.Code, inv.Qty, inv.Unitprice
}

func dbRowsToInvestments(rows *sql.Rows) ([]Investment, error) {
	var investments []Investment

	// for each row, assign column data to struct fields and append struct to slice
	for rows.Next() {
		var inv Investment
		if err := rows.Scan(&inv.Id, &inv.Date, &inv.Code, &inv.Qty, &inv.Unitprice); err != nil {
			return investments, err
		}
		investments = append(investments, inv)
	}

	// check for errors then return
	if err := rows.Err(); err != nil {
		return investments, err
	}
	return investments, nil
}

func dbRowsToRecords(rows *sql.Rows) ([]Record, error) {
	var records []Record

	// for each row, assign column data to struct fields and append struct to slice
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.Id, &rec.Date, &rec.Desc, &rec.Amt, &rec.CatId); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	// check for errors then return
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func dbRowsToCategories(rows *sql.Rows) ([]Category, error) {
	var categories []Category

	// for each row, assign column data to struct fields and append struct to slice
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.Id, &cat.Name, &cat.Desc, &cat.IsIncome); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	// check for errors then return
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
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
	startDate, _ := helper.MakeDate(2000, 1, 1)
	endDate, _ := helper.MakeDate(3000, 1, 1)

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
