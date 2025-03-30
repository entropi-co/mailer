package main

import (
	"github.com/sirupsen/logrus"
	"mailer/internal"
	"mailer/internal/api"
	"mailer/internal/imap"
	"mailer/internal/instance"
	"mailer/internal/smtp"
)

func main() {
	internal.LoadConfig()

	logrus.SetLevel(logrus.DebugLevel)

	inst := instance.CreateInstance()
	go smtp.ServeSMTP(inst)
	go imap.ServeIMAP(inst)
	api.ServeAPI(inst)
}
