package smtp

import (
	"github.com/go-gomail/gomail"
	"mailer/internal"
	"testing"
)

func TestServeSMTP(t *testing.T) {
	internal.LoadTestConfig()
	go ServeSMTP()

	dialer := gomail.NewDialer("smtp.gmail.com", 465, "", "")

	m := gomail.NewMessage()
	m.SetAddressHeader("From", "kappa@entropi.kr", "Kappa")
	m.SetAddressHeader("To", "arcranion@gmail.com", "Arcranion")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	err := dialer.DialAndSend(m)
	if err != nil {
		t.Fatalf("failed to dial and send: %+v", err)
	}
}
