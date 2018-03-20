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

// surface
func newSession(cs *coms, id string, cnf sessionConfig) (*session, error) {
	s := &session{
		coms: cs,
		cnf:  cnf,
		fid:  fmt.Sprintf("%s: ", id),
	}

	s.logf("connecting to %s on port %s", s.cnf.server, s.cnf.port)
	if err := s.dial(); err != nil {
		s.logerr(err)
		return nil, err
	}

	s.logf("logging in as %s", s.cnf.account)
	if err := s.login(); err != nil {
		s.logerr(err)
		return nil, err
	}

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
	if err := s.term(); err != nil {
		return err
	}

	if s.cnf.port != "993" {
		c, err := client.Dial(fmt.Sprintf("%s:%s", s.cnf.server, s.cnf.port))
		if err != nil {
			return err
		}

		s.cl = c

		return nil
	}

	c, err := client.DialTLS(fmt.Sprintf("%s:%s", s.cnf.server, s.cnf.port), nil)
	if err != nil {
		return err
	}

	s.cl = c

	return nil
}

func (s *session) login() error {
	if err := s.term(); err != nil {
		return err
	}

	if s.cl == nil {
		return fmt.Errorf("missing client in session")
	}

	if err := s.cl.Login(s.cnf.account, s.cnf.password); err != nil {
		return err
	}

	return s.setDelim()
}

func (s *session) setDelim() error {
	if err := s.term(); err != nil {
		return err
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

// surface
func (s *session) regularize(dst *session) error {
	s.logf("obtaining mailbox info")
	srcMis, err := mailboxInfos(s, "")
	if err != nil {
		s.logerr(err)
		return err
	}

	dst.logf("regularizing mailboxes")
	if err = addMissingBoxes(dst, srcMis); err != nil {
		s.logerr(err)
		return err
	}

	if true {
		// die for now
		return nil
	}

	dstMis, err := mailboxInfos(s, "")
	if err != nil {
		return err
	}

	for _, mi := range dstMis {
		hs, err := msgHashes(dst, mi.Name)
		if err != nil {
			return err
		}

		mc := make(chan *message)
		ec := make(chan error)
		defer close(ec)

		go func() {
			ec <- missingMsgsFeed(s, mi, hs, mc)
		}()

		if err = addMsgs(dst, mi.Name, mc); err != nil {
			return err
		}

		if err = <-ec; err != nil {
			return err
		}
	}

	return nil
}
