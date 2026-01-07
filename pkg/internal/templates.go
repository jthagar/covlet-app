package internal

import (
	"bytes"
	"strings"
	"text/template"

	"codeberg.org/go-pdf/fpdf"
)

func wrapText(s string) string {
	return "\n" + s + "\n"
}

func RenderEditor(t *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// SaveTextAsPDF renders the provided plain text into a simple PDF file.
// It performs basic word wrapping and supports multiple pages.
func SaveTextAsPDF(title, text, outPath string) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(title, true)
	pdf.SetAuthor("Covlet", true)
	pdf.AddPage()

	// Margins and font
	left, top, right := 20.0, 20.0, 20.0
	pdf.SetMargins(left, top, right)
	pdf.SetAutoPageBreak(true, 20.0)
	pdf.SetFont("Helvetica", "", 12)

	// Title
	if strings.TrimSpace(title) != "" {
		pdf.SetFont("Helvetica", "B", 14)
		pdf.CellFormat(0, 8, title, "", 1, "L", false, 0, "")
		pdf.Ln(2)
		pdf.SetFont("Helvetica", "", 12)
	}

	// Body text with MultiCell for wrapping
	// Use 0 width to make it wrap to the page width minus margins
	pdf.MultiCell(0, 6, text, "", "L", false)

	return pdf.OutputFileAndClose(outPath)
}
