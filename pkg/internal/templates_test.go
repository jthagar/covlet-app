package internal

import (
    "os"
    "path/filepath"
    "testing"
    "text/template"
)

func TestRenderEditor(t *testing.T) {
    tpl, err := template.New("t").Parse("Hello, {{.Name}}!")
    if err != nil {
        t.Fatalf("parse template: %v", err)
    }
    out, err := RenderEditor(tpl, map[string]string{"Name": "World"})
    if err != nil {
        t.Fatalf("RenderEditor error: %v", err)
    }
    if string(out) != "Hello, World!" {
        t.Fatalf("unexpected output: %q", string(out))
    }
}

func TestSaveTextAsPDF(t *testing.T) {
    dir := t.TempDir()
    out := filepath.Join(dir, "doc.pdf")
    if err := SaveTextAsPDF("Title", "Some text body", out); err != nil {
        t.Fatalf("SaveTextAsPDF error: %v", err)
    }
    st, err := os.Stat(out)
    if err != nil {
        t.Fatalf("expected file to exist: %v", err)
    }
    if st.Size() == 0 {
        t.Fatalf("expected non-empty pdf file")
    }
}
