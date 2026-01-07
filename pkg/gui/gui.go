package gui

import (
	"covlet/pkg/config"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type smallTheme struct{ base fyne.Theme }

func (s *smallTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return s.base.Color(n, v)
}
func (s *smallTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return s.base.Icon(n) }
func (s *smallTheme) Font(f fyne.TextStyle) fyne.Resource     { return s.base.Font(f) }
func (s *smallTheme) Size(n fyne.ThemeSizeName) float32 {
	sz := s.base.Size(n)
	if n == theme.SizeNameText {
		if sz > 1 {
			return sz - 1
		}
	}
	return sz
}

func Run() error {
	config.InitMainDir()
	a := app.New()
	// apply smaller text theme (one size smaller)
	a.Settings().SetTheme(&smallTheme{base: theme.DefaultTheme()})

	w := a.NewWindow("Covlet")
	windowSize := fyne.NewSize(1000, 700)
	w.Resize(windowSize)

	// Build editor
	editorContainer, editor := NewEditor(w)

	// ensure templates directory exists and helper to compute left/right roots
	_, _ = config.EnsureTemplatesDir()

	// trees and panes
	lRoot, rRoot := computeRoots()
	leftTree := fileTree(lRoot, editor, w)
	rightTree := fileTree(rRoot, editor, w)

	leftPane := container.NewBorder(widget.NewLabel("Templates Left"), nil, nil, nil, container.NewVScroll(leftTree))
	rightPane := container.NewBorder(widget.NewLabel("Templates Right"), nil, nil, nil, container.NewVScroll(rightTree))

	// dynamic resizable split panes
	centerSplit := container.NewHSplit(leftPane, editorContainer)
	centerSplit.Offset = 0.25
	mainSplit := container.NewHSplit(centerSplit, rightPane)
	mainSplit.Offset = 0.75

	// Build main menu bar
	toggleLeft := true
	toggleRight := true

	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("New", func() { editor.new(w) }),
		fyne.NewMenuItem("Open Folder…", func() {
			dlg := dialog.NewFolderOpen(func(listable fyne.ListableURI, err error) {
				if err != nil {
					log.Println("Folder open error:", err)
					return
				}
				if listable == nil {
					return
				}
				if err := config.SetMainDir(listable.Path()); err != nil {
					dialog.ShowError(err, w)
					return
				}
				_, _ = config.EnsureTemplatesDir()
				lRoot, rRoot := computeRoots()
				leftTree.Root = lRoot
				rightTree.Root = rRoot
				leftTree.Refresh()
				rightTree.Refresh()
			}, w)
			dlg.Show()
		}),
		fyne.NewMenuItem("Save", func() { _ = editor.save(w) }),
		fyne.NewMenuItem("Save As…", func() { _ = editor.saveAs(w) }),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() { w.Close() }),
	)

	editMenu := fyne.NewMenu("Edit",
		fyne.NewMenuItem("Undo", func() { editor.undo() }),
		fyne.NewMenuItem("Redo", func() { editor.redo() }),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Cut", func() {
			if c := w.Canvas(); c != nil {
				if foc, ok := c.Focused().(*widget.Entry); ok {
					foc.TypedShortcut(&fyne.ShortcutCut{Clipboard: w.Clipboard()})
				}
			}
		}),
		fyne.NewMenuItem("Copy", func() {
			if c := w.Canvas(); c != nil {
				if foc, ok := c.Focused().(*widget.Entry); ok {
					foc.TypedShortcut(&fyne.ShortcutCopy{Clipboard: w.Clipboard()})
				}
			}
		}),
		fyne.NewMenuItem("Paste", func() {
			if c := w.Canvas(); c != nil {
				if foc, ok := c.Focused().(*widget.Entry); ok {
					foc.TypedShortcut(&fyne.ShortcutPaste{Clipboard: w.Clipboard()})
				}
			}
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Find", func() { editor.find(w) }),
	)

	viewMenu := fyne.NewMenu("View",
		fyne.NewMenuItem("Toggle Left Sidebar", func() {
			toggleLeft = !toggleLeft
			if toggleLeft {
				leftPane.Show()
			} else {
				leftPane.Hide()
			}
		}),
		fyne.NewMenuItem("Toggle Right Sidebar", func() {
			toggleRight = !toggleRight
			if toggleRight {
				rightPane.Show()
			} else {
				rightPane.Hide()
			}
		}),
		fyne.NewMenuItem("Reset Panes", func() {
			leftPane.Show()
			rightPane.Show()
			centerSplit.Offset = 0.25
			mainSplit.Offset = 0.75
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Light Theme", func() { a.Settings().SetTheme(&smallTheme{base: theme.LightTheme()}) }),
		fyne.NewMenuItem("Dark Theme", func() { a.Settings().SetTheme(&smallTheme{base: theme.DarkTheme()}) }),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() { dialog.ShowInformation("About", "Covlet - Cover Letter Templates Editor", w) }),
		fyne.NewMenuItem("Shortcuts", func() {
			dialog.ShowInformation("Shortcuts", "Ctrl+S Save\nCtrl+N New\nCtrl+F Find\nCtrl+U Undo\nCtrl+R Redo", w)
		}),
	)

	w.SetMainMenu(fyne.NewMainMenu(fileMenu, editMenu, viewMenu, helpMenu))

	w.SetContent(mainSplit)
	w.ShowAndRun()
	return nil
}

// fileTree builds a filesystem tree rooted at the given directory and loads
// file content into the provided editor on selection.
func fileTree(root string, e *TextEditor, w fyne.Window) *widget.Tree {
	child := func(uid string) []string {
		// uid is a full path
		entries, err := os.ReadDir(uid)
		if err != nil {
			return []string{}
		}
		var ids []string
		for _, e := range entries {
			// hide dot files
			name := e.Name()
			if len(name) > 0 && name[0] == '.' {
				continue
			}
			ids = append(ids, filepath.Join(uid, name))
		}
		return ids
	}
	isBranch := func(uid string) bool {
		fi, err := os.Stat(uid)
		if err != nil {
			return false
		}
		return fi.IsDir()
	}
	var t *widget.Tree
	create := func(branch bool) fyne.CanvasObject {
		name := widget.NewLabel("")
		more := widget.NewButtonWithIcon("", theme.MoreVerticalIcon(), func() {})
		more.Importance = widget.LowImportance
		c := container.New(layout.NewHBoxLayout(), name, layout.NewSpacer(), more)
		return c
	}
	update := func(uid string, branch bool, obj fyne.CanvasObject) {
		c := obj.(*fyne.Container)
		name := c.Objects[0].(*widget.Label)
		more := c.Objects[2].(*widget.Button)
		name.SetText(filepath.Base(uid))
		if branch {
			more.Show()
			// bind click to create a file in this directory
			dir := uid
			more.OnTapped = func() {
				showCreateFileDialog(w, dir, t)
			}
		} else {
			more.Hide()
		}
	}

	// todo: figure out how to create animation showing brief selection and highlight of a file option
	t = widget.NewTree(child, isBranch, create, update)
	t.Root = root
	// Filter only template-like files for selection
	t.OnSelected = func(uid string) {
		fi, err := os.Stat(uid)
		if err != nil || fi.IsDir() {
			return
		}
		// only open known text types (.tpl, .tmpl, .txt, .md, .gohtml)
		ext := filepath.Ext(uid)
		switch ext {
		case ".tpl", ".tmpl", ".txt", ".md", ".gohtml", ".html":
			b, err := os.ReadFile(uid)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			e.editor.SetText(string(b))
			e.currentPath = uid
			// refresh variable sidebar to reflect newly loaded content
			e.refreshVarSidebar()
		}
		// add in auto-unselect to allow options to be re-selected
		// this also should mean any unsaved editor changes will get overwritten even if the same file gets re-selected
		t.Unselect(uid)
	}

	return t
}

// showCreateFileDialog prompts for a new file name in the given directory and creates it.
func showCreateFileDialog(w fyne.Window, dir string, t *widget.Tree) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("new_file.tpl")
	dialog.ShowForm("Create New File", "Create", "Cancel",
		[]*widget.FormItem{{Text: "File name", Widget: nameEntry}},
		func(confirm bool) {
			if !confirm {
				return
			}
			name := strings.TrimSpace(nameEntry.Text)
			if name == "" {
				dialog.ShowError(fmt.Errorf("file name cannot be empty"), w)
				return
			}
			if strings.ContainsAny(name, "/\\") {
				dialog.ShowError(fmt.Errorf("file name must not contain path separators"), w)
				return
			}
			// default extension if none provided
			if !strings.Contains(name, ".") {
				name += ".tpl"
			}
			full := filepath.Join(dir, name)
			if _, err := os.Stat(full); err == nil {
				dialog.ShowError(fmt.Errorf("file already exists: %s", name), w)
				return
			}
			if err := os.WriteFile(full, []byte(""), 0o644); err != nil {
				dialog.ShowError(err, w)
				return
			}
			// refresh tree and reveal new file
			t.OpenBranch(dir)
			t.Refresh()
			// optionally select the new file to make it visible
			t.Select(full)
		}, w)
}
