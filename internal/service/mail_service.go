package service

import (
	"crypto/tls"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/rotisserie/eris"
	"gopkg.in/gomail.v2"
)

type MailService interface {
	Send(msg dto.MailMessage) error
}

type mailService struct {
	dialer *gomail.Dialer
	sender string
}

func NewMailService(cfg config.Mail) MailService {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &mailService{dialer, cfg.Sender}
}

func (ms *mailService) Send(msg dto.MailMessage) error {
	message := gomail.NewMessage()
	message.SetHeader("From", ms.sender)
	message.SetHeader("To", msg.Recipient)
	message.SetHeader("Subject", msg.Subject)
	message.SetBody(msg.BodyType, msg.Body)

	if err := ms.dialer.DialAndSend(message); err != nil {
		return eris.Wrap(err, "error sending mail")
	}
	return nil
}
