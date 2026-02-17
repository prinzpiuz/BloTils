package mailer

import (
	"fmt"
	"log"
)

// SendResetPasswordEmail sends a password reset email
func SendResetPasswordEmail(toEmail, token string) error {
	resetLink := fmt.Sprintf("%s/reset_password/%s", GetBaseURL(), token)

	email := Email{
		To:      toEmail,
		Subject: "BloTils Password Reset",
		HTMLContent: fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2>Reset Your Password</h2>
				<p>Click the button below to reset your password:</p>
				<p style="margin: 30px 0;">
					<a href="%s" style="background-color: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px;">
						Reset Password
					</a>
				</p>
				<p style="color: #666; font-size: 14px;">
					If you didn't request this, please ignore this email.<br>
					This link will expire in 1 hour.
				</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
				<p style="color: #999; font-size: 12px;">BloTils - Utilities for your static blog</p>
			</div>
		`, resetLink),
	}

	if err := Send(email); err != nil {
		log.Printf("Failed to send reset password email to %s: %v", toEmail, err)
		return err
	}

	log.Printf("Reset password email sent to %s", toEmail)
	return nil
}

// SendWelcomeEmail sends a welcome email to new users
func SendWelcomeEmail(toEmail, name string) error {
	email := Email{
		To:      toEmail,
		Subject: "Welcome to BloTils!",
		HTMLContent: fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2>Welcome, %s!</h2>
				<p>Thanks for joining BloTils. We're excited to have you.</p>
				<p>Get started by adding your first domain.</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
				<p style="color: #999; font-size: 12px;">BloTils - Utilities for your static blog</p>
			</div>
		`, name),
	}

	return Send(email)
}

// SendAccountApprovedEmail notifies user their account was approved
func SendAccountApprovedEmail(toEmail, name string) error {
	loginLink := fmt.Sprintf("%s/login", GetBaseURL())

	email := Email{
		To:      toEmail,
		Subject: "Your BloTils Account Has Been Approved!",
		HTMLContent: fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2>Good news, %s!</h2>
				<p>Your BloTils account has been approved. You can now log in and start using the platform.</p>
				<p style="margin: 30px 0;">
					<a href="%s" style="background-color: #28a745; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px;">
						Log In Now
					</a>
				</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
				<p style="color: #999; font-size: 12px;">BloTils - Utilities for your static blog</p>
			</div>
		`, name, loginLink),
	}

	return Send(email)
}
