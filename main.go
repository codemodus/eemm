package main

import (
	"github.com/codemodus/sigmon"
	"github.com/sirupsen/logrus"
)

func main() {
	// wire into system signals and ignore during startup
	sm := sigmon.New(nil)
	sm.Run()

	// setup logging and main circuit breaker
	cs := newComs(logrus.New())
	trip := tripFn(cs)

	// configure shutdown sequence
	sm.Set(func(s *sigmon.SignalMonitor) {
		cs.close()
	})

	cs.Info("hello")
	// TODO: add sub-command for migration
	cs.Info("starting migration tool")

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

	// setup bonded session
	bs, err := newBondedSession(cs, 11, dstConf, srcConf)
	trip(err)
	defer bs.close()

	// run regularization procedure
	trip(bs.regularize())

	cs.Info("goodbye")

	// disconnect from system signals
	sm.Stop()
	// TODO: add flag to control concurrency
	// TODO: add flag(s) to restrict message handling to span (i.e. "after", "before")

	// TODO: add sub-command for duplicate removal
	// TODO: config = slice of dstCnf
}
