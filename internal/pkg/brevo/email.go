package brevo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type BrevoService interface {
	SendPasswordResetEmail(toEmail, toName, resetLink string) error
}

type brevoService struct {
	apiKey      string
	senderEmail string
	senderName  string
	client      *http.Client
}

func NewBrevoService(apiKey, senderEmail, senderName string) BrevoService {
	return &brevoService{
		apiKey:      apiKey,
		senderEmail: senderEmail,
		senderName:  senderName,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *brevoService) SendPasswordResetEmail(toEmail, toName, resetLink string) error {
	if s.apiKey == "" {
		return errors.New("brevo api key is not configured")
	}

	htmlContent := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 5px;">
			<h2 style="color: #333;">Reset Password Sembako App</h2>
			<p>Halo <b>%s</b>,</p>
			<p>Kami menerima permintaan untuk mereset kata sandi akun Anda. Klik tombol di bawah ini untuk melanjutkan:</p>
			<div style="text-align: center; margin: 30px 0;">
				<a href="%s" style="background-color: #4CAF50; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; font-weight: bold; display: inline-block;">Reset Password</a>
			</div>
			<p>Atau salin tautan berikut ke browser Anda:</p>
			<p style="word-break: break-all; color: #555;"><a href="%s">%s</a></p>
			<p style="color: #777; font-size: 12px; margin-top: 30px;">Tautan ini akan kedaluwarsa dalam 30 menit. Jika Anda tidak merasa meminta reset password, abaikan email ini.</p>
		</div>
	`, toName, resetLink, resetLink, resetLink)

	payload := BrevoPayload{
		Sender: BrevoEmailSender{
			Email: s.senderEmail,
			Name:  s.senderName,
		},
		To: []BrevoEmailRecipient{
			{
				Email: toEmail,
				Name:  toName,
			},
		},
		Subject:     "Reset Password - Sembako App",
		HtmlContent: htmlContent,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", s.apiKey)
	req.Header.Set("content-type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send email via brevo, status code: %d", resp.StatusCode)
	}

	return nil
}
