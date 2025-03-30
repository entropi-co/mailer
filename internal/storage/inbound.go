package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/godruoyi/go-snowflake"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	InboundsTable          = "inbounds"
	InboundsColID          = "id"
	InboundsColBody        = "body"
	InboundsColSender      = "sender"
	InboundsColDeliveredAt = "delivered_at"

	InboundsMailboxesTable      = "inbounds_mailboxes"
	InboundsMailboxesColInbound = "inbound"
	InboundsMailboxesColMailbox = "mailbox"
	InboundsMailboxesColUID     = "uid"
)

const QueryInsertInboundToPrimaryMailbox = `
WITH updated AS (
    UPDATE mailboxes AS m
        SET uid_next = m.uid_next + 1
        WHERE id IN (SELECT id
                     FROM (SELECT m.id,
                                  ROW_NUMBER() OVER (PARTITION BY owner ORDER BY priority, m.created_at) AS _row
                           FROM mailboxes m
                                    LEFT JOIN public.users
                                              on users.id = m.owner
                           WHERE users.local IN $1) as _sub
                     WHERE _row = 1)
        RETURNING m.id, uid_next)
INSERT
INTO inbounds_mailboxes (inbound, mailbox, uid)
SELECT $2, updated.id, updated.uid_next
FROM updated;
`

type InboundHeader map[string]interface{}

type Inbound struct {
	ID          uint64    `json:"id"`
	Body        []byte    `json:"body"`
	Sender      string    `json:"sender"`
	DeliveredAt time.Time `json:"delivered_at"`
}

func (i *InboundHeader) Scan(src any) error {
	logrus.Printf("[InboundHeader#Scan] Begin")
	data, ok := src.([]uint8)
	if !ok {
		return errors.New("source must be []uint8")
	}

	var header InboundHeader
	err := json.Unmarshal(data, &header)
	if err != nil {
		return err
	}

	logrus.Printf("[InboundHeader#Scan] Decoded header: %#v", header)

	*i = header
	return nil
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

func NewInbound(body []byte, sender string, deliveredAt time.Time) (*Inbound, error) {
	return &Inbound{
		ID:          snowflake.ID(),
		Body:        body,
		Sender:      sender,
		DeliveredAt: deliveredAt,
	}, nil
}

// AddInboundToRecipientsPrimaryMailbox adds given inbound to target recipients primary mailboxes
// The inbound is inserted in same transaction, rolling back if either insert inbound or insert relation has failed
func (s *Storage) AddInboundToRecipientsPrimaryMailbox(inbound *Inbound, recipients []string) error {
	tx, err := s.Database.Begin()
	if err != nil {
		return err
	}

	// Insert to inbounds
	_, err = squirrel.
		Insert("inbounds").
		Columns(
			InboundsColID,
			InboundsColBody,
			InboundsColSender,
			InboundsColDeliveredAt,
		).
		Values(inbound.ID, inbound.Body, inbound.Sender, inbound.DeliveredAt).
		RunWith(tx).
		Exec()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	// Insert to inbounds_mailboxes
	_, err = s.Database.Query(QueryInsertInboundToPrimaryMailbox, recipients, inbound.ID)
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
		Select(
			InboundsColID,
			InboundsColBody,
			InboundsColSender,
			InboundsColDeliveredAt,
			"sequence",
			"uid",
		).
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
		Select(
			InboundsColID,
			InboundsColBody,
			InboundsColSender,
			InboundsColDeliveredAt,
			"uid",
		).
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
