package internal

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type MailerConfig struct {
	SMTPEnabled        bool     `env:"SMTP_ENABLED" envDefault:"true"`
	SMTPHost           string   `env:"SMTP_HOST" envDefault:"0.0.0.0:587"`
	SMTPDomains        []string `env:"SMTP_DOMAINS"`
	SMTPNoAuth         bool     `env:"SMTP_NO_AUTH" envDefault:"false"`
	SMTPServiceKey     string   `env:"SMTP_SERVICE_KEY" envDefault:""`
	TLSEnabled         bool     `env:"TLS_ENABLED" envDefault:"false"`
	TLSCertificatePath string   `env:"TLS_CERTIFICATE_PATH" envDefault:"./mailer.crt"`
	TLSKeyPath         string   `env:"TLS_KEY_PATH" envDefault:"./mailer.key"`
}

var Config MailerConfig

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		logrus.Infoln("Error loading .env file, skipping")
	}

	err = env.Parse(&Config)
	if err != nil {
		logrus.Fatalf("Error parsing environment variables: %+v\n", err)
	}
}

func LoadTestConfig() {
	Config = MailerConfig{
		SMTPEnabled:    true,
		SMTPHost:       "localhost:5789",
		SMTPDomains:    []string{"test.test"},
		SMTPNoAuth:     true,
		SMTPServiceKey: "",
	}
}
