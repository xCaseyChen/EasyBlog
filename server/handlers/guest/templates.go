package guest

import (
	"html/template"
	"log"
)

// <templates name, template path> map
var guestTemplates = map[string]string{
	"setup": "templates/setup/index.html",
	"post":  "templates/post/index.html",
}

// root template for html
var Template *template.Template

// init(): templates
func init() {
	Template = template.New("root")
	for name, path := range guestTemplates {
		_, err := Template.New(name).ParseFiles(path)
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", path, err)
		}
	}
}
