package smtp

import (
	"github.com/mhale/smtpd"
	"github.com/sirupsen/logrus"
	"mailer/internal"
	"net"
)

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	//message, _ := mail.ReadMessage(bytes.NewReader(data))
	//fmt.Printf("From %+v\n", from)
	//fmt.Printf("To %+v\n", to)
	//handleOutbound(message)
	return nil
}

func rcptHandler(remoteAddr net.Addr, from string, to string) bool {
	//logrus.Printf("[@rcptHandler] FROM %s TO %s REMOTE %s\n", from, to, remoteAddr)
	//components := strings.Split(to, "@")
	//domain := components[1]
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

	err := server.ListenAndServe()
	if err != nil {
		logrus.Fatalf("Failed to start SMTP server: %s", err)
		return
	}
}
