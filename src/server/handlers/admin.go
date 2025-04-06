package handlers

import (
	"BloTils/src/db"
	"net/http"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
	db_connection := GetDbConnection(r)
	users := db.GetAllUsers(db_connection)
	templateData := setCommonTemplateTata(TemplateData{}, r, w)
	templateData.Data = map[string]any{"users": users}
	generateHTML(w, templateData, "layout", "admin")
}
