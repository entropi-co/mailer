package smtp

import (
	"net/mail"
)

const (
	MessageInbound = 1 << iota
	MessageOutbound
)

// handleOutbound delivers message to other smtp server
func handleOutbound(message *mail.Message, from string, to []string, domain string) {

}
