package email

import (
	"github.com/LucienLSA/go-blog/settings"
	"go.uber.org/zap"
	"gopkg.in/mail.v2"
)

// Send 发送邮件
func Send(mailTex, emailTo, subject, emailFrom string) error {
	m := mail.NewMessage()
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", mailTex)
	d := mail.NewDialer(settings.Conf.EmailConfig.SmtpHost, 465, settings.Conf.EmailConfig.SmtpEmail, settings.Conf.EmailConfig.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		zap.L().Error("DialAndSend failed", zap.Error(err))
		return err
	}
	return nil
}
