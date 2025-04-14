package handlers

import (
	"BloTils/src/db"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetUserDomains(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		db_connection := GetDbConnection(r)
		sessionData := GetSessionData(r)
		domains := db.GetAllUserDomains(db_connection, sessionData.User.Id)
		templateData := setCommonTemplateTata(TemplateData{}, r, w)
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
	case http.MethodGet:
		db_connection := GetDbConnection(r)
		domainId, err := strconv.Atoi(getUrlVars(r, "domain_id"))
		commonIntParsingError(w, r, err)
		domain := db.GetDomainById(db_connection, domainId)
		templateData := setCommonTemplateTata(TemplateData{}, r, w)
		templateData.Data = map[string]any{"domain": domain, "edit_page": true}
		generateHTML(w, templateData, "layout", "edit_domain_details")
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
		redirectUrl := fmt.Sprintf("/domain/%v", domainId)
		http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	}
}


func DeleteDomain(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		db_connection := GetDbConnection(r)
		domainId, err := strconv.Atoi(getUrlVars(r, "domain_id"))
		commonIntParsingError(w, r, err)
		db.DeleteDomain(db_connection, domainId)
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
	}
}
