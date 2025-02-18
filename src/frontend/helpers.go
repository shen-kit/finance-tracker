package frontend

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// UI

// remove all items except the first one
func clearScreen() {
	for range flex.GetItemCount() - 1 {
		flex.RemoveItem(flex.GetItem(1))
	}
	app.SetFocus(flex)
}

// Utility Functions

func isPartialDate(s string, _ rune) bool {
	regex0 := regexp.MustCompile(`^\d{0,4}$`)
	regex1 := regexp.MustCompile(`^\d{4}-\d{0,2}$`)
	regex2 := regexp.MustCompile(`^\d{4}-\d{2}-\d{0,2}$`)
	return regex0.MatchString(s) || regex1.MatchString(s) || regex2.MatchString(s)
}

func stringToInt(s string) int {
	res, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	if err != nil {
		panic(err)
	}
	return int(res)
}

/*
Was the key pressed one that should cause a 'back' navigation?
For all views except forms
*/
func isBackKey(event *tcell.EventKey) bool {
	return event.Rune() == 'q' ||
		event.Rune() == 'h' ||
		event.Key() == tcell.KeyCtrlC ||
		event.Key() == tcell.KeyCtrlQ ||
		event.Key() == tcell.KeyEscape
}

// Form Input

func formInputCapture(onCancel, onSubmit func()) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyCtrlQ || event.Key() == tcell.KeyEscape {
			onCancel()
			return nil
		} else if event.Key() == tcell.KeyEnter && event.Modifiers() == tcell.ModCtrl {
			onSubmit()
			return nil
		}
		return event
	}
}
