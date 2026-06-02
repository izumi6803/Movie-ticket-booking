package services

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort:     getEnv("SMTP_PORT", "587"),
		smtpUsername: getEnv("SMTP_USERNAME", ""),
		smtpPassword: getEnv("SMTP_PASSWORD", ""),
		fromEmail:    getEnv("FROM_EMAIL", "cinema@example.com"),
	}
}

func (s *EmailService) IsConfigured() bool {
	return s.smtpUsername != "" && s.smtpPassword != ""
}

// SendTicketConfirmation gửi email xác nhận vé
func (s *EmailService) SendTicketConfirmation(toEmail, userName string, bookingDetails map[string]interface{}) error {
	if !s.IsConfigured() {
		fmt.Println("Email service not configured. Skipping email send.")
		fmt.Printf("Would send ticket confirmation to %s for booking %v\n", toEmail, bookingDetails)
		return nil
	}

	subject := "Your Cinema Ticket Booking Confirmation"

	// Tạo nội dung email HTML
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #e50019; color: white; padding: 20px; text-align: center; }
        .content { background: #f9f9f9; padding: 20px; margin: 20px 0; }
        .footer { text-align: center; color: #666; font-size: 12px; margin-top: 20px; }
        .ticket-info { background: white; padding: 15px; margin: 10px 0; border-left: 4px solid #e50019; }
        .qr-code { text-align: center; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Booking Confirmation</h1>
        </div>
        
        <div class="content">
            <p>Dear %s,</p>
            <p>Thank you for your booking! Your payment has been confirmed.</p>
            
            <div class="ticket-info">
                <h3>Ticket Details</h3>
                <p><strong>Booking Code:</strong> %v</p>
                <p><strong>Movie:</strong> %v</p>
                <p><strong>Showtime:</strong> %v</p>
                <p><strong>Seats:</strong> %v</p>
                <p><strong>Total Amount:</strong> $%v</p>
            </div>
            
            <div class="qr-code">
                <p><strong>Your QR Code for Entry:</strong></p>
                <p style="font-family: monospace; background: #f0f0f0; padding: 10px; word-break: break-all;">
                    %v
                </p>
            </div>
            
            <p>Please arrive at least 15 minutes before the showtime.</p>
            <p>Show this QR code at the entrance for verification.</p>
        </div>
        
        <div class="footer">
            <p>This is an automated email. Please do not reply.</p>
            <p>© 2024 Cinema Booking System</p>
        </div>
    </div>
</body>
</html>`,
		userName,
		bookingDetails["bookingCode"],
		bookingDetails["movie"],
		bookingDetails["showtime"],
		bookingDetails["seats"],
		bookingDetails["amount"],
		bookingDetails["qrCode"],
	)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", toEmail, subject, body))

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)

	return smtp.SendMail(addr, auth, s.fromEmail, []string{toEmail}, msg)
}

// SendPaymentFailed gửi email thông báo thanh toán thất bại
func (s *EmailService) SendPaymentFailed(toEmail, userName string, bookingDetails map[string]interface{}) error {
	if !s.IsConfigured() {
		fmt.Println("Email service not configured. Skipping email send.")
		return nil
	}

	subject := "Payment Failed - Cinema Booking"
	body := fmt.Sprintf(`
<html>
<body>
    <h2>Payment Failed</h2>
    <p>Dear %s,</p>
    <p>We regret to inform you that your payment for booking %v could not be processed.</p>
    <p>Please try again or contact support if you need assistance.</p>
</body>
</html>`, userName, bookingDetails["bookingCode"])

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", toEmail, subject, body))

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)

	return smtp.SendMail(addr, auth, s.fromEmail, []string{toEmail}, msg)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
