package view

import (
	"html/template"
	"log/slog"
	"os"
)

func ParseTemplates() *template.Template {
	fm := template.FuncMap{
		"Localize": func(messageID string) string { return "" },
	}

	tmpl, err := template.New("").Funcs(fm).ParseGlob("./web/templates/*.gohtml")
	if err != nil {
		slog.Error("failed to parse templates", "error", err)
		os.Exit(1)
	}
	return tmpl
}
