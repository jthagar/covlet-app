package gui

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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
				dialog.ShowInformation("Found", fmt.Sprintf("Found '%s' at %d (row: %d, col: %d)", s.Text, index, row, col), w)
			} else {
				dialog.ShowInformation("Not found", fmt.Sprintf("Could not find '%s'", s.Text), w)
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
		return saveAs(w, entry)
	}

	return nil
}

func saveAs(w fyne.Window, entry *widget.Entry) error {
	txt := entry.Text
	txtbyte := []byte(txt)
	dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if uc != nil {
			if uc.URI() != nil {
				// ensure parent dir exists via fyne storage
				_, _ = uc.Write(txtbyte)
				_ = uc.Close()
				currentPath = uc.URI().Path()
			} else {
				// fallback
				newerr := ioutil.WriteFile(uc.URI().Path(), txtbyte, 0644)
				if newerr != nil {
					log.Fatal(newerr)
					return
				}
				currentPath = uc.URI().Path()
			}
		}
	}, w)
	return nil
}

// wrapAtCursor wraps the current word at the cursor with prefix/suffix.
func wrapAtCursor(entry *widget.Entry, prefix, suffix string) {
	// Calculate absolute index from CursorRow/Column
	text := entry.Text
	row := entry.CursorRow
	col := entry.CursorColumn
	// find index of start of row
	abs := 0
	rows := strings.Split(text, "\n")
	if row > len(rows)-1 {
		row = len(rows) - 1
	}
	for i := 0; i < row; i++ {
		abs += len(rows[i]) + 1 // include newline
	}
	if col > 0 {
		if row < len(rows) {
			if col > len(rows[row]) {
				col = len(rows[row])
			}
			abs += col
		}
	}
	// find word boundaries
	start := abs
	for start > 0 {
		ch := text[start-1]
		if ch == ' ' || ch == '\n' || ch == '\t' {
			break
		}
		start--
	}
	end := abs
	for end < len(text) {
		ch := text[end]
		if ch == ' ' || ch == '\n' || ch == '\t' {
			break
		}
		end++
	}
	var b bytes.Buffer
	b.WriteString(text[:start])
	b.WriteString(prefix)
	b.WriteString(text[start:end])
	b.WriteString(suffix)
	b.WriteString(text[end:])
	entry.SetText(b.String())
	// restore cursor somewhere sensible
	entry.CursorRow = row
	entry.CursorColumn = col + len(prefix)
	entry.Refresh()
}

// custom simple icons for Bold/Italic/Strikethrough using inline SVGs
func BoldIcon() fyne.Resource {
	svg := `<?xml version="1.0" encoding="UTF-8"?>
<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
  <rect x="4" y="3" width="6" height="18" fill="#222"/>
  <circle cx="12" cy="8" r="5" fill="#222"/>
  <circle cx="12" cy="16" r="5" fill="#222"/>
  <rect x="4" y="10" width="8" height="4" fill="#222"/>
</svg>`
	return fyne.NewStaticResource("bold_icon.svg", []byte(svg))
}

func ItalicIcon() fyne.Resource {
	svg := `<?xml version="1.0" encoding="UTF-8"?>
<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
  <polygon points="10,3 20,3 18,7 12,7 8,21 4,21 6,17 12,17" fill="#222"/>
</svg>`
	return fyne.NewStaticResource("italic_icon.svg", []byte(svg))
}

func StrikeIcon() fyne.Resource {
	svg := `<?xml version="1.0" encoding="UTF-8"?>
<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
  <rect x="3" y="11" width="18" height="2" fill="#222"/>
  <circle cx="9" cy="7" r="4" fill="#222"/>
  <circle cx="15" cy="17" r="4" fill="#222"/>
</svg>`
	return fyne.NewStaticResource("strike_icon.svg", []byte(svg))
}

// ensure file path helper for potential future use
func ensureParent(uri fyne.URI) {
	_ = uri
	// Placeholder for future directory checks using storage APIs
	_ = storage.NewURI("file:///")
}
