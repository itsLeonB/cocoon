package dto

type MailMessage struct {
	Recipient string `validate:"required,min=3,email"`
	Subject   string `validate:"required,min=3"`
	BodyType  string `validate:"required,min=3,oneof=text/plain"`
	Body      string `validate:"required,min=3"`
}
