package email

import (
	"fmt"
	"log"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// todo move this to env
const FromName string = "BloTilsMailer"
const FromMail string = "noreply@blotils.com"

func sendMail(toMail string, subject string, htmlContent string) error {
	to := mail.NewEmail(toMail, toMail)
	from := mail.NewEmail(FromName, FromMail)
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient("")
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: %d", response.StatusCode)
	}
	return nil
}

func resetPasswordMail(email string, token string) {
	resetLink := fmt.Sprintf("http://blotils.com:8000/reset_password/%s", token)
	subject := "Blotils Password Rest Link"
	htmlContent := fmt.Sprintf(`<h3>Rest Your Password <a href=%s>Here</a></h3>`, resetLink)
	err := sendMail(email, subject, htmlContent)
	if err != nil {
		log.Printf("Failed to send reset password mail to %s, err:%v", email, err)
	}
	log.Printf("Reset password mail sent to %s", email)
}
