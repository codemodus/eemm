package main

import (
	"fmt"

	imap "github.com/emersion/go-imap"
)

type bondedSession struct {
	*coms
	dst *session
	src *session
}

func makeBondedSession(cs *coms, dstConf, srcConf sessionConfig) (bondedSession, error) {
	bs := bondedSession{
		coms: cs,
	}

	dst, err := newSession(cs, dstConf)
	if err != nil {
		return bs, fmt.Errorf("cannot create destination session: %s", err)
	}
	bs.dst = dst

	src, err := newSession(cs, srcConf)
	if err != nil {
		return bs, fmt.Errorf("cannot create source session: %s", err)
	}
	bs.src = src

	return bs, nil
}

func (s *bondedSession) close() {
	if s.dst != nil {
		s.dst.close()
	}

	if s.src != nil {
		s.src.close()
	}
}

func (s *bondedSession) replicateMailboxes() ([]*imap.MailboxInfo, error) {
	return s.src.replicateMailboxes(s.dst)
}

func (s *bondedSession) replicateMessages() error {
	return s.src.replicateMessages(s.dst)
}
