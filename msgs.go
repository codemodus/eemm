package main

import (
	imap "github.com/emersion/go-imap"
)

type hashLookup map[string]struct{}

type message struct {
	mailbox string
	*imap.Message
}

func missingMsgsFeed(s *session, mi *imap.MailboxInfo, hashes hashLookup, mc chan *message) error {
	select {
	case <-s.done():
		return s.ErrShutdown
	default:
	}

	return nil
}

func msgHashes(s *session, name string) (hashLookup, error) {
	return nil, nil
}

func addMsgs(s *session, name string, mc chan *message) error {
	select {
	case <-s.done():
		return s.ErrShutdown
	default:
	}

	return nil
}
