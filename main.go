package main

import (
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Project represents a single project entry.
type Project struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
}

func main() {
	// err := cli.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}

	a := app.New()
	w := a.NewWindow("Covlet")
	w.Resize(fyne.NewSize(640, 480))

	message := canvas.NewText("Hello", color.White)
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New document")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)
	button := widget.NewButton("Update", func() {
		formatted := time.Now().Format("Time: 03:04:05")
		message.Text = formatted
		message.Refresh()
	})

	left := widget.NewLabel("Left")
	content := container.NewBorder(toolbar, nil, left, nil, button, message)
	w.SetContent(content)
	w.ShowAndRun()
}
