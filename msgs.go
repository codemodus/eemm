package main

import (
	"crypto/md5"
	"strconv"
	"unicode"

	imap "github.com/emersion/go-imap"
)

var (
	shortFetchItem = []imap.FetchItem{
		imap.FetchInternalDate,
		imap.FetchEnvelope,
		imap.FetchRFC822Size,
		imap.FetchUid,
	}
	fullFetchItem = []imap.FetchItem{
		imap.FetchFlags,
		imap.FetchInternalDate,
		imap.FetchRFC822Size,
		imap.FetchEnvelope,
		imap.FetchUid,
		"BODY.PEEK[]",
	}
)

const (
	hashLen = md5.Size
)

func msgHash(m *imap.Message) [hashLen]byte {
	var b []byte

	b = append(b, []byte(strconv.FormatInt(m.Envelope.Date.Unix(), 10))...)
	b = append(b, stripSpaceToBytes(m.Envelope.Subject)...)
	b = append(b, []byte(strconv.FormatUint(uint64(m.Size), 10))...)

	b = addresses(m.Envelope.From).appendBytesTo(b)
	b = addresses(m.Envelope.To).appendBytesTo(b)
	b = addresses(m.Envelope.Cc).appendBytesTo(b)
	b = addresses(m.Envelope.Bcc).appendBytesTo(b)

	return md5.Sum(b)
}

func msgHashes(cl *imapClient, mi *imapMailboxInfo, fromSafe bool) (map[[hashLen]byte]uint32, error) {
	hs := make(map[[hashLen]byte]uint32)

	prepFn := mi.preparedName
	if fromSafe {
		prepFn = mi.preparedSafeName
	}

	mb, err := cl.Select(prepFn(cl.delim, cl.pathprfx), false)
	if err != nil {
		return hs, err
	}

	if mb.Messages == 0 {
		return hs, nil
	}

	seq := &imap.SeqSet{}
	seq.AddRange(1, mb.Messages)

	msgs := make(chan *imap.Message, 20)
	msgsErr := make(chan error, 1)
	go func() {
		msgsErr <- cl.Fetch(seq, shortFetchItem, msgs)
	}()

	for msg := range msgs {
		if msg.Size == 0 {
			continue
		}

		hs[msgHash(msg)] = msg.Uid
	}

	return hs, <-msgsErr
}

func missingUIDs(dst, src *imapClient, mi *imapMailboxInfo) ([]uint32, error) {
	srcHs, err := msgHashes(src, mi, false)
	if err != nil {
		return nil, err
	}

	dstHs, err := msgHashes(dst, mi, true)
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

func messages(done chan struct{}, cl *imapClient, mi *imapMailboxInfo, seq *imap.SeqSet) ([]*imap.Message, error) {
	var ms []*imap.Message

	mb, err := cl.Select(mi.preparedName(cl.delim, cl.pathprfx), false)
	if err != nil {
		return nil, err
	}

	if mb.Messages == 0 {
		return nil, nil
	}

	msgs := make(chan *imap.Message, 20)
	msgsErr := make(chan error, 1)
	go func() {
		msgsErr <- cl.UidFetch(seq, fullFetchItem, msgs)
	}()

	for msg := range msgs {
		ms = append(ms, msg)
	}

	return ms, <-msgsErr
}

func addMsgs(done chan struct{}, cl *imapClient, mi *imapMailboxInfo, msgs []*imap.Message) error {
	mbName := mi.preparedSafeName(cl.delim, cl.pathprfx)

	bsn, err := imap.ParseBodySectionName(imap.FetchRFC822)
	bsn.Peek = true
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		if err := cl.Append(mbName, msg.Flags, msg.InternalDate, msg.GetBody(bsn)); err != nil {
			return err
		}
	}

	return nil
}

func stripSpaceToBytes(s string) []byte {
	b := make([]byte, 0, len(s))

	for _, r := range s {
		if !unicode.IsSpace(r) {
			b = append(b, []byte(string(r))...)
		}

	}

	return b
}
