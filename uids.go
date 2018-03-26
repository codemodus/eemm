package main

import (
	"sort"

	imap "github.com/emersion/go-imap"
)

type uniqIDs []uint32

func (ids uniqIDs) chunks(length int) []*imap.SeqSet {
	sort.Sort(&ids)

	var c []*imap.SeqSet

	for first, last := 0, length; first < len(ids); last += length {
		if last > len(ids) {
			last = len(ids)
		}

		s := &imap.SeqSet{}
		s.AddNum(ids[first:last]...)

		c = append(c, s)

		first = last
	}

	return c
}

func (ids *uniqIDs) Len() int {
	return len(*ids)
}

func (ids *uniqIDs) Swap(i, j int) {
	s := *ids
	s[i], s[j] = s[j], s[i]
}

func (ids *uniqIDs) Less(i, j int) bool {
	s := *ids
	return s[i] < s[j]
}
