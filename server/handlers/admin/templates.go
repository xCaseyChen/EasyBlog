package admin

import (
	"html/template"
	"log"
)

var adminTemplates = map[string][]string{
	"dashboard": {"templates/layout_admin/index.html", "templates/dashboard/index.html"},
	"preview":   {"templates/layout_preview/index.html", "templates/post/index.html"},
}

var tmpl map[string]*template.Template

func init() {
	tmpl = make(map[string]*template.Template)
	for name, paths := range adminTemplates {
		t, err := template.ParseFiles(paths...)
		if err != nil {
			log.Fatalf("Failed to parse template %v: %v", paths, err)
		}
		tmpl[name] = t
	}
}
