package email

import (
	"github.com/LucienLSA/go-blog/settings"
	"gopkg.in/mail.v2"
)

type EmailSender struct {
	SmtpHost      string `json:"smtpHost"`
	SmtpEmailFrom string `json:"smtpEmail"`
	SmtpPass      string `json:"smtpPass"`
}

func NewEmailSender() *EmailSender {
	eConfig := settings.Conf.EmailConfig
	return &EmailSender{
		SmtpHost:      eConfig.SmtpHost,
		SmtpEmailFrom: eConfig.SmtpEmail,
		SmtpPass:      eConfig.SmtpPass,
	}
}

// Send 发送邮件
func (s *EmailSender) Send(data, emailTo, subject string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.SmtpEmailFrom)
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", data)
	d := mail.NewDialer(s.SmtpHost, 465, s.SmtpEmailFrom, s.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
