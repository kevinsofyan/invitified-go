package utils

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendHTMLEmail(to string, subject string, htmlBody string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Debugging print statements
	fmt.Println("From address: ", from)
	fmt.Println("To address: ", to)
	fmt.Println("Password: ", password)
	fmt.Println("SMTP Host: ", smtpHost)
	fmt.Println("SMTP Port: ", smtpPort)

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("missing required environment variables")
	}

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(smtpHost, port, from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func GetOrderConfirmationEmail(orderNumber string, amount string) string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    <tr>
                        <td style="padding: 40px;">
                            <h1 style="color: #333333; margin-bottom: 30px; text-align: center;">Payment Confirmed! 🎉</h1>
                            
                            <p style="color: #666666; font-size: 16px; line-height: 24px; margin-bottom: 20px;">
                                Thank you for your order with Invitified! We're pleased to confirm that your payment has been successfully processed.
                            </p>

                            <div style="background-color: #f8f9fa; border-radius: 6px; padding: 20px; margin: 30px 0;">
                                <p style="margin: 0; color: #333333; font-size: 16px;">
                                    <strong>Order Number:</strong> ` + orderNumber + `<br>
                                    <strong>Amount Paid:</strong> ` + amount + `
                                </p>
                            </div>

                            <p style="color: #666666; font-size: 16px; line-height: 24px;">
                                You can track your order status by logging into your Invitified account. If you have any questions, please don't hesitate to contact our support team.
                            </p>

                            <div style="text-align: center; margin-top: 40px;">
                                <a href="https://invitified.com/orders" style="background-color: #4CAF50; color: #ffffff; padding: 12px 30px; text-decoration: none; border-radius: 4px; font-weight: bold;">View Order Details</a>
                            </div>
                        </td>
                    </tr>
                    <tr>
                        <td style="background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 0 0 8px 8px;">
                            <p style="color: #999999; font-size: 14px; margin: 0;">
                                This is an automated message, please do not reply directly to this email.<br>
                                © 2024 Invitified. All rights reserved.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`
}
