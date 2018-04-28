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

func makeBondedSession(cs *coms, l imap.Logger, dstConf, srcConf sessionConfig) (bondedSession, error) {
	bs := bondedSession{
		coms: cs,
	}

	dst, err := newSession(cs, l, dstConf)
	if err != nil {
		return bs, fmt.Errorf("cannot create destination session: %s", err)
	}
	bs.dst = dst

	src, err := newSession(cs, l, srcConf)
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

func (s *bondedSession) replicateMailboxes(glblExcl, lclExcl []string) ([]*imapMailboxInfo, error) {
	return s.src.replicateMailboxes(s.dst, glblExcl, lclExcl)
}

func (s *bondedSession) replicateMessages(glblExcl, lclExcl []string) error {
	return s.src.replicateMessages(s.dst, glblExcl, lclExcl)
}
