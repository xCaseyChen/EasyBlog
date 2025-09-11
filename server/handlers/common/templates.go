package common

import (
	"html/template"
	"log"
)

var commonTemplates = map[string]string{
	"info": "templates/info/index.html",
}

var tmpl *template.Template

func init() {
	tmpl = template.New("root")
	for name, path := range commonTemplates {
		_, err := tmpl.New(name).ParseFiles(path)
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", path, err)
		}
	}
}
