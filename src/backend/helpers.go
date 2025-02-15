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

func getMonthStartAndEnd(t time.Time) (time.Time, time.Time) {
	year, month, _ := t.Date()
	loc := t.Location()

	mStart := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	mEnd := mStart.AddDate(0, 1, 0)

	return mStart, mEnd
}
