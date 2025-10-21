package service

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChrisChantszto/futurefrontier/internal/config"
)

type SMTPService struct {
	config config.Config
}

func NewSMTPService(cfg config.Config) *SMTPService {
	return &SMTPService{config: cfg}
}

// GetSMTPConfig returns the appropriate SMTP configuration
// Uses custom SMTP if configured, otherwise falls back to company SMTP
func (s *SMTPService) GetSMTPConfig() config.SMTPConfig {
	// Check if custom SMTP is configured and allowed
	if s.config.AllowCustomSMTP && s.config.CustomSMTP.Host != "" {
		return s.config.CustomSMTP
	}
	// Fall back to company SMTP
	return s.config.CompanySMTP
}

// SendOTPEmail sends an OTP code via email with the specified locale
func (s *SMTPService) SendOTPEmail(to, code, locale string) error {
	smtpConfig := s.GetSMTPConfig()
	
	if smtpConfig.Host == "" || smtpConfig.Username == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	// Normalize locale
	locale = s.normalizeLocale(locale)

	// Get email subject based on locale
	subject := s.getSubjectByLocale("Your One-Time Login Code", locale)

	// Prepare email data
	data := map[string]interface{}{
		"OTP":        code,
		"TTLMinutes": s.config.OTPTTLMinutes,
	}

	// Get email body from template
	body, err := s.renderTemplate("otp", locale, data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Prepare message
	message := s.buildEmailMessage(smtpConfig.FromEmail, smtpConfig.FromName, to, subject, body, smtpConfig.ReplyTo, smtpConfig.BCC)

	// Send email
	return s.sendEmail(smtpConfig, to, message)
}

// buildEmailMessage constructs the email message with headers
func (s *SMTPService) buildEmailMessage(fromEmail, fromName, to, subject, body, replyTo, bcc string) string {
	var msg strings.Builder
	
	// From header
	if fromName != "" {
		msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, fromEmail))
	} else {
		msg.WriteString(fmt.Sprintf("From: %s\r\n", fromEmail))
	}
	
	// To header
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	
	// Reply-To header
	if replyTo != "" {
		msg.WriteString(fmt.Sprintf("Reply-To: %s\r\n", replyTo))
	}
	
	// BCC header (for internal tracking)
	if bcc != "" {
		msg.WriteString(fmt.Sprintf("Bcc: %s\r\n", bcc))
	}
	
	// Subject and content type
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)
	
	return msg.String()
}

// sendEmail sends the email using SMTP
func (s *SMTPService) sendEmail(smtpConfig config.SMTPConfig, to, message string) error {
	// SMTP server address
	addr := fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port)
	
	// Recipients (include BCC if specified)
	recipients := []string{to}
	if smtpConfig.BCC != "" {
		recipients = append(recipients, smtpConfig.BCC)
	}
	
	// Authentication
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)
	
	// For port 465 (SMTPS), we need to use TLS from the start
	if smtpConfig.Port == 465 && smtpConfig.UseTLS {
		return s.sendEmailTLS(addr, auth, smtpConfig.FromEmail, recipients, message)
	}
	
	// For other ports, use regular SMTP (with STARTTLS if needed)
	return smtp.SendMail(addr, auth, smtpConfig.FromEmail, recipients, []byte(message))
}

// sendEmailTLS sends email using TLS connection (for port 465)
func (s *SMTPService) sendEmailTLS(addr string, auth smtp.Auth, from string, to []string, message string) error {
	// Create TLS connection
	tlsConfig := &tls.Config{
		ServerName: strings.Split(addr, ":")[0],
	}
	
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()
	
	// Create SMTP client
	client, err := smtp.NewClient(conn, tlsConfig.ServerName)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()
	
	// Authenticate
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}
	
	// Set sender
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	
	// Set recipients
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}
	
	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}
	
	_, err = writer.Write([]byte(message))
	if err != nil {
		writer.Close()
		return fmt.Errorf("failed to write message: %w", err)
	}
	
	return writer.Close()
}

// normalizeLocale normalizes locale strings to match directory structure
func (s *SMTPService) normalizeLocale(locale string) string {
	// If locale is empty or invalid, default to English
	if locale == "" {
		return s.config.I18n.DefaultLocale
	}
	
	// Handle common variations
	switch strings.ToLower(locale) {
	case "zh", "zh-cn", "zh_cn":
		return "zh-hans"
	case "zh-tw", "zh_tw":
		return "zh-hant"
	default:
		// Check if it's one of our supported locales
		for _, supported := range s.config.I18n.SupportedLocales {
			if strings.EqualFold(locale, supported) {
				return supported
			}
		}
		// Default to the configured default locale
		return s.config.I18n.DefaultLocale
	}
}

// getSubjectByLocale returns the appropriate email subject based on locale
func (s *SMTPService) getSubjectByLocale(defaultSubject, locale string) string {
	subjects := map[string]map[string]string{
		"Your One-Time Login Code": {
			"en":      "Your One-Time Login Code",
			"zh-hans": "您的验证码",
			"zh-hant": "您的驗證碼",
		},
	}
	
	// Look up the subject in our map
	if localeMap, exists := subjects[defaultSubject]; exists {
		if subject, exists := localeMap[locale]; exists {
			return subject
		}
	}
	
	// Fallback to the default subject
	return defaultSubject
}

// renderTemplate renders an email template with the provided data
func (s *SMTPService) renderTemplate(templateName, locale string, data map[string]interface{}) (string, error) {
	// Get template path
	templatePath := filepath.Join(s.config.I18n.TemplatesDir, "emails", locale, templateName+".html")
	
	// Check if template exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		// Fallback to default locale
		templatePath = filepath.Join(s.config.I18n.TemplatesDir, "emails", s.config.I18n.DefaultLocale, templateName+".html")
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			return "", fmt.Errorf("template not found: %s", templateName)
		}
	}
	
	// Parse template
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	
	// Execute template
	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	
	return buf.String(), nil
}
