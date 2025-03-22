package imap

import (
	"errors"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	"mailer/internal/storage"
	"sync"
)

type Mailbox struct {
	Source   *storage.Mailbox
	Metadata *storage.MailboxMetadata
	Tracker  *imapserver.MailboxTracker

	Views []*MailboxView

	lock sync.Mutex
}

type MailboxView struct {
	Mailbox *Mailbox
	Tracker *imapserver.SessionTracker
}

func (i *IMAP) LoadMailbox(source *storage.Mailbox, metadata *storage.MailboxMetadata) {
	existing, hasExisting := i.Mailboxes[source.ID]
	if hasExisting {
		existing.Source = source
		existing.Metadata = metadata
		existing.Tracker.QueueNumMessages(uint32(metadata.InboundCount))

		return
	}

	mailbox := &Mailbox{
		Source:   source,
		Metadata: metadata,
		Tracker:  imapserver.NewMailboxTracker(uint32(metadata.InboundCount)),

		lock: sync.Mutex{},
	}

	i.Mailboxes[mailbox.Source.ID] = mailbox
}

func (v *MailboxView) Close() {
	v.Tracker.Close()
	if len(v.Mailbox.Views) == 1 {
		// TODO: Unload mailbox as no views are open
	}

	// Remove this view from mailbox
	for i := range v.Mailbox.Views {
		if v.Mailbox.Views[i] == v {
			v.Mailbox.Views = append(v.Mailbox.Views[:i], v.Mailbox.Views[i+1:]...)
			break
		}
	}
}

func (v *MailboxView) Fetch(w *imapserver.FetchWriter, numSet imap.NumSet, options *imap.FetchOptions) error {
	switch numSet := numSet.(type) {
	case imap.SeqSet:
		nums, ok := numSet.Nums()
		if !ok {
			return errors.New("failed to retrieve nums from numSet")
		}

		decodedNums := make([]uint32, len(nums))
		for n := range nums {
			num := nums[n]
			decoded := v.Tracker.DecodeSeqNum(num)
			if decoded != 0 {
				decodedNums = append(decodedNums, num)
			}
		}
	case imap.NumSet:

	}

	panic("unimplemented")

	return nil
}
