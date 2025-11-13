package gui

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"cover-letter-templates/pkg/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Run() error {
	a := app.New()
	w := a.NewWindow("Covlet")
	windowSize := fyne.NewSize(1000, 700)
	w.Resize(windowSize)

	// Load or initialize config with home_dir and ensure templates subfolder
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		// If file doesn't exist or fails to parse, start with empty config and initialize
		cfg = &config.Config{}
	}

	if cfg.HomeDir == "" || !dirExists(cfg.HomeDir) {
		if err := promptForHomeDir(w, cfg); err != nil {
			// If the user cancels, still show the app, but many actions will prompt again
			log.Println("home_dir not set:", err)
		}
	}
	// Ensure templates subfolder under home_dir
	homeTemplates := ""
	if cfg.HomeDir != "" {
		homeTemplates = filepath.Join(cfg.HomeDir, "templates")
		_ = os.MkdirAll(homeTemplates, 0755)
	}

 state := &editorState{
		cfg:           cfg,
		homeTemplates: homeTemplates,
	}

	// Build editor + preview split view and toolbar wired to state
	editor, preview, topBar := buildEditorUI(w, a, state)
	state.editor = editor
	state.preview = preview

	editorContainer := container.NewHSplit(editor, preview)
	editorContainer.Offset = 0.55

	content := container.NewBorder(topBar, nil, nil, nil, editorContainer)
	w.SetContent(content)
	w.ShowAndRun()
	return nil
}

// promptForHomeDir asks the user to select a folder for home_dir and saves config.yml.
func promptForHomeDir(w fyne.Window, cfg *config.Config) error {
	// Offer a folder open dialog
	var completed = make(chan struct{})
	var retErr error
	dlg := dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
		defer close(completed)
		if err != nil {
			retErr = err
			return
		}
		if list == nil {
			retErr = fmt.Errorf("selection canceled")
			return
		}
		p := list.Path()
		if p == "" {
			retErr = fmt.Errorf("invalid directory")
			return
		}
		if err := os.MkdirAll(p, 0755); err != nil {
			retErr = err
			return
		}
		cfg.HomeDir = p
		if err := cfg.SaveConfig("config.yml"); err != nil {
			retErr = err
			return
		}
		// create templates subfolder
		_ = os.MkdirAll(filepath.Join(cfg.HomeDir, "templates"), 0755)
		dialog.ShowInformation("Initialized", fmt.Sprintf("home_dir set to: %s", cfg.HomeDir), w)
	}, w)
	// Try to default to user home
	if uhome, _ := os.UserHomeDir(); uhome != "" {
		if l, err := storage.ListerForURI(storage.NewFileURI(uhome)); err == nil {
			dlg.SetLocation(l)
		}
	}
	dlg.Resize(fyne.NewSize(600, 400))
	dlg.Show()
	<-completed
	return retErr
}

// state for the editor
 type editorState struct {
	cfg           *config.Config
	homeTemplates string
	editor        *widget.Entry
	preview      *widget.Entry
}

func buildEditorUI(w fyne.Window, a fyne.App, st *editorState) (*widget.Entry, *widget.Entry, *widget.Toolbar) {
	// Editor
	editor := widget.NewMultiLineEntry()
	editor.SetPlaceHolder("Enter Go template here, e.g., {{ .Name }}")
	editor.Wrapping = fyne.TextWrapWord
	editor.OnChanged = func(s string) {
		if len(s) > 0 {
			lastChar := s[len(s)-1:]
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

	// Preview (read-only)
	preview := widget.NewMultiLineEntry()
	preview.SetPlaceHolder("Render preview will appear here")
	preview.Wrapping = fyne.TextWrapWord
	preview.Disable()

	// Shortcuts
	ctrlS := desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlS, func(shortcut fyne.Shortcut) {
		if err := guiSave(w, st, editor.Text); err != nil {
			dialog.ShowError(err, w)
		}
	})
	ctrlO := desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlO, func(shortcut fyne.Shortcut) { guiOpen(w, st, editor) })
	ctrlShiftS := desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierShift | fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlShiftS, func(shortcut fyne.Shortcut) {
		guiSaveAs(w, st, editor.Text)
	})
	ctrlU := desktop.CustomShortcut{KeyName: fyne.KeyU, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlU, func(shortcut fyne.Shortcut) { undo(editor) })
	ctrlR := desktop.CustomShortcut{KeyName: fyne.KeyR, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(&ctrlR, func(shortcut fyne.Shortcut) { redo(editor) })

	// Toolbar actions
	renderAction := widget.NewToolbarAction(theme.MediaPlayIcon(), func() { guiRender(w, st, editor.Text, preview) })
	openAction := widget.NewToolbarAction(theme.FolderOpenIcon(), func() { guiOpen(w, st, editor) })
	saveAction := widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		if err := guiSave(w, st, editor.Text); err != nil {
			dialog.ShowError(err, w)
		}
	})
 saveAsAction := widget.NewToolbarAction(theme.DocumentCreateIcon(), func() { guiSaveAs(w, st, editor.Text) })
	settingsAction := widget.NewToolbarAction(theme.SettingsIcon(), func() {
		if err := promptForHomeDir(w, st.cfg); err != nil {
			dialog.ShowError(err, w)
			return
		}
		st.homeTemplates = filepath.Join(st.cfg.HomeDir, "templates")
		_ = os.MkdirAll(st.homeTemplates, 0755)
	})
	help := widget.NewToolbarAction(theme.HelpIcon(), func() {
		dialog.ShowInformation("Help", "Type templates using Go text/template syntax. Click Render to preview using data from config.yml.", w)
	})
	menu := widget.NewToolbar(openAction, saveAction, saveAsAction, renderAction, settingsAction, &widget.ToolbarSeparator{}, help)

	return editor, preview, menu
}

// Render current editor content to preview using cfg.Resume
func guiRender(w fyne.Window, st *editorState, tplText string, preview *widget.Entry) {
	if st.cfg == nil {
		dialog.ShowError(fmt.Errorf("config not loaded"), w)
		return
	}
	// If home_dir is not set, prompt
	if st.cfg.HomeDir == "" || !dirExists(st.cfg.HomeDir) {
		if err := promptForHomeDir(w, st.cfg); err != nil {
			dialog.ShowError(err, w)
			return
		}
		st.homeTemplates = filepath.Join(st.cfg.HomeDir, "templates")
		_ = os.MkdirAll(st.homeTemplates, 0755)
	}
	// Execute template
	t, err := template.New("editor").Parse(tplText)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	var out bytes.Buffer
	if err := t.Execute(&out, st.cfg.Resume); err != nil {
		dialog.ShowError(err, w)
		return
	}
	preview.SetText(out.String())
}

// Open a template file under homeTemplates
func guiOpen(w fyne.Window, st *editorState, editor *widget.Entry) {
	if st.cfg.HomeDir == "" || !dirExists(st.cfg.HomeDir) {
		if err := promptForHomeDir(w, st.cfg); err != nil {
			dialog.ShowError(err, w)
			return
		}
		st.homeTemplates = filepath.Join(st.cfg.HomeDir, "templates")
		_ = os.MkdirAll(st.homeTemplates, 0755)
	}

	dlg := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if r == nil {
			return
		}
		defer r.Close()
		data, err := io.ReadAll(r)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		editor.SetText(string(data))
		currentPath = r.URI().Path()
	}, w)
	// Filter extensions
	dlg.SetFilter(storage.NewExtensionFileFilter([]string{".tpl", ".tmpl", ".gotmpl"}))
	// Limit to homeTemplates as starting location
	if st.homeTemplates != "" {
		if l, err := storage.ListerForURI(storage.NewFileURI(st.homeTemplates)); err == nil {
			dlg.SetLocation(l)
		}
	}
	dlg.Resize(fyne.NewSize(800, 600))
	dlg.Show()
}

// Save current content; if no currentPath, falls back to Save As
func guiSave(w fyne.Window, st *editorState, text string) error {
	if currentPath == "" {
		guiSaveAs(w, st, text)
		return nil
	}
	// ensure file remains within homeTemplates (defensive)
	if st.homeTemplates != "" {
		if !strings.HasPrefix(filepath.Clean(currentPath), filepath.Clean(st.homeTemplates)+string(os.PathSeparator)) && filepath.Clean(currentPath) != filepath.Clean(st.homeTemplates) {
			return fmt.Errorf("saving outside home_dir/templates is not allowed")
		}
	}
	return os.WriteFile(currentPath, []byte(text), 0644)
}

func guiSaveAs(w fyne.Window, st *editorState, text string) {
	if st.cfg.HomeDir == "" || !dirExists(st.cfg.HomeDir) {
		if err := promptForHomeDir(w, st.cfg); err != nil {
			dialog.ShowError(err, w)
			return
		}
		st.homeTemplates = filepath.Join(st.cfg.HomeDir, "templates")
		_ = os.MkdirAll(st.homeTemplates, 0755)
	}
	var completed = make(chan struct{})
	dlg := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		defer close(completed)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if uc == nil {
			return
		}
		p := uc.URI().Path()
		// Enforce saving under homeTemplates
		if st.homeTemplates != "" {
			cleanHome := filepath.Clean(st.homeTemplates)
			cleanP := filepath.Clean(p)
			if !strings.HasPrefix(cleanP, cleanHome+string(os.PathSeparator)) && cleanP != cleanHome {
				dialog.ShowError(fmt.Errorf("please save within %s", st.homeTemplates), w)
				return
			}
		}
		if !strings.HasSuffix(strings.ToLower(p), ".tpl") && !strings.HasSuffix(strings.ToLower(p), ".tmpl") && !strings.HasSuffix(strings.ToLower(p), ".gotmpl") {
			p = p + ".tpl"
		}
		if err := os.WriteFile(p, []byte(text), 0644); err != nil {
			dialog.ShowError(err, w)
			return
		}
		currentPath = p
		uc.Close()
	}, w)
	// Default location to homeTemplates
	if st.homeTemplates != "" {
		if l, err := storage.ListerForURI(storage.NewFileURI(st.homeTemplates)); err == nil {
			dlg.SetLocation(l)
		}
	}
	dlg.SetFileName("untitled.tpl")
	dlg.Resize(fyne.NewSize(800, 600))
	dlg.Show()
}

// helpers
func dirExists(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}
