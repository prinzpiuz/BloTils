package server

import (
	"BloTils/src/db"
	"log"
	"net/http"
	"strconv"
)

func GetUserDomains(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		db_connection := GetDbConnection(r)
		sessionData := GetSessionData(r)
		domains := db.GetAllUserDomains(db_connection, sessionData.User.Id)
		templateData := setCommonTemplateData(r, w)
		templateData.Data = map[string]any{"domains": domains}
		generateHTML(w, templateData, "layout", "domains")
	} else {
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
	}

}

func AddDomain(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
	case http.MethodPost:
		db_connection := GetDbConnection(r)
		sessionData := GetSessionData(r)
		parseForm(r, w)
		domainName := r.FormValue("domainName")
		enableLike := getToggleValues(r, "enableLike")
		enableComment := getToggleValues(r, "enableComments")
		domain := db.DomainFactory(sessionData.User, domainName, enableLike, enableComment, 0)
		db.AddDomainAndSettings(db_connection, domain)
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
	}

}

func EditDomainSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		db_connection := GetDbConnection(r)
		domainId, err := strconv.Atoi(getUrlVars(r, "domain_id"))
		commonIntParsingError(w, r, err)
		sessionData := GetSessionData(r)
		parseForm(r, w)
		domainIdFromForm, err := strconv.Atoi(r.FormValue("domainId"))
		commonIntParsingError(w, r, err)
		if domainId != domainIdFromForm {
			log.Printf("Domain ID mismatch: %v != %v", domainId, domainIdFromForm)
		} else {
			domainName := r.FormValue("domainName")
			enableLike := getToggleValues(r, "enableLike")
			enableComment := getToggleValues(r, "enableComments")
			domain := db.DomainFactory(sessionData.User, domainName, enableLike, enableComment, domainId)
			db.UpdateDomainDetails(db_connection, domain)
		}
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
	}
}

func DeleteDomain(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		db_connection := GetDbConnection(r)
		domainId, err := strconv.Atoi(getUrlVars(r, "domain_id"))
		commonIntParsingError(w, r, err)
		err = db.DeleteDomain(db_connection, domainId)
		if err != nil {
			log.Printf("Error Deleting Domain %d: %v", domainId, err)
			// TODO set message
		}
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
	}
}
