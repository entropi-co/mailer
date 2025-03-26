package imap

import (
	"bytes"
	"errors"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	"github.com/sirupsen/logrus"
	"mailer/internal/storage"
	"sync"
	"unsafe"
)

type Mailbox struct {
	IMAP *IMAP

	Source   *storage.Mailbox
	Metadata *storage.MailboxMetadata
	Tracker  *imapserver.MailboxTracker

	Views []*MailboxView

	lock sync.Mutex
}

type MailboxView struct {
	Mailbox *Mailbox
	Session *Session
	Tracker *imapserver.SessionTracker
}

func (i *IMAP) LoadMailbox(source *storage.Mailbox, metadata *storage.MailboxMetadata) *Mailbox {
	existing, hasExisting := i.Mailboxes[source.ID]
	if hasExisting {
		existing.Source = source
		existing.Metadata = metadata
		existing.Tracker.QueueNumMessages(uint32(metadata.InboundCount))

		return existing
	}

	mailbox := &Mailbox{
		Source:   source,
		Metadata: metadata,
		Tracker:  imapserver.NewMailboxTracker(uint32(metadata.InboundCount)),

		lock: sync.Mutex{},
	}

	i.Mailboxes[mailbox.Source.ID] = mailbox

	return mailbox
}

// Close cleans up children view and unload from IMAP registry
func (m *Mailbox) Close() {
	m.lock.Lock()
	for _, view := range m.Views {
		view.Close()
	}

	delete(m.IMAP.Mailboxes, m.Source.ID)

	m.lock.Unlock()
}

func (m *Mailbox) CreateView(session *Session) *MailboxView {
	view := &MailboxView{
		Mailbox: m,
		Session: session,
		Tracker: m.Tracker.NewSession(),
	}

	m.lock.Lock()
	m.Views = append(m.Views, view)
	m.lock.Unlock()

	return view
}

func (v *MailboxView) Close() {
	v.Tracker.Close()
	if len(v.Mailbox.Views) == 1 {
		v.Mailbox.Close()
	}

	// Remove this view from mailbox
	for i := range v.Mailbox.Views {
		if v.Mailbox.Views[i] == v {
			v.Mailbox.Views = append(v.Mailbox.Views[:i], v.Mailbox.Views[i+1:]...)
			break
		}
	}
}

func (v *MailboxView) queryInboundsFromSeqSet(seqSet *imap.SeqSet) ([]*storage.InboundWithMetadata, error) {
	nums, ok := seqSet.Nums()
	if !ok {
		return nil, errors.New("failed to retrieve nums from SeqSet")
	}

	decodedNums := make([]uint32, len(nums))
	for n := range nums {
		num := nums[n]
		decoded := v.Tracker.DecodeSeqNum(num)
		if decoded != 0 {
			decodedNums = append(decodedNums, num)
		}
	}

	inbounds, err := v.Session.Storage().QueryInboundsBySequences(v.Mailbox.Source.ID, decodedNums)
	if err != nil {
		return nil, err
	}

	return inbounds, nil
}

func (v *MailboxView) queryInboundsFromUIDSet(uidSet *imap.UIDSet) ([]*storage.InboundWithMetadata, error) {
	nums, ok := uidSet.Nums()
	if !ok {
		return nil, errors.New("failed to retrieve nums from UIDSet")
	}
	casted := unsafe.Slice((*uint32)(&nums[0]), len(nums))

	inbounds, err := v.Session.Storage().QueryInboundsByUIDS(v.Mailbox.Source.ID, casted)
	if err != nil {
		return nil, err
	}

	return inbounds, nil
}

func (v *MailboxView) Fetch(w *imapserver.FetchWriter, numSet imap.NumSet, options *imap.FetchOptions) error {
	logrus.Printf("[mailbox@%d] fetch", v.Mailbox.Source.ID)

	var err error
	var inbounds []*storage.InboundWithMetadata

	switch numSet := numSet.(type) {
	case imap.SeqSet:
		inbounds, err = v.queryInboundsFromSeqSet(&numSet)
		if err != nil {
			return err
		}
	case imap.UIDSet:
		//inbounds, err = v.queryInboundsFromUIDSet(&numSet)
		//if err != nil {
		//	return err
		//}
	}

	for _, inbound := range inbounds {
		mw := w.CreateMessage(v.Tracker.EncodeSeqNum(inbound.Sequence))
		mw.WriteUID(imap.UID(inbound.UID))

		for _, section := range options.BodySection {
			buffer := imapserver.ExtractBodySection(bytes.NewReader([]byte(inbound.Body)), section)
			sw := mw.WriteBodySection(section, int64(len(buffer)))
			_, writeErr := sw.Write(buffer)
			closeErr := sw.Close()
			if writeErr != nil {
				return writeErr
			}
			if closeErr != nil {
				return closeErr
			}
		}

		for _, section := range options.BinarySection {
			buffer := imapserver.ExtractBinarySection(bytes.NewReader([]byte(inbound.Body)), section)
			sw := mw.WriteBinarySection(section, int64(len(buffer)))
			_, writeErr := sw.Write(buffer)
			closeErr := sw.Close()
			if writeErr != nil {
				return writeErr
			}
			if closeErr != nil {
				return closeErr
			}
		}

		for _, size := range options.BinarySectionSize {
			n := imapserver.ExtractBinarySectionSize(bytes.NewReader([]byte(inbound.Body)), size)
			mw.WriteBinarySectionSize(size, n)
		}

		err := mw.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
