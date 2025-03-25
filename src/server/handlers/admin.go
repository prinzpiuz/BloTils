package handlers

import (
	"BloTils/src/db"
	"net/http"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
	db_connection := GetDbConnection(r)
	sessionData := GetSessionData(r)
	users := db.GetAllUsers(db_connection)
	templateData := setCommonTemplateTata(TemplateData{}, r, w)
	templateData.Session = sessionData
	templateData.Data = map[string]interface{}{"users": users}
	generateHTML(w, templateData, "layout", "admin")
}
