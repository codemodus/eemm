package main

import (
	"fmt"
	"strings"

	imap "github.com/emersion/go-imap"
)

func mailboxInfos(s *session, name string) ([]*imap.MailboxInfo, error) {
	if err := checkDepth(name, s.dlm); err != nil {
		return nil, err
	}

	ic := make(chan *imap.MailboxInfo, 10)
	ec := make(chan error, 1)
	defer close(ec)

	go func() {
		ec <- s.c.List("", listArg(name, s.dlm), ic)
	}()

	var mis []*imap.MailboxInfo
	var anyErr bool

	for mi := range ic {
		mis = append(mis, mi)

		children, err := mailboxInfos(s, mi.Name)
		if err != nil {
			s.logerr(err)
			anyErr = true

			continue
		}

		mis = append(mis, children...)
	}

	if anyErr {
		return nil, fmt.Errorf("cannot recurse boxes; check log")
	}

	return mis, <-ec
}

func addMissingBoxes(s *session, mis []*imap.MailboxInfo) error {
	curMis, err := mailboxInfos(s, "")
	if err != nil {
		return err
	}

	for _, mi := range mis {
		name := delimAdjustedName(mi, s.dlm)
		var found bool

		for _, curMi := range curMis {
			if curMi.Name == name {
				found = true
				break
			}
		}

		if found {
			continue
		}

		if err := s.c.Create(name); err != nil {
			return err
		}
	}

	return nil
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
