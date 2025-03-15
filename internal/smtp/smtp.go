package smtp

import (
	"bytes"
	"fmt"
	"github.com/mhale/smtpd"
	"github.com/sirupsen/logrus"
	"mailer/internal"
	"net"
	"net/mail"
	"strings"
)

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	_, _ = mail.ReadMessage(bytes.NewReader(data))
	fmt.Printf("From %+v\n", from)
	fmt.Printf("To %+v\n", to)
	//handleOutbound(message)
	return nil
}

func rcptHandler(remoteAddr net.Addr, from string, to string) bool {
	logrus.Printf("[@rcptHandler] FROM %s TO %s REMOTE %s\n", from, to, remoteAddr)
	components := strings.Split(to, "@")
	domain := components[1]
	fmt.Printf("[@rcptHandler] Domain %s \n", domain)
	return true
}

// authHandler checks if password matches service key
func authHandler(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
	return string(username) == "service" && string(password) == internal.Config.SMTPServiceKey, nil
}

func ServeSMTP() {
	if !internal.Config.SMTPEnabled {
		return
	}

	logrus.Infoln("Initialize SMTP")

	server := &smtpd.Server{
		Addr:         internal.Config.SMTPHost,
		Handler:      mailHandler,
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

	logrus.Infoln("Listening")
	err := server.ListenAndServe()
	if err != nil {
		logrus.Fatalf("Failed to start SMTP server: %s", err)
		return
	}
}
