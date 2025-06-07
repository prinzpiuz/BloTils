// Package handlers provides HTTP request handlers for the application.
package server

import (
	"net/http"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	templateData := setCommonTemplateData(r, w)
	templateData.Data = map[string]interface{}{"indexPage": true}
	generateHTML(w, templateData, "layout", "index")

}
