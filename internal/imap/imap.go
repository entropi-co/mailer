package imap

import (
	"github.com/emersion/go-imap/v2/imapserver"
	"github.com/sirupsen/logrus"
	"mailer/internal/instance"
	"net"
)

type IMAP struct {
	Instance  *instance.Instance
	Mailboxes map[uint64]*Mailbox
}

func ServeIMAP(instance *instance.Instance) {
	logrus.Infoln("Initialize IMAP")

	imap := &IMAP{Instance: instance}

	server := imapserver.New(&imapserver.Options{
		NewSession: func(conn *imapserver.Conn) (imapserver.Session, *imapserver.GreetingData, error) {
			return imap.NewSession(), nil, nil
		},
		Caps:         nil,
		Logger:       nil,
		TLSConfig:    nil,
		InsecureAuth: false,
		DebugWriter:  nil,
	})

	ln, err := net.Listen("tcp", ":993")
	if err != nil {
		logrus.Fatalf("Failed to create IMAP listener: %+v", err)
	}

	logrus.Infoln("Serving IMAP")
	if err := server.Serve(ln); err != nil {
		logrus.Fatalf("Failed to start IMAP server: %+v", err)
	}
}
