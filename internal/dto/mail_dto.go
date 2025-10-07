package dto

type MailMessage struct {
	RecipientMail string `validate:"required,min=3,email"`
	RecipientName string `validate:"required,min=1"`
	Subject       string `validate:"required,min=3"`
	HTMLContent   string
	TextContent   string
}
