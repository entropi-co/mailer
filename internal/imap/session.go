package imap

import (
	"database/sql"
	"errors"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	"mailer/internal/instance"
	"mailer/internal/storage"
	"time"
)

type Session struct {
	IMAP *IMAP
	User *storage.User
	View *MailboxView
}

func (i *IMAP) NewSession() *Session {
	return &Session{
		IMAP: i,
	}
}

// Instance is shorthand for IMAP.Instance
func (s *Session) Instance() *instance.Instance {
	return s.IMAP.Instance
}

// Storage is shorthand for IMAP.Instance.Storage
func (s *Session) Storage() *storage.Storage {
	return s.IMAP.Instance.Storage
}

func (s *Session) Close() error {
	return nil
}

func (s *Session) Login(username, password string) error {
	user, err := s.IMAP.Instance.Storage.QueryUserByLocalAndKeyValue(username, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return imapserver.ErrAuthFailed
		}

		return err
	}

	s.User = user

	return nil
}

func (s *Session) Select(name string, options *imap.SelectOptions) (*imap.SelectData, error) {
	source, metadata, err := s.IMAP.Instance.Storage.QueryMailboxWithMetadata(s.User.ID, name)
	if err != nil {
		return nil, err
	}

	mailbox := s.IMAP.LoadMailbox(source, metadata)
	view := mailbox.CreateView(s)

	s.View = view

	return &imap.SelectData{
		Flags:          nil,
		PermanentFlags: nil,
		NumMessages:    uint32(metadata.InboundCount),
		UIDNext:        imap.UID(metadata.UID + 1),
		UIDValidity:    uint32(time.Now().UnixMilli()),
		List:           nil,
		HighestModSeq:  0,
	}, nil
}

func (s *Session) Create(mailbox string, options *imap.CreateOptions) error {
	//TODO implement me
	panic("implement me")
}

func (s *Session) Delete(mailbox string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Session) Rename(mailbox, newName string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Session) Subscribe(mailbox string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Session) Unsubscribe(mailbox string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Session) List(w *imapserver.ListWriter, ref string, patterns []string, options *imap.ListOptions) error {
	//TODO implement me
	panic("implement me")
}

func (s *Session) Status(mailbox string, options *imap.StatusOptions) (*imap.StatusData, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Session) Append(mailbox string, r imap.LiteralReader, options *imap.AppendOptions) (*imap.AppendData, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Session) Poll(w *imapserver.UpdateWriter, allowExpunge bool) error {
	err := s.View.Tracker.Poll(w, allowExpunge)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Idle(w *imapserver.UpdateWriter, stop <-chan struct{}) error {
	err := s.View.Tracker.Idle(w, stop)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Unselect() error {
	s.View.Close()
	s.View = nil

	return nil
}

func (s *Session) Expunge(w *imapserver.ExpungeWriter, uids *imap.UIDSet) error {
	//TODO implement me
	panic("implement me")
}

func (s *Session) Search(kind imapserver.NumKind, criteria *imap.SearchCriteria, options *imap.SearchOptions) (*imap.SearchData, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Session) Fetch(w *imapserver.FetchWriter, numSet imap.NumSet, options *imap.FetchOptions) error {
	return s.View.Fetch(w, numSet, options)
}

func (s *Session) Store(w *imapserver.FetchWriter, numSet imap.NumSet, flags *imap.StoreFlags, options *imap.StoreOptions) error {
	//TODO implement me
	panic("implement me")
}

func (s *Session) Copy(numSet imap.NumSet, dest string) (*imap.CopyData, error) {
	//TODO implement me
	panic("implement me")
}
