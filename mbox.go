package main

import (
	"fmt"
	"strings"

	imap "github.com/emersion/go-imap"
)

func mailboxInfos(cl *imapClient, name string) ([]*imap.MailboxInfo, error) {
	if err := checkDepth(name, cl.delim); err != nil {
		return nil, err
	}

	ic := make(chan *imap.MailboxInfo, 10)
	ec := make(chan error, 1)
	defer close(ec)

	go func() {
		ec <- cl.List("", listArg(name, cl.delim), ic)
	}()

	var mis []*imap.MailboxInfo

	for mi := range ic {
		mis = append(mis, mi)

		children, err := mailboxInfos(cl, mi.Name)
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

func addMissingBox(cl *imapClient, dstMis []*imap.MailboxInfo, srcMi *imap.MailboxInfo) error {
	dstName := delimAdjustedName(srcMi, cl.delim)

	for _, dstMi := range dstMis {
		if dstMi.Name == dstName {
			return nil
		}
	}

	return cl.Create(dstName)
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
