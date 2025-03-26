package storage

import (
	"fmt"
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

type InboundMetadata struct {
	// Not present on the actual table
	// Calculated with row_number() in query
	Sequence uint32 `json:"sequence"`

	// Located in inbounds_mailboxes table
	UID uint64 `json:"uid"`
}

type InboundWithMetadata struct {
	Inbound
	InboundMetadata
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

func (s *Storage) CreateInbound(inbound *Inbound, recipients []uint64) error {
	tx, err := s.Database.Begin()
	if err != nil {
		return err
	}

	// Insert inbound
	_, err = squirrel.
		Insert("inbounds").
		Columns("id", "header", "content", "sender", "delivered_at").
		Values(inbound.ID, inbound.Header, inbound.Body, inbound.Sender, inbound.DeliveredAt).
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

func (s *Storage) QueryInboundsBySequences(mailbox uint64, sequences []uint32) ([]*InboundWithMetadata, error) {
	withQuery := fmt.Sprintf(
		"with _ranked as "+
			"(select inbounds.*, inbounds_mailboxes.uid as uid, row_number() over (order by id) as sequence"+
			" from inbounds"+
			" left join inbounds_mailboxes on inbounds.id = inbounds_mailboxes.inbound"+
			" where mailbox = %d"+
			")",
		mailbox,
	)

	rows, err := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("id", "header", "content", "sender", "delivered_at", "sequence", "uid").
		From("_ranked").
		Prefix(withQuery).
		Where(squirrel.Eq{"sequence": sequences}).
		RunWith(s.Database).
		Query()
	if err != nil {
		return nil, err
	}

	inbounds := make([]*InboundWithMetadata, len(sequences))
	i := 0
	for rows.Next() {
		inbound := InboundWithMetadata{}
		if err := rows.Scan(
			&inbound.ID,
			&inbound.Header,
			&inbound.Body,
			&inbound.Sender,
			&inbound.DeliveredAt,
			&inbound.Sequence,
			&inbound.UID,
		); err != nil {
			return nil, err
		}

		inbounds[i] = &inbound

		i++
	}

	return inbounds, nil
}

func (s *Storage) QueryInboundsByUIDS(mailbox uint64, uids []uint32) ([]*InboundWithMetadata, error) {
	rows, err := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("id", "header", "content", "sender", "delivered_at", "uid").
		From("_ranked").
		Prefix("with _ranked as "+
			"(select inbounds.*, inbounds_mailboxes.uid as uid, row_number() over (order by id) as sequence"+
			" from inbounds"+
			" left join inbounds_mailboxes on inbounds.id = inbounds_mailboxes.inbound"+
			" where mailbox = $1"+
			")",
			mailbox,
		).
		Where(squirrel.Eq{"uid": uids}).
		RunWith(s.Database).
		Query()
	if err != nil {
		return nil, err
	}

	inbounds := make([]*InboundWithMetadata, len(uids))
	i := 0
	for rows.Next() {
		if err := rows.Scan(
			&inbounds[i].ID,
			&inbounds[i].Header,
			&inbounds[i].Body,
			&inbounds[i].Sender,
			&inbounds[i].DeliveredAt,
			&inbounds[i].Sequence,
			&inbounds[i].UID,
		); err != nil {
			return nil, err
		}

		i++
	}

	return inbounds, nil
}
