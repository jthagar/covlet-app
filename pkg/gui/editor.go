package gui

import (
	"cover-letter-templates/pkg/config"
	"cover-letter-templates/pkg/internal"
	"fmt"
	"text/template"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TextEditor struct {
	editor      *widget.Entry
	currentPath string
	undoText    []string
	redoText    []string
}

// NewEditor returns the editor container and the underlying text entry widget
func NewEditor(w fyne.Window) (*fyne.Container, *TextEditor) {
	// TODO: modify this section into the active template editor
	editor := newEditor()

	// ctrl + S to Save
	ctrlS := desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlS, func(shortcut fyne.Shortcut) {
		err := editor.save(w)
		if err != nil {
			return
		}
	})
	// ctrl + N to create a new File
	ctrlN := desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlN, func(shortcut fyne.Shortcut) {
		editor.new(w)
	})
	// ctrl + F to find
	ctrlF := desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlF, func(shortcut fyne.Shortcut) {
		editor.find(w)
	})
	// ctrl + U to undo
	ctrlU := desktop.CustomShortcut{KeyName: fyne.KeyU, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlU, func(shortcut fyne.Shortcut) {
		editor.undo()
	})
	// ctrl + R to redo
	ctrlR := desktop.CustomShortcut{KeyName: fyne.KeyR, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlR, func(shortcut fyne.Shortcut) {
		editor.redo()
	})
	// toolbar versions of the above
	help := widget.NewToolbarAction(theme.HelpIcon(), func() {
		dialog.ShowInformation("Help", "Help", w)
	})
	newFile := widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
		editor.new(w)
	})
	findItem := widget.NewToolbarAction(theme.SearchIcon(), func() {
		editor.find(w)
	})
	undoChange := widget.NewToolbarAction(theme.ContentUndoIcon(), func() {
		editor.undo()
	})
	redoChange := widget.NewToolbarAction(theme.ContentRedoIcon(), func() {
		editor.redo()
	})

	menu := widget.NewToolbar(help, &widget.ToolbarSeparator{}, newFile, findItem, undoChange, redoChange)

	// create tabs to allow multiple templates to be opened and edited
	tabs := container.NewDocTabs(container.NewTabItem("Editor", editor.editor))
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.CreateTab = func() *container.TabItem { return container.NewTabItem("NewTab", newEditor().editor) }

	// create a new window to host rendering at the bottom of the app
	renderButton := widget.NewButton("Render", func() {
		configFile, err := config.LoadConfig("config.yml")
		if err != nil {
			fmt.Printf("error loading config: %v", err)
			return
		}

		// todo: modify into custom func maybe, and modify to allow input vals on top of config default
		r, err := internal.RenderEditor(editor.ConvertText(), configFile.Resume)
		if err != nil {
			fmt.Printf("error rendering template: %v", err)
			return
		}
		// todo: render to new bottom window (idk)
		text := widget.NewMultiLineEntry()
		text.SetText(string(r))
		rContent := container.NewBorder(nil, nil, nil, nil, text)
		rWindow := fyne.CurrentApp().NewWindow("Rendered Cover Letter")
		rWindow.SetMainMenu(renderMenu(rWindow))
		rWindow.SetContent(rContent)
		rWindow.Resize(fyne.NewSize(1000, 700))
		rWindow.Show()
	})

	// Return a simple container that hosts the template editor. Additional controls
	// (preview, render, save) can be added by the caller using container composites.
	return container.NewBorder(menu, renderButton, nil, nil, tabs), editor
}

// newEditor returns the underlying text entry widget
func newEditor() *TextEditor {

	// TODO: modify this section into the active template editor
	e := &TextEditor{
		editor:      widget.NewMultiLineEntry(),
		currentPath: "",
		undoText:    []string{},
		redoText:    []string{},
	}

	e.editor.SetPlaceHolder("Enter Go template here, e.g., {{ .Name }}")
	///////////////////

	e.editor.Wrapping = fyne.TextWrapWord
	e.editor.OnChanged = func(s string) {
		if len(s) > 0 {
			lastChar := s[len(s)-1:]
			//TODO: improve later
			switch lastChar {
			case "{":
				if len(e.undoText) > 0 {
					if e.undoText[len(e.undoText)-1] != s {
						e.editor.SetText(s + "}")
					}
				} else {
					e.editor.SetText(s + "}")
				}
			case "(":
				if len(e.undoText) > 0 {
					if e.undoText[len(e.undoText)-1] != s {
						e.editor.SetText(s + ")")
					}
				} else {
					e.editor.SetText(s + ")")
				}
			case "[":
				if len(e.undoText) > 0 {
					if e.undoText[len(e.undoText)-1] != s {
						e.editor.SetText(s + "]")
					}
				} else {
					e.editor.SetText(s + "]")
				}
			}

			if len(e.undoText) < 6 {
				e.undoText = append(e.undoText, s)
			} else {
				e.undoText = append(e.undoText[1:], s)
			}
		}
	}

	return e
}

func (e *TextEditor) ConvertText() *template.Template {
	text := e.editor.Text
	return template.Must(template.New("").Parse(text))
}
