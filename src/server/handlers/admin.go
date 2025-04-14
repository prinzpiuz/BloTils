package handlers

import (
	"BloTils/src/db"
	"net/http"
	"strconv"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		db_connection := GetDbConnection(r)
		users := db.GetAllUsers(db_connection)
		templateData := setCommonTemplateTata(TemplateData{}, r, w)
		templateData.Data = map[string]any{"users": users}
		generateHTML(w, templateData, "layout", "admin")
	}
}

func ApproveUserRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		db_connection := GetDbConnection(r)
		userId, err := strconv.Atoi(getUrlVars(r, "user_id"))
		commonIntParsingError(w, r, err)
		db.ApproveUser(db_connection, userId)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		db_connection := GetDbConnection(r)
		userId, err := strconv.Atoi(getUrlVars(r, "user_id"))
		commonIntParsingError(w, r, err)
		session := GetSessionData(r)
		if session.User.IsAdmin() && session.User.Id != userId {
			db.DeleteUser(db_connection, userId)
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}
