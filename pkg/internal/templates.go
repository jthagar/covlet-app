package internal

import (
	"bytes"
	"text/template"
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
