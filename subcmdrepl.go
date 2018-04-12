package main

import (
	"sync"
)

type replicateArgs struct {
	cs      *coms
	l       *replScopedLog
	dstConf sessionConfig
	srcConf sessionConfig
}

type replicateCtx struct {
	t      *replTracker
	srvrID int
	accts  *replAccountsConf
	args   *replicateArgs
}

func generateReplicateCtxs(cs *coms, l Logger, t *replTracker, cnf replConf) (chan *replicateCtx, error) {
	c := make(chan *replicateCtx)

	go func() {
		defer close(c)

		for sID, g := range cnf.Servers {
			srvrCnf := g
			t.logServers(&srvrCnf)

			for asID, as := range g.Accounts {
				acctsCnf := as

				ctx := &replicateCtx{
					t:      t,
					srvrID: sID,
					accts:  &acctsCnf,
				}

				ctx.args = &replicateArgs{
					cs: cs,
					l:  newReplScopedLog(l, "REPL", sID, asID),
					dstConf: sessionConfig{
						server:   g.DstSrvrport[0],
						port:     g.DstSrvrport[1],
						pathprfx: g.DstPathprefix,
						account:  as.DstAcctpass[0],
						password: as.DstAcctpass[1],
					},
					srcConf: sessionConfig{
						server:   g.SrcSrvrport[0],
						port:     g.SrcSrvrport[1],
						pathprfx: g.SrcPathprefix,
						account:  as.SrcAcctpass[0],
						password: as.SrcAcctpass[1],
					},
				}

				c <- ctx
			}
		}
	}()

	return c, nil
}

func replicate(a *replicateArgs) error {
	a.l.Infof(
		"bonding to %s(%s) from %s(%s)",
		a.dstConf.account, a.dstConf.server, a.srcConf.account, a.srcConf.server,
	)

	bs, err := makeBondedSession(a.cs, a.dstConf, a.srcConf)
	if err != nil {
		a.l.Errorf("cannot bond sessions: %s", err)
		return err
	}
	defer func() {
		a.l.Info("closing bond")

		bs.close()
	}()

	a.l.Info("replicating mailboxes")

	if _, err := bs.replicateMailboxes(); err != nil {
		a.l.Errorf("cannot replicate mailboxes: %s", err)
		return err
	}

	a.l.Info("replicating messages")

	if err := bs.replicateMessages(); err != nil {
		a.l.Errorf("cannot replicate messages: %s", err)
		return err
	}

	return nil
}

func runReplication(cs *coms, l Logger, width int, cnf replConf) error {
	if !cnf.run {
		return nil
	}

	l.Info("replication started")
	defer func() { l.Info("replication ended") }()

	if err := cnf.Normalize(); err != nil {
		return err
	}

	t, err := newReplTracker()
	if err != nil {
		return err
	}
	defer func() {
		if derr := t.close(); derr != nil {
			l.Error(derr)
		}
	}()

	c, err := generateReplicateCtxs(cs, l, t, cnf)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < width; i++ {
		wg.Add(1)

		go func() {
			for ctx := range c {
				if err := replicate(ctx.args); err != nil {
					ctx.t.logInvalidAccts(ctx.srvrID, ctx.accts)
					continue
				}

				ctx.t.logValidAccts(ctx.srvrID, ctx.accts)
			}

			wg.Done()
		}()
	}
	wg.Wait()

	return nil
}
