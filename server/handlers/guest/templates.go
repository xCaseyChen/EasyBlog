package guest

import (
	"html/template"
	"log"
)

// <templates name, template path> map
var guestTemplates = map[string]string{
	"setup": "templates/setup/index.html",
	"post":  "templates/post/index.html",
	"home":  "templates/home/index.html",
}

// root template for html
var tmpl *template.Template

// init(): templates
func init() {
	tmpl = template.New("root")
	for name, path := range guestTemplates {
		_, err := tmpl.New(name).ParseFiles(path)
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", path, err)
		}
	}
}
