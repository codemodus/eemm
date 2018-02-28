package main

import (
	"fmt"
	"strings"

	imap "github.com/emersion/go-imap"
)

func mailboxInfos(s *session, name string) ([]*imap.MailboxInfo, error) {
	select {
	case <-s.done():
		return nil, s.ErrShutdown
	default:
	}

	if err := checkDepth(name, s.dlm); err != nil {
		return nil, err
	}

	ic := make(chan *imap.MailboxInfo, 10)
	ec := make(chan error, 1)
	defer close(ec)

	go func() {
		ec <- s.cl.List("", listArg(name, s.dlm), ic)
	}()

	var mis []*imap.MailboxInfo

	for mi := range ic {
		mis = append(mis, mi)

		children, err := mailboxInfos(s, mi.Name)
		if err != nil {
			drainMailboxInfo(ic, ec)
			return nil, err
		}

		mis = append(mis, children...)
	}

	if err := <-ec; err != nil {
		return nil, err
	}

	return mis, nil
}

func addMissingBoxes(s *session, mis []*imap.MailboxInfo) error {
	dstMis, err := mailboxInfos(s, "")
	if err != nil {
		return err
	}

	for _, mi := range mis {
		if err := addMissingBox(s, dstMis, mi); err != nil {
			return err
		}
	}

	return nil
}

func addMissingBox(s *session, dstMis []*imap.MailboxInfo, srcMi *imap.MailboxInfo) error {
	select {
	case <-s.done():
		return s.ErrShutdown
	default:
	}

	dstName := delimAdjustedName(srcMi, s.dlm)

	for _, dstMi := range dstMis {
		if dstMi.Name == dstName {
			return nil
		}
	}

	return s.cl.Create(dstName)
}

func drainMailboxInfo(c chan *imap.MailboxInfo, ec chan error) {
	for range c {
	}
	<-ec
}

func delimAdjustedName(mi *imap.MailboxInfo, delim string) string {
	return strings.Replace(mi.Name, mi.Delimiter, delim, -1)
}

func checkDepth(name, delim string) error {
	if strings.Count(name, delim) > 18 {
		return fmt.Errorf("max depth encountered %s", name)
	}

	return nil
}

func listArg(name, delim string) string {
	return strings.TrimLeft(name+delim+"%", delim)
}
