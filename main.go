package main

import (
	"mailer/internal"
	"mailer/internal/smtp"
)

func main() {
	internal.LoadConfig()
	smtp.ServeSMTP()
}
