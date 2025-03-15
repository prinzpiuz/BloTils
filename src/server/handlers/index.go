// Package handlers provides HTTP request handlers for the application.
package handlers

import (
	"net/http"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	templateData := set_common_template_data(TemplateData{}, r, w)
	generateHTML(w, templateData, "layout", "index")

}
