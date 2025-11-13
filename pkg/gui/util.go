package gui

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Char struct {
	Char string
	Row  int
	Col  int
}

var currentPath string = ""
var undoText []string = []string{}
var redoText []string = []string{}

func undo(entry *widget.Entry) {
	if len(undoText) > 0 {
		if len(redoText) < 6 {
			redoText = append(redoText, entry.Text)
		} else {
			redoText = redoText[1:]
			redoText = append(redoText, entry.Text)
		}
		entry.SetText(undoText[len(undoText)-1])
		undoText = undoText[:len(undoText)-1]
	}
}

func redo(entry *widget.Entry) {
	if len(redoText) > 0 {
		undoText = append(undoText, entry.Text)
		entry.SetText(redoText[len(redoText)-1])
		redoText = redoText[:len(redoText)-1]
	}
}

// untested, and might break
func find(entry *widget.Entry, w fyne.Window) {
	s := &widget.Entry{Text: "Enter text to search for"}
	dialog.ShowForm("Find", "Search", "Cancel",
		[]*widget.FormItem{
			{Text: "Find:", Widget: s},
		}, func(bool) {
			index := strings.Index(entry.Text, s.Text)
			if index != -1 {
				row := strings.Count(entry.Text[:index], "\n")
				foundedLine := strings.Split(entry.Text[:index], "\n")[row]
				col := len(foundedLine)
				w.Canvas().Focus(entry)
				entry.CursorRow = row
				entry.CursorColumn = col
				entry.Refresh()
				dialog.ShowInformation("Found", fmt.Sprintf("Found '%s' at %d (row: %d, col: %d)", s, index, row, col), w)
			} else {
				dialog.ShowInformation("Not found", fmt.Sprintf("Could not find '%s'", s), w)
			}
		}, w)
}

// I don't fully understand this function, other than it would create a new file I guess
func new(entry *widget.Entry, w fyne.Window) {
	currentPath = ""
	redoText = []string{}

	dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		if uc != nil {
			currentPath = uc.URI().Path()
		}
	}, w)
}

func save(w fyne.Window, entry *widget.Entry) error {
	txt := entry.Text
	txtbyte := []byte(txt)

	if len(currentPath) > 0 {
		newerr := ioutil.WriteFile(currentPath, txtbyte, 0644)
		if newerr != nil {
			log.Fatal(newerr)
			return newerr
		}
	} else {
		dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uc != nil {
				newerr := ioutil.WriteFile(uc.URI().Path(), txtbyte, 0644)
				if newerr != nil {
					log.Fatal(newerr)
					return
				}
				currentPath = uc.URI().Path()
			}
		}, w)
	}

	return nil
}
