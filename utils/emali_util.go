package utils

import (
	"fmt"
	"github.com/jordan-wright/email"
	"gopkg.in/yaml.v3"
	"net/smtp"
	"os"
	"wxbot4g/logger"
)

func init() {
	file, err := os.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println("fail to read file:", err)
	}
	err = yaml.Unmarshal(file, &config)
	qqInfo := config.Email
	if qqInfo.MailRec == "" || qqInfo.Code == "" || qqInfo.MailSend == "" {
		logger.Log.Info("QQ邮箱配置异常，无法使用通知功能")
	} else {
		logger.Log.Info(fmt.Sprintf("通知发送邮箱%s, 通知接收邮箱%s", qqInfo.MailSend, qqInfo.MailRec))
	}
	if err != nil {
		fmt.Println("error init config,", err.Error())
	}
}

var config = &notifyConfig{}
var qqInfo = &emailQQ{}

type notifyConfig struct {
	Email emailQQ `yaml:"notify"`
}
type emailQQ struct {
	MailSend string `yaml:"mail_send"`
	MailRec  string `yaml:"mail_rec"`
	Code     string
}

func Notify(msg string, subTitle string) {
	qqInfo := config.Email
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "wxbot4g通知 <" + qqInfo.MailSend + ">"
	// 设置接收方的邮箱
	e.To = []string{qqInfo.MailRec}
	//设置主题
	e.Subject = subTitle
	//设置文件发送的内容
	e.Text = []byte(msg)
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", qqInfo.MailSend, qqInfo.Code, "smtp.qq.com"))
	if err != nil {
		logger.Log.Fatal(err)
	}
}
