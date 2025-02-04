package backend

import (
	"time"
)

type Record struct {
	id    int
	date  time.Time
	desc  string
	amt   float32
	catId int
}

func (rec Record) spread() (int, time.Time, string, float32, int) {
	return rec.id, rec.date, rec.desc, rec.amt, rec.catId
}

type Category struct {
	id       int
	name     string
	isIncome bool
	desc     string
}

func (cat Category) spread() (int, string, bool, string) {
	return cat.id, cat.name, cat.isIncome, cat.desc
}

type Investment struct {
	id        int
	date      time.Time
	code      string
	qty       float32
	unitprice float32
}

func (inv Investment) spread() (int, time.Time, string, float32, float32) {
	return inv.id, inv.date, inv.code, inv.qty, inv.unitprice
}

type FilterOpts struct {
	minCost   float32
	maxCost   float32
	startDate time.Time
	endDate   time.Time
	catId     []int
	code      string
}

func NewFilterOpts() FilterOpts {
	/*
	  Set default options for filters, allow functions to be passed to modify these
	*/
	startDate, _ := time.Parse("2006-01-02", "2000-01-01")
	endDate, _ := time.Parse("2006-01-02", "3000-01-01")

	opts := &FilterOpts{
		minCost:   -10000000,
		maxCost:   10000000,
		startDate: startDate,
		endDate:   endDate,
		catId:     []int{},
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
	opts.catId = val
	return opts
}

func (opts FilterOpts) WithCode(val string) FilterOpts {
	opts.code = val
	return opts
}
