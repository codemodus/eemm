package main

import (
	"fmt"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type sessionConfig struct {
	server, port, account, password string
}

type session struct {
	*coms
	cl  *client.Client
	dlm string
	fid string
	cnf sessionConfig
}

func newSession(cs *coms, id string, cnf sessionConfig) (*session, error) {
	s := &session{
		coms: cs,
		cnf:  cnf,
		fid:  fmt.Sprintf("%s: ", id),
	}

	if err := s.dial(); err != nil {
		return nil, err
	}
	s.logf("connected to %s on port %s", s.cnf.server, s.cnf.port)

	if err := s.login(); err != nil {
		return nil, err
	}
	s.logf("logged in as %s", s.cnf.account)

	if err := s.setDelim(); err != nil {
		return nil, err
	}
	s.logf("obtained delimiter")

	return s, nil
}

func (s *session) logf(format string, args ...interface{}) {
	s.Infof(s.fid+format, args...)
}

func (s *session) logerr(err error) {
	s.Error(s.fid + err.Error())
}

func (s *session) close() {
	if s.cl != nil {
		_ = s.cl.Logout()
		_ = s.cl.Close()
	}
}

func (s *session) dial() error {
	select {
	case <-s.done():
		return s.ErrShutdown
	default:
	}

	c, err := client.DialTLS(fmt.Sprintf("%s:%s", s.cnf.server, s.cnf.port), nil)
	if err != nil {
		return err
	}

	s.cl = c

	return nil
}

func (s *session) login() error {
	select {
	case <-s.done():
		return s.ErrShutdown
	default:
	}

	if s.cl == nil {
		return fmt.Errorf("missing client in session")
	}

	return s.cl.Login(s.cnf.account, s.cnf.password)
}

func (s *session) setDelim() error {
	select {
	case <-s.done():
		return s.ErrShutdown
	default:
	}

	ic := make(chan *imap.MailboxInfo, 10)
	ec := make(chan error)
	defer close(ec)

	go func() {
		ec <- s.cl.List("", "", ic)
	}()

	for mi := range ic {
		s.dlm = mi.Delimiter
	}

	return <-ec
}

func (s *session) syncTo(dst *session) error {
	mis, err := mailboxInfos(s, "")
	if err != nil {
		return err
	}
	s.logf("obtained mailbox info")

	if err = addMissingBoxes(dst, mis); err != nil {
		return err
	}
	dst.logf("normalized mailboxes")

	if true {
		// die for now
		return nil
	}

	mc := make(chan *message)
	ec := make(chan error)

	go func() {
		ec <- msgFeed(s, mis, mc)
	}()

	if err = addMissingMsgs(dst, mc); err != nil {
		return err
	}

	return <-ec
}
