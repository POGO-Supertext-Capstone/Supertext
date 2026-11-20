package web

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
)

var templates = template.Must(template.ParseFS(embedFS, "templates/*.tmpl.html"))

var devTemplates = false

func UseDevTemplates(dev bool) {
	devTemplates = dev
}

func loadTemplates() (*template.Template, error) {
	if devTemplates {
		return template.ParseFS(os.DirFS("web"), "templates/*.tmpl.html")
	}
	return templates, nil
}

// render executes the named template into a buffer first, so a template error
// returns a clean 500 instead of a half-written page.
func render(w http.ResponseWriter, name string, data any) {
	tmpl, err := loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
