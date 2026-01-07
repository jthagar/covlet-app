package gui

import (
	"covlet/pkg/config"
	"covlet/pkg/internal"
	"fmt"
	"strings"
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
	// sidebar state
	varSidebar *fyne.Container
	varForm    *widget.Form
	// tracked variables and user overrides
	tmplVars  []string
	overrides map[string]string
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

		// Build data by applying overrides on top of config defaults
		data := applyOverrides(configFile.Resume, editor.overrides)

		// render the template
		r, err := internal.RenderEditor(editor.ConvertText(), data)
		if err != nil {
			fmt.Printf("error rendering template: %v", err)
			return
		}
		// todo: render to new bottom window (idk)
		text := widget.NewMultiLineEntry()
		text.SetText(string(r))
        rContent := container.NewBorder(nil, nil, nil, nil, text)
        rWindow := fyne.CurrentApp().NewWindow("Rendered Cover Letter")
        // pass a getter so the menu can export the latest text
        rWindow.SetMainMenu(renderMenu(rWindow, func() string { return text.Text }))
        rWindow.SetContent(rContent)
        rWindow.Resize(fyne.NewSize(1000, 700))
        rWindow.Show()
    })

	// Build variable sidebar on the right
	editor.initVarSidebar()
	mainWindow := container.NewHSplit(tabs, editor.varSidebar)
	mainWindow.Offset = 0.75
	// Compose main editor area with right sidebar
	content := container.NewBorder(menu, renderButton, nil, nil, mainWindow)

	// Return a simple container that hosts the template editor and the var sidebar.
	return content, editor
}

// newEditor returns the underlying text entry widget
func newEditor() *TextEditor {

	// TODO: modify this section into the active template editor
	e := &TextEditor{
		editor:      widget.NewMultiLineEntry(),
		currentPath: "",
		undoText:    []string{},
		redoText:    []string{},
		varSidebar:  nil,
		varForm:     nil,
		tmplVars:    []string{},
		overrides:   map[string]string{},
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
		// Update variable sidebar on any change
		e.refreshVarSidebar()
	}

	return e
}

func (e *TextEditor) ConvertText() *template.Template {
	text := e.editor.Text
	return template.Must(template.New("").Parse(text))
}

// initVarSidebar initializes the sidebar container used to track {{ }} variables
func (e *TextEditor) initVarSidebar() {
	title := widget.NewLabel("Template Variables")
	title.Alignment = fyne.TextAlignLeading
	e.varForm = &widget.Form{}
	info := widget.NewLabel("override any default values")
	clearBtn := widget.NewButton("Clear Overrides", func() {
		e.overrides = map[string]string{}
		e.refreshVarSidebar()
	})
	e.varSidebar = container.NewBorder(title, container.NewVBox(info, clearBtn), nil, nil, container.NewVScroll(e.varForm))
	e.varSidebar.Resize(fyne.NewSize(50, 100))
	// initial populate
	e.refreshVarSidebar()
}

// refreshVarSidebar reparses the text and rebuilds the form entries
func (e *TextEditor) refreshVarSidebar() {
	if e.varForm == nil {
		return
	}
	vars := parseTopLevelVars(e.editor.Text)
	// keep order stable
	e.tmplVars = vars

	// rebuild form
	e.varForm.Items = nil
	for _, v := range vars {
		name := v // capture
		val := e.overrides[name]
		entry := widget.NewEntry()
		entry.SetText(val)
		entry.OnChanged = func(s string) {
			if e.overrides == nil {
				e.overrides = map[string]string{}
			}
			e.overrides[name] = s
		}
		e.varForm.Append(name, entry)
	}
	e.varForm.Refresh()
}

// parseTopLevelVars extracts top-level variable names referenced as {{ .Name }} etc.
// It attempts to collect unique identifiers that immediately follow a dot.
func parseTopLevelVars(s string) []string {
	// very light-weight parse using scanning between {{ and }}
	type void struct{}
	seen := map[string]void{}
	order := []string{}
	i := 0
	for i < len(s) {
		start := strings.Index(s[i:], "{{")
		if start < 0 { // no more
			break
		}
		start += i + 2
		endRel := strings.Index(s[start:], "}}")
		if endRel < 0 {
			break
		}
		end := start + endRel
		expr := s[start:end]
		// look for .Identifier pattern
		// skip spaces
		j := 0
		for j < len(expr) && (expr[j] == ' ' || expr[j] == '\n' || expr[j] == '\t') {
			j++
		}
		for j < len(expr) {
			if expr[j] == '.' {
				// read identifier after dot
				k := j + 1
				// if starts with '(' or other, skip
				if k < len(expr) && isIdentStart(expr[k]) {
					startName := k
					k++
					for k < len(expr) && isIdentPart(expr[k]) {
						k++
					}
					name := expr[startName:k]
					if name != "" {
						if _, ok := seen[name]; !ok {
							seen[name] = void{}
							order = append(order, name)
						}
					}
				}
			}
			j++
		}

		i = end + 2
	}
	return order
}

func isIdentStart(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || b == '_'
}

func isIdentPart(b byte) bool {
	return isIdentStart(b) || (b >= '0' && b <= '9')
}

// applyOverrides returns a copy of the Resume with string fields overridden by overrides map
func applyOverrides(in config.Resume, overrides map[string]string) config.Resume {
	if overrides == nil || len(overrides) == 0 {
		return in
	}
	// copy value
	out := in
	// manually override known top-level string fields
	if v, ok := overrides["Name"]; ok {
		out.Name = v
	}
	if v, ok := overrides["Email"]; ok {
		out.Email = v
	}
	if v, ok := overrides["Phone"]; ok {
		out.Phone = v
	}
	if v, ok := overrides["Address"]; ok {
		out.Address = v
	}
	if v, ok := overrides["Website"]; ok {
		out.Website = v
	}
	if v, ok := overrides["Github"]; ok {
		out.Github = v
	}
	if v, ok := overrides["CompanyToApplyTo"]; ok {
		out.CompanyToApplyTo = v
	}
	if v, ok := overrides["RoleToApplyTo"]; ok {
		out.RoleToApplyTo = v
	}
	return out
}
