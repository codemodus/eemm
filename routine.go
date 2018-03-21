package main

import "fmt"

func replicate(cs *coms, l Logger) {
	// TODO: config = slice of dstCnf/srcCnf pairs
	dstConf := sessionConfig{
		server:   "mail.host.invalid",
		port:     "993",
		account:  "dst@example.com",
		password: "invalid",
	}

	srcConf := sessionConfig{
		server:   "mail.host.invalid",
		port:     "993",
		account:  "src@example.com",
		password: "invalid",
	}

	tl := makeTrackingLog(l, "BSS", 11)

	tl.logf("bonding to %s(%s) from %s(%s)", dstConf.account, dstConf.server, srcConf.account, srcConf.server)
	bs, err := makeBondedSession(cs, dstConf, srcConf)
	if err != nil {
		tl.logerr(fmt.Errorf("cannot bond sessions: %s", err))
		return
	}
	defer bs.close()

	tl.logf("replicating mailboxes")
	mis, err := bs.replicateMailboxes()
	if err != nil {
		tl.logerr(fmt.Errorf("cannot replicate mailboxes: %s", err))
		return
	}

	tl.logf("replicating messages")
	if err := bs.replicateMessages(mis); err != nil {
		tl.logerr(fmt.Errorf("cannot replicate messages: %s", err))
		return
	}
}
