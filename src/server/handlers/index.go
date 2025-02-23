// Package handlers provides HTTP request handlers for the application.
package handlers

import (
	"log"
	"net/http"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	templateData := set_template_data(TemplateData{}, r)
	err := generateHTML(w, templateData, "layout", "index")
	if err != nil {
		log.Printf("Error Generating HTML: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
