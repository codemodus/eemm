package main

import (
	"flag"
	"fmt"
	"sync"
)

type replConf struct {
	fs  *flag.FlagSet
	run bool

	Groups []struct {
		DstSrvrport [2]string
		SrcSrvrport [2]string

		Accounts []struct {
			DstAcctpass [2]string
			SrcAcctpass [2]string
		}
	}
}

func makeReplConf() replConf {
	return replConf{
		fs: flag.NewFlagSet("replicate", flag.ContinueOnError),
	}
}

func (c *replConf) AttachFlags() {}

type replicateArgs struct {
	cs               *coms
	l                Logger
	id, ct           int
	dstConf, srcConf sessionConfig
}

func runReplication(cs *coms, l Logger, width int, cnf replConf) error {
	if !cnf.run {
		return nil
	}

	l.Info("starting replication tool")

	c := make(chan *replicateArgs)
	go func() {
		defer close(c)

		for id, g := range cnf.Groups {
			for ct, as := range g.Accounts {
				ct++

				a := &replicateArgs{
					cs: cs,
					l:  l,
					id: id,
					ct: ct,
					dstConf: sessionConfig{
						server:   g.DstSrvrport[0],
						port:     g.DstSrvrport[1],
						account:  as.DstAcctpass[0],
						password: as.DstAcctpass[1],
					},
					srcConf: sessionConfig{
						server:   g.SrcSrvrport[0],
						port:     g.SrcSrvrport[1],
						account:  as.SrcAcctpass[0],
						password: as.SrcAcctpass[1],
					},
				}

				c <- a
			}
		}
	}()

	wg := &sync.WaitGroup{}

	for i := 0; i < width; i++ {
		wg.Add(1)

		go func(id int) {
			for a := range c {
				replicate(a)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	return nil
}

func replicate(a *replicateArgs) {
	tl := makeTrackingLog(a.l, "REPL", a.id, a.ct)

	tl.logf(
		"bonding to %s(%s) from %s(%s)",
		a.dstConf.account, a.dstConf.server, a.srcConf.account, a.srcConf.server,
	)

	bs, err := makeBondedSession(a.cs, a.dstConf, a.srcConf)
	if err != nil {
		tl.logerr(fmt.Errorf("cannot bond sessions: %s", err))
		return
	}
	defer bs.close()

	tl.logf("replicating mailboxes")

	if _, err := bs.replicateMailboxes(); err != nil {
		tl.logerr(fmt.Errorf("cannot replicate mailboxes: %s", err))
		return
	}

	tl.logf("replicating messages")

	if err := bs.replicateMessages(); err != nil {
		tl.logerr(fmt.Errorf("cannot replicate messages: %s", err))
		return
	}
}
