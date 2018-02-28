package main

import (
	"errors"
	"fmt"
)

type bondedSession struct {
	*coms
	fid string
	dst *session
	src *session
}

func newBondedSession(cs *coms, id int, dstConf, srcConf sessionConfig) (*bondedSession, error) {
	bs := &bondedSession{
		coms: cs,
		fid:  fmt.Sprintf(" %03dBSS: ", id),
	}
	bs.logf("%s(%s) from %s(%s)", dstConf.account, dstConf.server, srcConf.account, srcConf.server)

	dst, err := newSession(cs, fmt.Sprintf("   %03ddst", id), dstConf)
	if err != nil {
		return nil, err
	}
	bs.dst = dst
	bs.logf("established session for destination")

	src, err := newSession(cs, fmt.Sprintf("     %03dsrc", id), srcConf)
	if err != nil {
		return nil, err
	}
	bs.src = src
	bs.logf("established session for source (%s)", src.cf.account)

	return bs, nil
}

func (s *bondedSession) logf(format string, args ...interface{}) {
	s.ic <- fmt.Sprintf(s.fid+format, args...)
}

func (s *bondedSession) logerr(err error) {
	s.ec <- errors.New(s.fid + err.Error())
}

func (s *bondedSession) close() {
	if s.dst != nil {
		s.dst.close()
	}

	if s.src != nil {
		s.src.close()
	}
}

func (s *bondedSession) sync() error {
	return s.src.syncTo(s.dst)
}
