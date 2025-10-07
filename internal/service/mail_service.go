package service

import (
	"context"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/rotisserie/eris"
)

type brevoMailService struct {
	client     *brevo.APIClient
	senderMail string
	senderName string
}

func NewMailService(mailConfig config.Mail) MailService {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", mailConfig.ApiKey)
	br := brevo.NewAPIClient(cfg)
	return &brevoMailService{br, mailConfig.SenderMail, mailConfig.SenderName}
}

func (ms *brevoMailService) Send(ctx context.Context, msg dto.MailMessage) error {
	mail := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  ms.senderName,
			Email: ms.senderMail,
		},
		To: []brevo.SendSmtpEmailTo{{
			Email: msg.RecipientMail,
			Name:  msg.RecipientName,
		}},
		Subject:     msg.Subject,
		HtmlContent: msg.HTMLContent,
		TextContent: msg.TextContent,
	}

	if _, _, err := ms.client.TransactionalEmailsApi.SendTransacEmail(ctx, mail); err != nil {
		return eris.Wrap(err, "error sending email")
	}
	return nil
}
