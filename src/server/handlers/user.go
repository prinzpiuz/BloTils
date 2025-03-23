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

const UserHandler = "UserHandler"

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
		confirmPassword := r.FormValue("confirmPassword")
		db_connection := get_db_connection(r)
		isEmailValidated, msg1 := userEmailValidated(db_connection, email)
		isValidPassword, msg2 := checkeckPassword(password, confirmPassword)
		if !isEmailValidated || !isValidPassword {
			templateData := set_common_template_data(TemplateData{}, r, w)
			templateData.Errors = []string{msg1, msg2}
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
		msg := `Account Creation Request Processed Succesfully <br>
								Wait for approval from admin`
		sendMessagePage(r, w, msg)
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

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templateData := set_common_template_data(TemplateData{}, r, w)
		generateHTML(w, templateData, "layout", "forgot_password")
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error Parsing Form: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		email := r.FormValue("email")
		db_connection := get_db_connection(r)
		user := db.GetUser(db_connection, email)
		if user.UserExist() {
			token := generateSecureToken()
			err := db.SetRestToken(db_connection, token, user.Id)
			if err == nil {
				resetPasswordMail(email, token)
				msg := `You'll recive your reset link in mail, if your email is verified`
				sendMessagePage(r, w, msg)
				return
			}
			log.Print("Failed To Save Reset Password Token")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Printf("User %s Not Exist", email)
	}
}

func validateToken(w http.ResponseWriter, r *http.Request, token string) db.PasswordResetToken {
	db_connection := get_db_connection(r)
	tokenObj := db.GetTokenUser(db_connection, token)
	if !tokenObj.IsValid() {
		setHTTPError(w, InvalidToken, UserHandler, http.StatusUnauthorized)
		return db.PasswordResetToken{}
	}
	return tokenObj
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		token := getUrlVars(r, "token")
		tokenObj := validateToken(w, r, token)
		if tokenObj.IsEmpty() {
			return
		}
		templateData := set_common_template_data(TemplateData{}, r, w)
		templateData.Data = map[string]interface{}{"tokenObj": tokenObj}
		generateHTML(w, templateData, "layout", "reset_password")
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error Parsing Form: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token := r.FormValue("reset_token")
		tokenObj := validateToken(w, r, token)
		if tokenObj.IsEmpty() {
			return
		}
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")
		isValidPassword, msg := checkeckPassword(password, confirmPassword)
		if !isValidPassword {
			templateData := set_common_template_data(TemplateData{}, r, w)
			templateData.Errors = []string{msg}
			templateData.Data = map[string]interface{}{"tokenObj": tokenObj}
			generateHTML(w, templateData, "layout", "reset_password")
			return

		}
		passwordHash := passwordHash(password)
		db_connection := get_db_connection(r)
		err = db.DeleteTokenAndSetPassword(db_connection, tokenObj.Token, passwordHash, tokenObj.User.Id)
		if err != nil {
			log.Printf("Failed to update password for user %s", tokenObj.User.Email)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendMessagePage(r, w, "Password Reset Successfull")
		log.Printf("Password Reset Successfully for user %s", tokenObj.User.Email)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templateData := set_common_template_data(TemplateData{}, r, w)
		generateHTML(w, templateData, "layout", "login")
	case http.MethodPost:

	}
}
