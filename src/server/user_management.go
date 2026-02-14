package server

import (
	"BloTils/src/db"
	mailer "BloTils/src/email"
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
	templateData := setCommonTemplateData(r, w)
	templateData.Data = map[string]interface{}{"createAccountPage": true}
	switch r.Method {
	case http.MethodGet:
		generateHTML(w, templateData, "layout", "create_account")
	case http.MethodPost:
		parseForm(r, w)
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")
		db_connection := GetDbConnection(r)
		isEmailValidated, msg1 := userEmailValidated(db_connection, email)
		isValidPassword, msg2 := checkeckPassword(password, confirmPassword)
		if !isEmailValidated || !isValidPassword {
			SetErrorFlash(w, r, msg1)
			SetErrorFlash(w, r, msg2)
			http.Redirect(w, r, "/create_account", http.StatusSeeOther)
			return
		}
		passwordHash := passwordHash(password)
		err := db.AddUserRequest(db_connection, email, passwordHash)
		if err != nil {
			msg := fmt.Sprintf("Error: Adding user request %s to DB failed", email)
			log.Print(msg)
			http.Error(w, msg, http.StatusInternalServerError)
		}
		log.Printf("Successfully Added user request %s", email)
		msg := `Account Creation Request Processed Successfully <br>
								Wait for approval from admin`
		SetSuccessFlash(w, r, msg)
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
	_, err := fmt.Scanln(&email)
	if err != nil {
		log.Fatal(err)
	}
	isEmailValidated, msg := userEmailValidated(db_connection, email)
	if isEmailValidated {
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
		fmt.Printf("Successfully Created Admin User %s", email)
	} else {
		fmt.Printf("Email Validations Failed\n,Error: %s", msg)
	}

}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templateData := setCommonTemplateData(r, w)
		generateHTML(w, templateData, "layout", "forgot_password")
	case http.MethodPost:
		parseForm(r, w)
		email := r.FormValue("email")
		db_connection := GetDbConnection(r)
		user := db.GetUser(db_connection, email)
		if user.UserExist() {
			token := generateSecureToken()
			err := db.SetRestToken(db_connection, token, user.Id)
			if err == nil {
				mailer.ResetPasswordMail(email, token)
				msg := `You'll recive your reset link in mail, if your email is verified`
				SetSuccessFlash(w, r, msg)
				return
			}
			log.Print("Failed To Save Reset Password Token")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Printf("User %s Not Exist", email)
	}
}

func validateToken(w http.ResponseWriter, r *http.Request, token string) db.PasswordResetToken {
	db_connection := GetDbConnection(r)
	tokenObj := db.GetTokenUser(db_connection, token)
	if !tokenObj.IsValid() {
		setHTTPError(w, ErrInvalidToken, UserHandler, http.StatusUnauthorized)
		return db.PasswordResetToken{}
	}
	return tokenObj
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		token := getUrlVars(r, "token")
		tokenObj := validateToken(w, r, token)
		if tokenObj.IsEmpty() {
			return
		}
		templateData := setCommonTemplateData(r, w)
		templateData.Data = map[string]interface{}{"tokenObj": tokenObj}
		generateHTML(w, templateData, "layout", "reset_password")
	case http.MethodPost:
		parseForm(r, w)
		token := r.FormValue("reset_token")
		tokenObj := validateToken(w, r, token)
		if tokenObj.IsEmpty() {
			return
		}
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")
		isValidPassword, msg := checkeckPassword(password, confirmPassword)
		if !isValidPassword {
			templateData := setCommonTemplateData(r, w)
			SetErrorFlash(w, r, msg)
			templateData.Data = map[string]interface{}{"tokenObj": tokenObj}
			generateHTML(w, templateData, "layout", "reset_password")
			return

		}
		passwordHash := passwordHash(password)
		db_connection := GetDbConnection(r)
		err := db.DeleteTokenAndSetPassword(db_connection, tokenObj.Token, passwordHash, tokenObj.User.Id)
		if err != nil {
			log.Printf("Failed to update password for user %s", tokenObj.User.Email)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		SetSuccessFlash(w, r, "Password Reset Successful")
		log.Printf("Password Reset Successfully for user %s", tokenObj.User.Email)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	templateData := setCommonTemplateData(r, w)
	templateData.Data = map[string]interface{}{"loginPage": true}
	switch r.Method {
	case http.MethodGet:
		generateHTML(w, templateData, "layout", "login")
	case http.MethodPost:
		parseForm(r, w)
		email := r.FormValue("email")
		password := r.FormValue("password")
		db_connection := GetDbConnection(r)
		user := db.GetUser(db_connection, email)
		if validLogin(user, password) {
			sessionToken := generateSessionToken(r, w)
			err := db.CreateSession(db_connection, sessionToken, user.Id)
			if err != nil {
				log.Printf("Error: Creating session for user %s failed", email)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			SetSuccessFlash(w, r, "Login Successfull")
			if user.IsAdmin() {
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
		log.Print("Error: Invalid Login")
		msg := "Invalid Email or Password"
		SetErrorFlash(w, r, msg)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return

	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	db_connection := GetDbConnection(r)
	sessionData := GetSessionData(r)
	if sessionData.IsEmpty() {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	err := db.DeleteSession(db_connection, sessionData.User.Id)
	if err != nil {
		log.Printf("Error: Deleting session for user %s failed", sessionData.User.Email)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	SetSuccessFlash(w, r, "Logout Successfull")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
