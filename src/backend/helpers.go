package backend

import (
	"fmt"
	"time"
)

/*
Creates a Time object from a given year, month, and day.
*/
func makeDate(year, month, day int) (time.Time, error) {
	return time.Parse("2006-01-02", fmt.Sprintf("%04d-%02d-%02d", year, month, day))
}

/* Returns times one second before the start and end of the month */
func getMonthStartAndEnd(t time.Time) (time.Time, time.Time) {
	year, month, _ := t.Date()
	loc := t.Location()
	mStart := time.Date(year, month, 1, 0, 0, -1, 0, loc)
	mEnd := time.Date(year, month+1, 1, 0, 0, -1, 0, loc)
	// mStart = mStart.Add(-time.Second)
	// mEnd := mStart.AddDate(0, 1, 0)
	// mEnd = mEnd.Add(-time.Second)
	return mStart, mEnd
}

// returns a string right-aligned, with '  $amt.xx' format
func rightAlign(amt float32, decimals, width int, prefix string) string {
	fmtStr1 := fmt.Sprintf("%%%ds", width)
	fmtStr2 := fmt.Sprintf("%s%%.%df", prefix, decimals)
	return fmt.Sprintf(fmtStr1, fmt.Sprintf(fmtStr2, amt))
}
