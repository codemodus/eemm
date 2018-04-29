package main

import "testing"

func TestPathNameMatches(t *testing.T) {
	ds := []struct {
		name string
		ref  string
		val  string
		want bool
	}{
		{"static inequal short", "INBOX/test", "oops/", false},
		{"static inequal", "INBOX/test", "INBOX/oops", false},
		{"static inequal long", "INBOX/test", "INBOX/testx", false},
		{"static equal", "INBOX/test", "INBOX/test", true},
		{"wildcard nonequiv prefix", "INBOX/*", "INBUZ/oops", false},
		{"wildcard nonequiv long", "INBOX/*", "INBOX/something/test", false},
		{"wildcard equiv", "INBOX/*", "INBOX/test", true},
		{"wildcard equiv long", "INBOX/*/test", "INBOX/something/test", true},
	}

	for _, d := range ds {
		got := pathNameMatches(d.ref, d.val)
		if got != d.want {
			t.Errorf("%s: got %v, want %v", d.name, got, d.want)
		}
	}
}
