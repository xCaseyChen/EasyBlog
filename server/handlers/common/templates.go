package common

import (
	"html/template"
	"log"
)

var commonTemplates = map[string][]string{
	"info":  {"templates/layout_common/index.html", "templates/info/index.html"},
	"setup": {"templates/layout_common/index.html", "templates/setup/index.html"},
}

var tmpl map[string]*template.Template

func init() {
	tmpl = make(map[string]*template.Template)
	for name, paths := range commonTemplates {
		t, err := template.ParseFiles(paths...)
		if err != nil {
			log.Fatalf("Failed to parse template %v: %v", paths, err)
		}
		tmpl[name] = t
	}
}
