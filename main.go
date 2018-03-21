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
	l := logrus.New()
	cs := newComs()
	trip := tripFn(cs, l)
	_ = trip

	// configure shutdown sequence
	sm.Set(func(s *sigmon.SignalMonitor) {
		cs.close()
	})

	l.Info("hello")
	// TODO: add sub-command for migration
	l.Info("starting migration tool")

	replicate(cs, l)

	l.Info("goodbye")

	// disconnect from system signals
	sm.Stop()
	// TODO: add flag to control concurrency
	// TODO: add flag(s) to restrict message handling to span (i.e. "after", "before")

	// TODO: add sub-command for duplicate removal
	// TODO: config = slice of dstCnf
}
