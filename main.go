package main

import (
	"os"

	"github.com/codemodus/sigmon"
	"github.com/sirupsen/logrus"
)

func main() {
	// wire into system signals and ignore during startup
	sm := sigmon.New(nil)
	sm.Run()

	// setup logging and main circuit breaker
	var l Logger = logrus.New()
	cs := newComs()
	trip := tripFn(cs, l)
	_ = trip

	// load config
	cnf, err := NewConf("config.cnf")
	if err != nil {
		if err == errFlagParse {
			os.Exit(1)
		}

		trip(err)
	}

	// act on config vals
	if !cnf.Main.verbose {
		l = &voidLog{}
	}
	width := runWidth(cnf.Main.rsrvd)

	// configure shutdown sequence
	sm.Set(func(s *sigmon.SignalMonitor) {
		cs.close()
	})

	trip(
		runReplication(cs, l, width, cnf.Repl),
	)

	// disconnect from system signals
	sm.Stop()

	// TODO: add flag(s) to restrict message handling to span (i.e. "after", "before")
	// TODO: add sub-command for duplicate removal
}
