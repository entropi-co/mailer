package internal

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type MailerConfig struct {
	GlobalServiceKey string `env:"GLOBAL_SERVICE_KEY" envDefault:""`

	DatabaseURL string `env:"DATABASE_URL"`

	IMAPHost string `env:"IMAP_HOST" envDefault:"0.0.0.0:993"`

	SMTPEnabled bool   `env:"SMTP_ENABLED" envDefault:"true"`
	SMTPHost    string `env:"SMTP_HOST" envDefault:"0.0.0.0:587"`
	SMTPDomain  string `env:"SMTP_DOMAIN"`
	SMTPNoAuth  bool   `env:"SMTP_NO_AUTH" envDefault:"false"`

	TLSEnabled         bool   `env:"TLS_ENABLED" envDefault:"false"`
	TLSCertificatePath string `env:"TLS_CERTIFICATE_PATH" envDefault:"./mailer.crt"`
	TLSKeyPath         string `env:"TLS_KEY_PATH" envDefault:"./mailer.key"`
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
		SMTPEnabled:      true,
		SMTPHost:         "localhost:5789",
		SMTPDomain:       "test.dev",
		SMTPNoAuth:       true,
		GlobalServiceKey: "",
	}
}
