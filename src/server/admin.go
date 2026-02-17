package server

import (
	"BloTils/src/db"
	mailer "BloTils/src/email"
	"log"
	"net/http"
	"strconv"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		db_connection := GetDbConnection(r)
		users := db.GetAllUsers(db_connection)
		templateData := setCommonTemplateData(r, w)
		templateData.Data = map[string]any{"users": users}
		generateHTML(w, templateData, "layout", "admin")
	}
}

func ApproveUserRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		db_connection := GetDbConnection(r)
		userId, err := strconv.Atoi(getUrlVars(r, "user_id"))
		commonIntParsingError(w, r, err)
		email, err := db.ApproveUser(db_connection, userId)
		if err != nil {
			log.Printf("Error Approving User %d: %v", userId, err)
			SetErrorFlash(w, r, "Error Approving User")
		}
		SetSuccessFlash(w, r, "User Request Approved")
		mailer.SendAccountApprovedEmail(email, email)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		db_connection := GetDbConnection(r)
		userId, err := strconv.Atoi(getUrlVars(r, "user_id"))
		commonIntParsingError(w, r, err)
		session := GetSessionData(r)
		if session.User.IsAdmin() && session.User.Id != userId {
			err = db.DeleteUser(db_connection, userId)
			if err != nil {
				log.Printf("Error Deleting User %d: %v", userId, err)
				SetErrorFlash(w, r, "Error Deleting User")
			}
		}
		SetSuccessFlash(w, r, "User Deleted Successfully")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}
