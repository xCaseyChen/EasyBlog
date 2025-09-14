package guest

import (
	"html/template"
	"log"
)

// <templates name, template path> map
var guestTemplates = map[string][]string{
	"post":  {"templates/layout_guest/index.html", "templates/post/index.html"},
	"home":  {"templates/layout_guest/index.html", "templates/home/index.html"},
	"about": {"templates/layout_guest/index.html", "templates/about/index.html"},
}

// <template name, *template.Template> map
var tmpl map[string]*template.Template

// init(): templates
func init() {
	tmpl = make(map[string]*template.Template)
	for name, paths := range guestTemplates {
		t, err := template.ParseFiles(paths...)
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", name, err)
		}
		tmpl[name] = t
	}
}
