package main

import (
	"fmt"
)

type bondedSession struct {
	*coms
	fid string
	dst *session
	src *session
}

// surface
func newBondedSession(cs *coms, id int, dstConf, srcConf sessionConfig) (*bondedSession, error) {
	bs := &bondedSession{
		coms: cs,
		fid:  fmt.Sprintf(" %03dBSS: ", id),
	}
	bs.logf("bonding to %s(%s) from %s(%s)", dstConf.account, dstConf.server, srcConf.account, srcConf.server)

	dst, err := newSession(cs, fmt.Sprintf("   %03ddst", id), dstConf)
	if err != nil {
		return nil, err
	}
	bs.dst = dst

	src, err := newSession(cs, fmt.Sprintf("     %03dsrc", id), srcConf)
	if err != nil {
		return nil, err
	}
	bs.src = src

	return bs, nil
}

func (s *bondedSession) logf(format string, args ...interface{}) {
	s.Infof(s.fid+format, args...)
}

func (s *bondedSession) logerr(err error) {
	s.Error(s.fid + err.Error())
}

func (s *bondedSession) close() {
	if s.dst != nil {
		s.dst.close()
	}

	if s.src != nil {
		s.src.close()
	}
}

// surface
func (s *bondedSession) regularize() error {
	s.logf("regularizing mailboxes")
	if err := s.src.regularize(s.dst); err != nil {
		s.logerr(fmt.Errorf("cannot regularize mailboxes"))
		return err
	}

	return nil
}
