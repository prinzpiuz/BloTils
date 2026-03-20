package server

import (
	"BloTils/src/db"
	"encoding/json"
	"errors"
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
		return

	case http.MethodPost:
		db_connection := GetDbConnection(r)
		sessionData := GetSessionData(r)
		if sessionData == nil || sessionData.User.Id == 0 {
			SetErrorFlash(w, r, "Please log in to add a domain")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		parseForm(r, w)

		domainName := r.FormValue("domainName")
		if domainName == "" {
			SetErrorFlash(w, r, "Domain name is required")
			http.Redirect(w, r, "/domains", http.StatusSeeOther)
			return
		}

		enableLike := getToggleValues(r, "enableLike")
		enableComment := getToggleValues(r, "enableComments")

		domain := db.DomainFactory(sessionData.User, domainName, enableLike, enableComment, 0)

		err := db.AddDomainAndSettings(db_connection, domain)
		if err != nil {
			log.Printf("Error adding domain: %v", err)

			if errors.Is(err, db.ErrDomainExists) {
				SetErrorFlash(w, r, "Domain already exists")
			} else {
				SetErrorFlash(w, r, "Failed to add domain. Please try again.")
			}
			http.Redirect(w, r, "/domains", http.StatusSeeOther)
			return
		}

		SetSuccessFlash(w, r, "Domain added successfully")
		http.Redirect(w, r, "/domains", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func EditDomainSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionData := GetSessionData(r)
	if sessionData == nil || sessionData.User.Id == 0 {
		SetErrorFlash(w, r, "Please log in")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	domainID, err := strconv.Atoi(getUrlVars(r, "domain_id"))
	if err != nil {
		log.Printf("Error parsing domain_id: %v", err)
		SetErrorFlash(w, r, "Invalid domain ID")
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
		return
	}
	db_connection := GetDbConnection(r)
	err = db.DeleteDomain(db_connection, domainID, sessionData.User.Id)
	if err != nil {
		log.Printf("Error deleting domain %d: %v", domainID, err)
		SetErrorFlash(w, r, "Failed to delete domain")
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
		return
	}

	SetSuccessFlash(w, r, "Domain deleted successfully")
	http.Redirect(w, r, "/domains", http.StatusSeeOther)
}

// DomainDetailPage renders the domain detail page showing all likes for posts under a domain.
func DomainDetailPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
		return
	}
	db_connection := GetDbConnection(r)
	domainId, err := strconv.Atoi(getUrlVars(r, "domain_id"))
	if err != nil {
		commonIntParsingError(w, r, err)
		return
	}
	domain, err := db.GetDomainById(db_connection, domainId)
	if err != nil {
		if errors.Is(err, db.ErrDomainNotFound) {
			SetErrorFlash(w, r, "Domain not found")
			http.Redirect(w, r, "/domains", http.StatusSeeOther)
			return
		}
		log.Printf("Error getting domain by ID %d: %v", domainId, err)
		SetErrorFlash(w, r, "Error retrieving domain details")
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
		return
	}

	likes := db.GetLikesByDomainId(db_connection, domainId)
	templateData := setCommonTemplateData(r, w)
	templateData.Data = map[string]any{
		"domain": domain,
		"likes":  likes,
	}
	generateHTML(w, templateData, "layout", "domain_detail")
}

// DomainLikesTimeline is an API handler that returns the likes timeline data
// for a specific post (URI) within a domain, used for rendering the graph.
func DomainLikesTimeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		setHTTPError(w, ErrMethodNotAllowed, "DomainLikesTimeline", http.StatusMethodNotAllowed)
		return
	}
	db_connection := GetDbConnection(r)
	domainId, err := strconv.Atoi(getUrlVars(r, "domain_id"))
	if err != nil {
		setHTTPError(w, ErrBadRequest, "DomainLikesTimeline", http.StatusBadRequest)
		return
	}
	uri := r.URL.Query().Get("uri")
	if uri == "" {
		setHTTPError(w, ErrBadRequest, "DomainLikesTimeline: missing uri param", http.StatusBadRequest)
		return
	}
	timeline := db.GetLikesTimeline(db_connection, domainId, uri)
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(timeline)
	if err != nil {
		log.Printf("Error Encoding JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
