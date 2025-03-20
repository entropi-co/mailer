package smtp

import (
	"bytes"
	"github.com/mhale/smtpd"
	"github.com/sirupsen/logrus"
	"mailer/internal"
	"mailer/internal/instance"
	"mailer/internal/storage"
	"net"
	"net/mail"
	"strings"
	"time"
)

type SMTP struct {
	Instance *instance.Instance
}

func (s *SMTP) mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	message, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return err
	}
	inbound, err := storage.NewInbound(message, from, time.Now())
	if err != nil {
		return err
	}

	recipientIds, err := s.Instance.Storage.QueryUserIDsByLocals(to)
	if err != nil {
		return err
	}

	if err = s.Instance.Storage.CreateInbound(inbound, recipientIds); err != nil {
		return err
	}

	return nil
}

// TODO: Add spam remote determination
func rcptHandler(remoteAddr net.Addr, from string, to string) bool {
	components := strings.Split(to, "@")
	domain := components[1]
	if domain != internal.Config.SMTPDomain {
		return false
	}
	return true
}

// authHandler checks if password matches service key
// TODO: Add API key access method
func authHandler(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
	return string(username) == "service" && string(password) == internal.Config.GlobalServiceKey, nil
}

func ServeSMTP(instance *instance.Instance) {
	if !internal.Config.SMTPEnabled {
		return
	}

	logrus.Infoln("Initialize SMTP")

	smtp := &SMTP{
		Instance: instance,
	}

	server := &smtpd.Server{
		Addr:         internal.Config.SMTPHost,
		Handler:      smtp.mailHandler,
		HandlerRcpt:  rcptHandler,
		AuthHandler:  authHandler,
		AuthRequired: !internal.Config.SMTPNoAuth,
	}

	if internal.Config.TLSEnabled {
		logrus.Infoln("Configuring TLS")
		if err := server.ConfigureTLS(internal.Config.TLSCertificatePath, internal.Config.TLSKeyPath); err != nil {
			logrus.Fatalf("Unable to configure TLS for SMTP: %v", err)
			return
		}
	}

	logrus.Infoln("Listening SMTP")
	err := server.ListenAndServe()
	if err != nil {
		logrus.Fatalf("Failed to start SMTP server: %s", err)
		return
	}
}
