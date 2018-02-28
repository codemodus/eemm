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
	c   *client.Client
	dlm string
	fid string
	cf  sessionConfig
}

func newSession(cs *coms, id string, cf sessionConfig) (*session, error) {
	s := &session{
		coms: cs,
		cf:   cf,
		fid:  fmt.Sprintf("%s: ", id),
	}

	if err := s.dial(); err != nil {
		return nil, err
	}
	s.logf("connected to %s on port %s", s.cf.server, s.cf.port)

	if err := s.login(); err != nil {
		return nil, err
	}
	s.logf("logged in as %s", s.cf.account)

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
	if s.c != nil {
		_ = s.c.Logout()
		_ = s.c.Close()
	}
}

func (s *session) dial() error {
	select {
	case <-s.dc:
		return s.sd
	default:
	}

	c, err := client.DialTLS(fmt.Sprintf("%s:%s", s.cf.server, s.cf.port), nil)
	if err != nil {
		return err
	}

	s.c = c

	return nil
}

func (s *session) login() error {
	select {
	case <-s.dc:
		return s.sd
	default:
	}

	if s.c == nil {
		return fmt.Errorf("missing client in session")
	}

	return s.c.Login(s.cf.account, s.cf.password)
}

func (s *session) setDelim() error {
	select {
	case <-s.dc:
		return s.sd
	default:
	}

	ic := make(chan *imap.MailboxInfo, 10)
	ec := make(chan error)
	defer close(ec)

	go func() {
		ec <- s.c.List("", "", ic)
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

	if err = <-ec; err != nil {
		_ = err
		return err
	}

	return nil
}
