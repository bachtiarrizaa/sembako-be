package brevo

type BrevoEmailSender struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type BrevoEmailRecipient struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type BrevoPayload struct {
	Sender      BrevoEmailSender      `json:"sender"`
	To          []BrevoEmailRecipient `json:"to"`
	Subject     string                `json:"subject"`
	HtmlContent string                `json:"htmlContent"`
}
