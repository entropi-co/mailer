package imap

import (
	"errors"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	"mailer/internal/instance"
	"mailer/internal/storage"
	"time"
)

type Session struct {
	Instance *instance.Instance
	User     *storage.User
	Mailbox  *storage.Mailbox
}

func (i *IMAP) NewSession() *Session {
	return &Session{
		Instance: i.Instance,
	}
}

func (s Session) Close() error {
	return nil
}

func (s Session) Login(username, password string) error {
	key, err := s.Instance.Storage.QueryKeyByLocal(username)
	if err != nil {
		return err
	}

	if key.Value != password {
		return errors.New("api key mismatch")
	}

	user, err := s.Instance.Storage.QueryUserByLocal(username)
	if err != nil {
		return err
	}

	s.User = user

	return nil
}

func (s Session) Select(name string, options *imap.SelectOptions) (*imap.SelectData, error) {
	mailbox, metadata, err := s.Instance.Storage.QueryMailboxWithMetadata(s.User.ID, name)
	if err != nil {
		return nil, err
	}

	s.Mailbox = mailbox

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

func (s Session) Create(mailbox string, options *imap.CreateOptions) error {
	//TODO implement me
	panic("implement me")
}

func (s Session) Delete(mailbox string) error {
	//TODO implement me
	panic("implement me")
}

func (s Session) Rename(mailbox, newName string) error {
	//TODO implement me
	panic("implement me")
}

func (s Session) Subscribe(mailbox string) error {
	//TODO implement me
	panic("implement me")
}

func (s Session) Unsubscribe(mailbox string) error {
	//TODO implement me
	panic("implement me")
}

func (s Session) List(w *imapserver.ListWriter, ref string, patterns []string, options *imap.ListOptions) error {
	//TODO implement me
	panic("implement me")
}

func (s Session) Status(mailbox string, options *imap.StatusOptions) (*imap.StatusData, error) {
	//TODO implement me
	panic("implement me")
}

func (s Session) Append(mailbox string, r imap.LiteralReader, options *imap.AppendOptions) (*imap.AppendData, error) {
	//TODO implement me
	panic("implement me")
}

func (s Session) Poll(w *imapserver.UpdateWriter, allowExpunge bool) error {
	//TODO implement me
	panic("implement me")
}

func (s Session) Idle(w *imapserver.UpdateWriter, stop <-chan struct{}) error {
	//TODO implement me
	panic("implement me")
}

func (s Session) Unselect() error {
	//TODO implement me
	panic("implement me")
}

func (s Session) Expunge(w *imapserver.ExpungeWriter, uids *imap.UIDSet) error {
	//TODO implement me
	panic("implement me")
}

func (s Session) Search(kind imapserver.NumKind, criteria *imap.SearchCriteria, options *imap.SearchOptions) (*imap.SearchData, error) {
	//TODO implement me
	panic("implement me")
}

func (s Session) Fetch(w *imapserver.FetchWriter, numSet imap.NumSet, options *imap.FetchOptions) error {

	panic("implement me")
}

func (s Session) Store(w *imapserver.FetchWriter, numSet imap.NumSet, flags *imap.StoreFlags, options *imap.StoreOptions) error {
	//TODO implement me
	panic("implement me")
}

func (s Session) Copy(numSet imap.NumSet, dest string) (*imap.CopyData, error) {
	//TODO implement me
	panic("implement me")
}
