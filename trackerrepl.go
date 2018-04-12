package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

type replTracker struct {
	v    *os.File
	i    *os.File
	c    *os.File
	vMu  sync.Mutex
	vCnf *replConf
	iMu  sync.Mutex
	iCnf *replConf
	cMu  sync.Mutex
	cCnf *replConf
}

func newReplTracker() (*replTracker, error) {
	tm := time.Now()
	ft := tm.Format("200601021504")

	v, err := os.Create("valid_" + ft)
	if err != nil {
		return nil, err
	}

	i, err := os.Create("invalid_" + ft)
	if err != nil {
		return nil, err
	}

	c, err := os.Create("original_" + ft)
	if err != nil {
		return nil, err
	}

	t := &replTracker{
		v:    v,
		i:    i,
		c:    c,
		vCnf: &replConf{},
		iCnf: &replConf{},
		cCnf: &replConf{},
	}

	return t, nil
}

func (t *replTracker) logServers(srvrs *replServersConf) {
	addReplTrackedServers(&t.vMu, t.vCnf, srvrs)
	addReplTrackedServers(&t.iMu, t.iCnf, srvrs)
	addReplTrackedServers(&t.cMu, t.cCnf, srvrs)
}

func addReplTrackedServers(mu sync.Locker, cnf *replConf, srvrs *replServersConf) {
	s := *srvrs
	s.Accounts = nil
	mu.Lock()
	cnf.Servers = append(cnf.Servers, s)
	mu.Unlock()
}

func (t *replTracker) logValidAccts(srvrID int, accts *replAccountsConf) {
	addReplTrackedAccts(&t.vMu, t.vCnf, srvrID, accts)
	addReplTrackedAccts(&t.cMu, t.cCnf, srvrID, accts)
}

func (t *replTracker) logInvalidAccts(srvrID int, accts *replAccountsConf) {
	addReplTrackedAccts(&t.iMu, t.iCnf, srvrID, accts)
	addReplTrackedAccts(&t.cMu, t.cCnf, srvrID, accts)
}

func addReplTrackedAccts(mu sync.Locker, cnf *replConf, srvrID int, accts *replAccountsConf) {
	mu.Lock()
	cnf.Servers[srvrID].Accounts = append(cnf.Servers[srvrID].Accounts, *accts)
	mu.Unlock()
}

func (t *replTracker) close() error {
	var errs []error

	if err := closeReplTrackerFile(&t.vMu, t.v, t.vCnf); err != nil {
		errs = append(errs, err)
	}

	if err := closeReplTrackerFile(&t.iMu, t.i, t.iCnf); err != nil {
		errs = append(errs, err)
	}

	if err := closeReplTrackerFile(&t.cMu, t.c, t.cCnf); err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		return nil
	}

	e := ""
	for _, err := range errs {
		e += "& " + err.Error()
	}

	return fmt.Errorf(e)
}

func closeReplTrackerFile(mu sync.Locker, f *os.File, cnf *replConf) error {
	mu.Lock()
	defer mu.Unlock()

	defer func() { _ = f.Close() }()

	cnf = pruneEmptyReplTrackedServers(cnf)
	if len(cnf.Servers) == 0 {
		return nil
	}

	type replwrap struct {
		Repl replConf `toml:"Repl" json:"Repl"`
	}

	enc := toml.NewEncoder(f)
	return enc.Encode(replwrap{*cnf})
}

func pruneEmptyReplTrackedServers(cnf *replConf) *replConf {
	ss := cnf.Servers[:0]

	for _, s := range cnf.Servers {
		if s.Accounts != nil && len(s.Accounts) > 0 {
			ss = append(ss, s)
		}
	}

	cnf.Servers = ss

	return cnf
}
