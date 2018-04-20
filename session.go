package main

import (
	"fmt"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type imapClient struct {
	*client.Client
	delim    string
	pathprfx string
}

type sessionConfig struct {
	server   string
	port     string
	pathprfx string
	account  string
	password string
}

type session struct {
	*coms
	cnf sessionConfig
	cl  *imapClient
}

func newSession(cs *coms, l imap.Logger, cnf sessionConfig) (*session, error) {
	s := &session{
		coms: cs,
		cnf:  cnf,
		cl: &imapClient{
			pathprfx: cnf.pathprfx,
		},
	}

	if err := s.dial(l); err != nil {
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
	if s.cl.Client == nil {
		return
	}

	if err := recover(); err != nil {
		fmt.Println(err)
		return
	}

	_ = s.cl.Logout()
	_ = s.cl.Terminate()
}

func (s *session) dial(l imap.Logger) error {
	if err := s.term(); err != nil {
		return err
	}

	if s.cnf.port != "993" {
		cl, err := client.Dial(fmt.Sprintf("%s:%s", s.cnf.server, s.cnf.port))
		if err != nil {
			return err
		}

		cl.ErrorLog = l
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

func (s *session) checkTerm() error {
	if err := s.term(); err != nil {
		return err
	}

	select {
	case <-s.cl.LoggedOut():
		return fmt.Errorf("unexpectedly logged out")
	default:
	}

	return nil
}

func (s *session) login() error {
	if err := s.ensureClient(); err != nil {
		return err
	}

	return s.cl.Login(s.cnf.account, s.cnf.password)
}

func (s *session) setDelim() error {
	if err := s.checkTerm(); err != nil {
		return err
	}

	ic := make(chan *imap.MailboxInfo, 20)
	ec := make(chan error, 1)
	defer close(ec)

	go func() {
		ec <- s.cl.List("", "", ic)
	}()

	for mi := range ic {
		s.cl.delim = mi.Delimiter
	}

	return <-ec
}

func (s *session) replicateMailboxes(dst *session) ([]*imapMailboxInfo, error) {
	if err := checkTerms(dst, s); err != nil {
		return nil, err
	}

	mis, err := missingMailboxInfos(dst.cl, s.cl)
	if err != nil {
		return nil, err
	}

	for k, mi := range mis {
		if err := s.term(); err != nil {
			return mis[:k], err
		}

		if err := addMailbox(dst.cl, mi); err != nil {
			return mis[:k], err
		}
	}

	return mis, nil
}

func (s *session) replicateMessages(dst *session) error {
	if err := checkTerms(dst, s); err != nil {
		return err
	}

	mis, err := mailboxInfos(s.cl, "")
	if err != nil {
		return err
	}

	setSafeNames(dst.cl, s.cl, mis)

	for _, mi := range mis {
		uids, err := missingUIDs(dst.cl, s.cl, mi)
		if err != nil {
			return err
		}

		for _, seq := range uniqIDs(uids).chunks(32) {
			fms, err := messages(s.donec, s.cl, mi, seq)
			if err != nil {
				return err
			}

			if err = addMsgs(s.donec, dst.cl, mi, fms); err != nil {
				return err
			}
		}
	}

	return nil
}

func checkTerms(dst, src *session) error {
	if err := src.checkTerm(); err != nil {
		return err
	}

	return dst.checkTerm()
}
