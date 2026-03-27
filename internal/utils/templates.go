package utils

import (
	"html/template"
	"net/http"
)

var templates *template.Template

// CarregarTemplates insere os templates HTML na variável templates
func CarregarTemplates() {
	templates = template.Must(template.ParseGlob("../../website/*.html"))
}

// ExecutarTemplate renderiza uma página HTML na tela
func ExecutarTemplate(w http.ResponseWriter, template string, dados interface{}) {
	templates.ExecuteTemplate(w, template+".html", dados)
}
