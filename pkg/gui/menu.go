package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// todo: figure out either making a struct or general function to handle creating multiple different menus for the windows

func renderMenu(w fyne.Window) *fyne.MainMenu {
	// todo: add menu ability to print generated texts as pdf or word, etc...
	fileMenu := fyne.NewMenu("File",
		// fyne.NewMenuItem("Save", func() { _ = editor.save(w) }),
		// fyne.NewMenuItem("Save As…", func() { _ = editor.saveAs(w) }),
		fyne.NewMenuItem("Quit", func() { w.Close() }),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() { dialog.ShowInformation("About", "Generated Document", w) }),
		fyne.NewMenuItem("Shortcuts", func() {
			dialog.ShowInformation("Shortcuts", "Ctrl+S Save\nCtrl+N New\nCtrl+F Find\nCtrl+U Undo\nCtrl+R Redo", w)
		}),
	)

	return fyne.NewMainMenu(fileMenu, helpMenu)
}
