package view

import (
	"html/template"
	"log"
)

func ParseTemplates() *template.Template {
	fm := template.FuncMap{
		"Localize": func(messageID string) string { return "" },
	}

	tmpl, err := template.New("").Funcs(fm).ParseGlob("./web/templates/*.gohtml")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}
	return tmpl
}
