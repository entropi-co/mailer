package storage

import (
	"github.com/Masterminds/squirrel"
	"io"
	"net/mail"
	"strings"
	"time"
)

type Inbound struct {
	ID          uint64              `json:"id"`
	Header      map[string][]string `json:"header"`
	Body        string              `json:"body"`
	Sender      string              `json:"sender"`
	DeliveredAt time.Time           `json:"delivered_at"`
}

func NewInbound(message *mail.Message, sender string, deliveredAt time.Time) (*Inbound, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, message.Body)
	if err != nil {
		return nil, err
	}
	return &Inbound{
		Header:      message.Header,
		Body:        buf.String(),
		Sender:      sender,
		DeliveredAt: deliveredAt,
	}, nil
}

func (s *Storage) CreateInbound(inbound Inbound, recipients []uint64) error {
	tx, err := s.Database.Begin()
	if err != nil {
		return err
	}

	// Insert inbound
	_, err = squirrel.
		Insert("inbounds").
		Columns("id", "content", "sender", "delivered_at").
		Values(inbound.ID, inbound.Body, inbound.Sender, inbound.DeliveredAt).
		RunWith(tx).
		Exec()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	// Insert inbound to recipients
	builder := squirrel.
		Insert("inbounds_recipients").
		Columns("inbound", "recipient").
		Values(inbound, recipients).
		RunWith(tx)
	for i := range recipients {
		recipient := recipients[i]
		builder.Values(inbound.ID, recipient)
	}
	_, err = builder.Exec()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return err
}
