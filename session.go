package main

import (
	"fmt"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type imapClient struct {
	*client.Client
	delim string
}

type sessionConfig struct {
	server   string
	port     string
	account  string
	password string
}

type session struct {
	*coms
	cnf sessionConfig
	cl  *imapClient
}

func newSession(cs *coms, cnf sessionConfig) (*session, error) {
	s := &session{
		coms: cs,
		cnf:  cnf,
		cl:   &imapClient{},
	}

	if err := s.dial(); err != nil {
		return nil, fmt.Errorf("cannot dial into %s on %s: %s", cnf.port, cnf.server, err)
	}

	if err := s.login(); err != nil {
		return nil, fmt.Errorf("cannot login as %s: %s", cnf.account, err)
	}

	if err := s.setDelim(); err != nil {
		return nil, fmt.Errorf("cannot set delimiter: %s", err)
	}

	return s, nil
}

func (s *session) close() {
	if s.cl.Client != nil {
		_ = s.cl.Logout()
		_ = s.cl.Close()
	}
}

func (s *session) dial() error {
	if err := s.term(); err != nil {
		return err
	}

	if s.cnf.port != "993" {
		cl, err := client.Dial(fmt.Sprintf("%s:%s", s.cnf.server, s.cnf.port))
		if err != nil {
			return err
		}

		s.cl.Client = cl

		return nil
	}

	cl, err := client.DialTLS(fmt.Sprintf("%s:%s", s.cnf.server, s.cnf.port), nil)
	if err != nil {
		return err
	}

	s.cl.Client = cl

	return nil
}

func (s *session) ensureClient() error {
	if err := s.term(); err != nil {
		return err
	}

	// TODO: implement ensureClient correctly
	if s.cl == nil || s.cl.Client == nil {
		return fmt.Errorf("missing client in session")
	}

	return nil
}

func (s *session) login() error {
	if err := s.ensureClient(); err != nil {
		return err
	}

	return s.cl.Login(s.cnf.account, s.cnf.password)
}

func (s *session) ensureLogin() error {
	// TODO: implement ensureLogin correctly
	return s.ensureClient()
}

func (s *session) setDelim() error {
	if err := s.ensureLogin(); err != nil {
		return err
	}

	ic := make(chan *imap.MailboxInfo, 10)
	ec := make(chan error)
	defer close(ec)

	go func() {
		ec <- s.cl.List("", "", ic)
	}()

	for mi := range ic {
		s.cl.delim = mi.Delimiter
	}

	return <-ec
}

func (s *session) replicateMailboxes(dst *session) ([]*imap.MailboxInfo, error) {
	if err := s.ensureLogin(); err != nil {
		return nil, err
	}

	srcMis, err := mailboxInfos(s.cl, "")
	if err != nil {
		return nil, err
	}

	if err = dst.ensureLogin(); err != nil {
		return nil, err
	}

	dstMis, err := mailboxInfos(dst.cl, "")
	if err != nil {
		return nil, err
	}

	for k, mi := range srcMis {
		if err := s.term(); err != nil {
			return srcMis[:k], err
		}

		if err := addMissingBox(dst.cl, dstMis, mi); err != nil {
			return srcMis[:k], err
		}
	}

	return srcMis, nil
}

func (s *session) replicateMessages(dst *session, mis []*imap.MailboxInfo) error {
	return nil
}
