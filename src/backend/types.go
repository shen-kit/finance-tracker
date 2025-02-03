package backend

import "time"

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
