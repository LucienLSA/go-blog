package email

import (
	"math/rand"
	"strconv"

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

// 生成随机验证码
func GetConfirmCode() string {
	var confirmCode int
	for i := 0; i < 6; i++ {
		confirmCode = confirmCode*10 + (rand.Intn(9) + 1) //随机函数获取值
	}
	// 转换成字符串
	confirmCodeStr := strconv.Itoa(confirmCode)
	return confirmCodeStr
}
