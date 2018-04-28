package main

import (
	"fmt"
	"strings"

	imap "github.com/emersion/go-imap"
)

type imapMailboxInfo struct {
	*imap.MailboxInfo
	altName string
}

func (i *imapMailboxInfo) alternateName(prfx string) string {
	if len(i.Name) <= len(prfx) {
		return i.Name
	}

	return prfx + "__" + i.Name[len(prfx)+1:]
}

func (i *imapMailboxInfo) safeName() string {
	if i.altName != "" {
		return i.altName
	}

	return i.Name
}

func (i *imapMailboxInfo) preparedName(delim, prfx string) string {
	return preparedName(i.Name, i.Delimiter, delim, prfx)
}

func (i *imapMailboxInfo) preparedSafeName(delim, prfx string) string {
	return preparedName(i.safeName(), i.Delimiter, delim, prfx)
}

func (i *imapMailboxInfo) trimmedName(prfx string) string {
	pLen := len(prfx)
	if pLen > 0 && len(i.Name) > pLen && i.Name[:pLen] == prfx {
		return i.Name[pLen+1:]
	}

	return i.Name
}

func (i *imapMailboxInfo) startsWith(delim, prfx string) bool {
	prfxAndDelim := prfx + delim

	return len(i.Name) >= len(prfxAndDelim) && i.Name[:len(prfxAndDelim)] == prfxAndDelim
}

func setSafeNames(dst, src *imapClient, srcMis []*imapMailboxInfo) {
	for _, smi := range srcMis {
		if src.pathprfx != dst.pathprfx && smi.startsWith(src.delim, dst.pathprfx) {
			smi.altName = smi.alternateName(dst.pathprfx)
		}
	}
}

func mailboxInfos(cl *imapClient, name string, glblExcl, lclExcl []string) ([]*imapMailboxInfo, error) {
	if err := checkDepth(name, cl.delim); err != nil {
		return nil, err
	}

	// TODO: skip if excluded.

	ic := make(chan *imap.MailboxInfo, 80)
	ec := make(chan error, 1)
	defer close(ec)

	go func() {
		ec <- cl.List("", listArg(name, cl.delim), ic)
	}()

	var mis []*imapMailboxInfo

	for imi := range ic {
		mi := &imapMailboxInfo{MailboxInfo: imi}
		imiName := mi.Name
		mi.Name = mi.trimmedName(cl.pathprfx)

		mis = append(mis, mi)

		if mailboxHasNoChildren(imi) {
			continue
		}

		children, err := mailboxInfos(cl, imiName, glblExcl, lclExcl)
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

func missingMailboxInfos(dst, src *imapClient, glblExcl, lclExcl []string) ([]*imapMailboxInfo, error) {
	srcMis, err := mailboxInfos(src, "", glblExcl, lclExcl)
	if err != nil {
		return nil, err
	}

	dstMis, err := mailboxInfos(dst, "", glblExcl, lclExcl)
	if err != nil {
		return nil, err
	}

	mis := srcMis[:0]

	setSafeNames(dst, src, srcMis)

	for _, smi := range srcMis {
		found := false

		for _, dmi := range dstMis {
			if smi.safeName() == dmi.preparedName(src.delim, "") {
				found = true
				break
			}
		}

		if !found {
			mis = append(mis, smi)
		}
	}

	return mis, nil
}

func addMailbox(cl *imapClient, srcMi *imapMailboxInfo) error {
	dstName := srcMi.preparedSafeName(cl.delim, cl.pathprfx)

	return cl.Create(dstName)
}

func drainMailboxInfo(c chan *imap.MailboxInfo, ec chan error) {
	for range c {
	}

	select {
	case <-ec:
	default:
	}
}

func mailboxHasNoChildren(mi *imap.MailboxInfo) bool {
	for _, f := range mi.Attributes {
		if strings.ToLower(f) == `\hasnochildren` {
			return true
		}
	}

	return false
}

func preparedName(miName, miDelim, delim, pathprfx string) string {
	s := strings.Replace(miName, miDelim, delim, -1)

	if pathprfx != "" {
		s = pathprfx + delim + s
	}

	return s
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
