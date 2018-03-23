package main

import (
	"crypto/md5"

	imap "github.com/emersion/go-imap"
)

var (
	shortFetchItem = []imap.FetchItem{
		imap.FetchInternalDate,
		imap.FetchEnvelope,
		imap.FetchUid,
	}
	fullFetchItem = []imap.FetchItem{
		imap.FetchFlags,
		imap.FetchInternalDate,
		imap.FetchRFC822Size,
		imap.FetchEnvelope,
		imap.FetchBody,
		imap.FetchUid,
	}
)

const (
	hashLen = md5.Size
)

func appendAddrBytes(b []byte, adsGrp ...[]*imap.Address) []byte {
	for _, ads := range adsGrp {
		for _, v := range ads {
			b = append(b, []byte(v.MailboxName)...)
			b = append(b, []byte(v.PersonalName)...)
			b = append(b, []byte(v.AtDomainList)...)
			b = append(b, []byte(v.HostName)...)
		}
	}

	return b
}

func msgHash(m *imap.Message) [hashLen]byte {
	var b []byte

	b = append(b, []byte(m.InternalDate.String())...)
	b = append(b, []byte(m.Envelope.Date.String())...)
	b = append(b, []byte(m.Envelope.Subject)...)

	b = appendAddrBytes(
		b,
		m.Envelope.Sender,
		m.Envelope.From,
		m.Envelope.ReplyTo,
		m.Envelope.To,
		m.Envelope.Cc,
		m.Envelope.Bcc,
	)

	return md5.Sum(b)
}

func msgHashes(cl *imapClient, mi *imap.MailboxInfo) (map[[hashLen]byte]uint32, error) {
	hs := make(map[[hashLen]byte]uint32)

	mbName := delimAdjustedName(mi, cl.delim)
	mb, err := cl.Select(mbName, false)
	if err != nil {
		return hs, err
	}

	if mb.Messages == 0 {
		return hs, nil
	}

	seq := &imap.SeqSet{}
	seq.AddRange(1, mb.Messages)

	msgs := make(chan *imap.Message, 10)
	msgsErr := make(chan error, 1)
	go func() {
		msgsErr <- cl.Fetch(seq, shortFetchItem, msgs)
	}()

	for msg := range msgs {
		hs[msgHash(msg)] = msg.Uid
	}

	return hs, <-msgsErr
}

func missingUIDs(dst, src *imapClient, mi *imap.MailboxInfo) ([]uint32, error) {
	srcHs, err := msgHashes(src, mi)
	if err != nil {
		return nil, err
	}

	dstHs, err := msgHashes(dst, mi)
	if err != nil {
		return nil, err
	}

	var uids []uint32

	for hash, uid := range srcHs {
		if _, ok := dstHs[hash]; !ok {
			uids = append(uids, uid)
		}
	}

	return uids, nil
}

func messages(done chan struct{}, cl *imapClient, uids []uint32) ([]*imap.Message, error) {
	// TODO: this
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
	return nil, nil
}

func addMsgs(done chan struct{}, cl *imapClient, mi *imap.MailboxInfo, msgs []*imap.Message) error {
	// TODO: this
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
	return nil
}
