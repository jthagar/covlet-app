package gui

import (
    "covlet/pkg/config"
    "covlet/pkg/internal"
    "fmt"
    "path/filepath"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
)

// TODO: Consider creating a RenderWindow struct that encapsulates window + data accessors
// and provides methods for building menus and toolbars. This will avoid passing closures
// around and make testing easier.

// renderMenu builds the menu for the render window. getText returns the latest rendered text.
func renderMenu(w fyne.Window, getText func() string) *fyne.MainMenu {
    // Export as PDF
    exportPDF := fyne.NewMenuItem("Export as PDF…", func() {
        titleEntry := widget.NewEntry()
        titleEntry.SetPlaceHolder("Document Title")
        dialog.ShowForm("Export as PDF", "Save", "Cancel",
            []*widget.FormItem{{Text: "Title", Widget: titleEntry}}, func(ok bool) {
                if !ok {
                    return
                }
                title := titleEntry.Text
                if title == "" {
                    title = "Document"
                }
                // Default output directory: ~/Downloads/covlet
                dir, err := config.EnsureDownloadsCovletDir()
                if err != nil {
                    dialog.ShowError(fmt.Errorf("could not prepare output directory: %w", err), w)
                    return
                }
                // File name from title
                base := sanitizeFileName(title)
                if base == "" {
                    base = "document"
                }
                out := filepath.Join(dir, base+".pdf")
                if err := internal.SaveTextAsPDF(title, getText(), out); err != nil {
                    dialog.ShowError(fmt.Errorf("failed to export PDF: %w", err), w)
                    return
                }
                dialog.ShowInformation("Saved", fmt.Sprintf("PDF saved to\n%s", out), w)
            }, w)
    })

    fileMenu := fyne.NewMenu("File",
        exportPDF,
        fyne.NewMenuItemSeparator(),
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

// sanitizeFileName provides a very small filter for creating file names from titles.
// It is not exhaustive but good enough for simple cases.
func sanitizeFileName(s string) string {
    r := make([]rune, 0, len(s))
    for _, ch := range s {
        switch ch {
        case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
            r = append(r, '_')
        default:
            r = append(r, ch)
        }
    }
    // trim spaces
    out := string(r)
    return out
}
