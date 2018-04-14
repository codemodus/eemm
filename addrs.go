package main

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	imap "github.com/emersion/go-imap"
)

const (
	msngMbx = "missing_mailbox"
	msngDom = "missing_domain"
	emailAt = "@"
)

type addresses []*imap.Address

func (ads addresses) appendBytesTo(b []byte) []byte {
	sort.Sort(&ads)

	for _, v := range ads {
		if addressIncomplete(v) || addressMailboxNameMalformed(v) {
			continue
		}

		b = append(b, []byte(v.MailboxName)...)
		b = append(b, []byte(v.HostName)...)
	}

	return b
}

func (ads *addresses) Len() int {
	return len(*ads)
}

func (ads *addresses) Swap(i, j int) {
	a := *ads
	a[i], a[j] = a[j], a[i]
}

func (ads *addresses) Less(i, j int) bool {
	a := *ads
	la, ra := a[i], a[j]

	lrc := utf8.RuneCountInString(la.MailboxName)
	rrc := utf8.RuneCountInString(ra.MailboxName)
	if lrc != rrc {
		return lrc < rrc
	}

	offset := 0
	for k, v := range la.MailboxName {
		r, w := utf8.DecodeRuneInString(ra.MailboxName[k-offset:])
		offset += utf8.RuneLen(v) - w

		if v != r {
			return v < r
		}
	}

	lrc = utf8.RuneCountInString(la.HostName)
	rrc = utf8.RuneCountInString(ra.HostName)
	if lrc != rrc {
		return lrc < rrc
	}

	offset = 0
	for k, v := range la.HostName {
		r, w := utf8.DecodeRuneInString(ra.HostName[k-offset:])
		offset += utf8.RuneLen(v) - w

		if v != r {
			return v < r
		}
	}

	return false
}

func addressIncomplete(a *imap.Address) bool {
	if a.MailboxName == "" || a.HostName == "" {
		return true
	}

	if strings.Contains(a.MailboxName, emailAt) || strings.Contains(a.HostName, emailAt) {
		return true
	}

	return strings.ToLower(a.MailboxName) == msngMbx || strings.ToLower(a.HostName) == msngDom
}

func addressMailboxNameMalformed(a *imap.Address) bool {
	rs := []rune(a.MailboxName)

	return len(rs) < 1 || !unicode.IsLetter(rs[0]) || !unicode.IsLetter(rs[len(rs)-1])
}
