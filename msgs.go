package main

import (
	imap "github.com/emersion/go-imap"
)

type hashLookup map[string]struct{}

type message struct {
	mailbox string
	*imap.Message
}

/*src := s
_ = src
wg := &sync.WaitGroup{}
wg.Add(2)

var dstMsgs, srcMsgs []*imap.Message

go func() {
	//dstMsgs = dst.feed()
	wg.Done()
}()
go func() {
	//srcMsgs = src.feed()
	wg.Done()
}()

wg.Wait()
dst.logf("done syncing")

	bsn, err := imap.ParseBodySectionName(imap.FetchBody)
	bsn.Peek = true
	if err != nil {
		return err
	}

	for _, msg := range srcMsgs {
		if err := dst.c.Append("INBOX", msg.Flags, msg.InternalDate, msg.GetBody(bsn)); err != nil {
			dst.logerr(fmt.Errorf("cannot append message: %s", err))
		}
	}
*/

func missingMsgsFeed(s *session, mi *imap.MailboxInfo, hashes hashLookup, mc chan *message) error {
	if err := s.term(); err != nil {
		return err
	}

	_ = s

	return nil
}

func msgHashes(s *session, name string) (hashLookup, error) {
	return nil, nil
}

//func msgFeed(s *session, mis []*imap.MailboxInfo, mc chan *message) error {
/*	mbName := "INBOX"
	mb, err := s.c.Select(mbName, false)
	if err != nil {
		s.logerr(err)
		return nil
	}

	// store flags with mailbox data instead
	s.logf("flags for %s: %s", mb.Name, mb.Flags)

	first, last := uint32(1), mb.Messages
	seq := &imap.SeqSet{}
	seq.AddRange(first, last)
	// store count instead
	s.logf("found %d messages", last)

	msgs := make(chan *imap.Message, 10)
	msgsErr := make(chan error, 1)
	go func() {
		fis := []imap.FetchItem{
			//imap.FetchRFC822,
			imap.FetchFlags,
			imap.FetchInternalDate,
			imap.FetchEnvelope,
			imap.FetchBody,
			imap.FetchUid,
		}
		msgsErr <- s.c.Fetch(seq, fis, msgs)
	}()

	var mbuf []*imap.Message

	for msg := range msgs {
		mbuf = append(mbuf, msg)
		// store messages instead
		s.logf("* %s as %d with %v", msg.Envelope.Subject, msg.Uid, msg.Flags)
	}
	if err := <-msgsErr; err != nil {
		s.logerr(err)
	}
*/

func addMsgs(s *session, name string, mc chan *message) error {
	if err := s.term(); err != nil {
		return err
	}

	_ = s

	return nil
}

func (s *session) sreplicateMessages(dst *session, mis []*imap.MailboxInfo) error {
	/*if err := s.ensureLogin(); err != nil {
		return err
	}

	for _, mi := range mis {

		hs, err := msgHashes(dst, mi.Name)
		if err != nil {
			return err
		}

		mc := make(chan *message)
		ec := make(chan error)
		defer close(ec)

		go func() {
			ec <- missingMsgsFeed(s, mi, hs, mc)
		}()

		if err = addMsgs(dst, mi.Name, mc); err != nil {
			return err
		}

		if err = <-ec; err != nil {
			return err
		}
	}
	*/
	return nil
}
