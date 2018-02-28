package main

import (
	imap "github.com/emersion/go-imap"
)

type message struct {
	mailbox string
	*imap.Message
}

func msgFeed(s *session, mis []*imap.MailboxInfo, mc chan *message) error {
	select {
	case <-s.dc:
		return s.sd
	default:
	}

	return nil
}

func addMissingMsgs(s *session, mc chan *message) error {
	select {
	case <-s.dc:
		return s.sd
	default:
	}

	return nil
}
