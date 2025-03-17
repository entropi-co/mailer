package main

import (
	"mailer/internal"
	"mailer/internal/api"
	"mailer/internal/imap"
	"mailer/internal/instance"
	"mailer/internal/smtp"
)

func main() {
	internal.LoadConfig()

	inst := instance.CreateInstance()
	go smtp.ServeSMTP(inst)
	go imap.ServeIMAP(inst)
	api.ServeAPI(inst)
}
