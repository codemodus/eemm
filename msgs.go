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
	if err := s.term(); err != nil {
		return err
	}

	_ = s

	return nil
}

func msgHashes(s *session, name string) (hashLookup, error) {
	return nil, nil
}

func addMsgs(s *session, name string, mc chan *message) error {
	if err := s.term(); err != nil {
		return err
	}

	_ = s

	return nil
}
