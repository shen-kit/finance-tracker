package helper

import (
	"fmt"
	"time"
)

/*
Creates a Time object from a given year, month, and day.
*/
func MakeDate(year, month, day int) (time.Time, error) {
	return time.Parse("2006-01-02", fmt.Sprintf("%04d-%02d-%02d", year, month, day))
}
