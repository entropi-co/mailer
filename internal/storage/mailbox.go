package storage

import (
	"github.com/Masterminds/squirrel"
	"time"
)

type MailboxType string

const (
	MailboxInbox MailboxType = "inbox"
	MailboxJunk  MailboxType = "junk"
	MailboxSent  MailboxType = "sent"
	MailboxDraft MailboxType = "draft"
	MailboxTrash MailboxType = "trash"
	MailboxUser  MailboxType = "user"
)

type Mailbox struct {
	ID          uint64      `json:"id"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Owner       uint64      `json:"owner"`
	Type        MailboxType `json:"type"`
	CreatedAt   time.Time   `json:"created_at"`
}

type MailboxMetadata struct {
	UID          uint64 `json:"uid"`
	InboundCount int64  `json:"inbound_count"`
}

func (s *Storage) QueryMailbox(user uint64, name string) (*Mailbox, error) {
	row := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("id", "name", "display_name", "owner", "type", "created_at").
		From("mailboxes").
		Where(squirrel.And{squirrel.Eq{"owner": user}, squirrel.Eq{"name": name}}).
		QueryRow()

	var mailbox Mailbox
	if err := row.Scan(&mailbox.ID, &mailbox.Name, &mailbox.DisplayName, &mailbox.Owner, &mailbox.Type, &mailbox.CreatedAt); err != nil {
		return nil, err
	}

	return &mailbox, nil
}

func (s *Storage) QueryMailboxMetadata(user uint64, name string) (*MailboxMetadata, error) {
	row := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("count(i.*)", "max(i.uid)").
		From("mailboxes m").
		LeftJoin("inbounds_mailboxes i on m.id = i.mailbox").
		Where(squirrel.Eq{
			"m.owner": user,
			"m.name":  name,
		}).
		GroupBy("m.id")

	var metadata MailboxMetadata
	if err := row.Scan(&metadata.InboundCount, &metadata.UID); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (s *Storage) QueryMailboxWithMetadata(user uint64, name string) (*Mailbox, *MailboxMetadata, error) {
	row := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("count(i.*)", "max(i.uid)", "id", "name", "display_name", "owner", "type", "created_at").
		From("mailboxes m").
		LeftJoin("inbounds_mailboxes i on m.id = i.mailbox").
		Where(squirrel.Eq{
			"m.owner": user,
			"m.name":  name,
		}).
		GroupBy("m.id")

	var mailbox Mailbox
	var metadata MailboxMetadata
	if err := row.Scan(&metadata.InboundCount, &metadata.UID, &mailbox.ID, &mailbox.Name, &mailbox.DisplayName, &mailbox.Owner, &mailbox.Type, &mailbox.CreatedAt); err != nil {
		return nil, nil, err
	}

	return &mailbox, &metadata, nil
}

func (s *Storage) QueryMailboxCurrentUID() (uint, error) {
	row := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("max(uid)").
		From("inbounds").
		QueryRow()

	var uid uint
	if err := row.Scan(&uid); err != nil {
		return 0, err
	}

	return uid, nil
}
