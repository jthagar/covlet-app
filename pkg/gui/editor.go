package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Editor returns the editor container and the underlying text entry widget
func Editor(w fyne.Window) (*fyne.Container, *widget.Entry) {
	// TODO: modify this section into the active template editor
	editor := widget.NewMultiLineEntry()
	editor.SetPlaceHolder("Enter Go template here, e.g., {{ .Name }}")
	///////////////////

	editor.Wrapping = fyne.TextWrapWord
	editor.OnChanged = func(s string) {
		if len(s) > 0 {
			lastChar := s[len(s)-1:]
			//TODO: improve later
			switch lastChar {
			case "{":
				if len(undoText) > 0 {
					if undoText[len(undoText)-1] != s {
						editor.SetText(s + "}")
					}
				} else {
					editor.SetText(s + "}")
				}
			case "(":
				if len(undoText) > 0 {
					if undoText[len(undoText)-1] != s {
						editor.SetText(s + ")")
					}
				} else {
					editor.SetText(s + ")")
				}
			case "[":
				if len(undoText) > 0 {
					if undoText[len(undoText)-1] != s {
						editor.SetText(s + "]")
					}
				} else {
					editor.SetText(s + "]")
				}
			}

			if len(undoText) < 6 {
				undoText = append(undoText, s)
			} else {
				undoText = append(undoText[1:], s)
			}
		}
	}

	// ctrl + S to Save
	ctrlS := desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlS, func(shortcut fyne.Shortcut) {
		err := save(w, editor)
		if err != nil {
			return
		}
	})
	// ctrl + N to create a new File
	ctrlN := desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlN, func(shortcut fyne.Shortcut) {
		new(editor, w)
	})
	// ctrl + F to find
	ctrlF := desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlF, func(shortcut fyne.Shortcut) {
		find(editor, w)
	})
	// ctrl + U to undo
	ctrlU := desktop.CustomShortcut{KeyName: fyne.KeyU, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlU, func(shortcut fyne.Shortcut) {
		undo(editor)
	})
	// ctrl + R to redo
	ctrlR := desktop.CustomShortcut{KeyName: fyne.KeyR, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlR, func(shortcut fyne.Shortcut) {
		redo(editor)
	})
	// toolbar versions of the above
	help := widget.NewToolbarAction(theme.HelpIcon(), func() {
		dialog.ShowInformation("Help", "Help", w)
	})
	newFile := widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
		new(editor, w)
	})
	findItem := widget.NewToolbarAction(theme.SearchIcon(), func() {
		find(editor, w)
	})
	undoChange := widget.NewToolbarAction(theme.ContentUndoIcon(), func() {
		undo(editor)
	})
	redoChange := widget.NewToolbarAction(theme.ContentRedoIcon(), func() {
		redo(editor)
	})

	menu := widget.NewToolbar(help, &widget.ToolbarSeparator{}, newFile, findItem, undoChange, redoChange)
	// Return a simple container that hosts the template editor. Additional controls
	// (preview, render, save) can be added by the caller using container composites.
	return container.NewBorder(menu, nil, nil, nil, editor), editor
}
