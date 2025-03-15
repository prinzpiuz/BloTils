package handlers

import (
	"BloTils/src/db"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"

	"golang.org/x/term"
)

// CreateAccountPage renders the HTML template for the account creation page.
func CreateAccountPage(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		templateData := set_common_template_data(TemplateData{}, r, w)
		generateHTML(w, templateData, "layout", "create_account")
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error Parsing Form: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		db_connection := get_db_connection(r)
		isEmailValidated, msg1 := userEmailValidated(db_connection, email)
		isValidPassword, msg2 := checkeckPassword(password, confirmPassword)
		errorMessages := []string{msg1, msg2}
		if !isEmailValidated || !isValidPassword {
			userTemplateData := errorAndMessages{ErrorMsg: errorMessages}
			templateData := set_common_template_data(TemplateData{}, r, w)
			templateData.Data = userTemplateData
			generateHTML(w, templateData, "layout", "create_account")
			return
		}
		passwordHash := passwordHash(password)
		err = db.AddUserRequest(db_connection, email, passwordHash)
		if err != nil {
			msg := fmt.Sprintf("Error: Adding user request %s to DB failed", email)
			log.Printf(msg)
			http.Error(w, msg, http.StatusInternalServerError)
		}
		log.Printf("Succesfully Added user request %s", email)
		// success message
		//sucess page

	}
}

func userEmailValidated(db_connection *sql.DB, email string) (bool, string) {
	if !isValidEmail(email) {
		return false, "Error: Invalid Email"
	}
	if db.GetUser(db_connection, email).UserExist() {
		return false, "Error: Email Already Exists"
	}
	return true, ""
}

func checkeckPassword(password string, confirm string) (bool, string) {
	if password != confirm {
		msg := "Error: Passwords do not match"
		log.Println(msg)
		return false, msg
	}
	if len(password) < 8 {
		msg := "Error: Password must be at least 8 characters long"
		log.Println(msg)
		return false, msg
	}
	return true, ""
}

func CreateAdmin(db_connection *sql.DB) {
	fmt.Println("Create Admin")
	var (
		email    string
		password string
	)
	fmt.Print("Email: ")
	fmt.Scanln(&email)
	isEmailValidated, msg := userEmailValidated(db_connection, email)
	if !isEmailValidated {
		fmt.Println("Note: Create A Strong Password")
		fmt.Print("Password: ")
		bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
		password = string(bytePassword)
		fmt.Println()
		fmt.Print("Confirm Password: ")
		byteConfirm, _ := term.ReadPassword(int(syscall.Stdin))
		confirm := string(byteConfirm)
		fmt.Println()
		isValidPassword, msg := checkeckPassword(password, confirm)
		if !isValidPassword {
			fmt.Println(msg)
			os.Exit(1)
		}
		passwordHash := passwordHash(password)
		err := db.CreatSuperUser(db_connection, email, passwordHash)
		if err != nil {
			fmt.Printf("Error: Adding admin user %s to DB failed", email)
		}
		fmt.Printf("Succesfully Created Admin User %s", email)
	} else {
		fmt.Printf("Email Validations Failed\n,Error: %s", msg)
	}

}
