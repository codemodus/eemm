package main

import (
	"errors"
	"os"
)

// Logger describes basic logging functions.
type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

type voidLog struct{}

func (l *voidLog) Info(...interface{})           {}
func (l *voidLog) Infof(string, ...interface{})  {}
func (l *voidLog) Error(...interface{})          {}
func (l *voidLog) Errorf(string, ...interface{}) {}

type coms struct {
	Logger
	ErrShutdown error

	donec chan struct{}
}

func newComs(log Logger) *coms {
	if log == nil {
		log = &voidLog{}
	}

	return &coms{
		Logger:      log,
		ErrShutdown: errors.New("shutting down"),
		donec:       make(chan struct{}),
	}
}

func (c *coms) done() <-chan struct{} {
	return c.donec
}

func (c *coms) close() {
	select {
	case <-c.donec:
	default:
		close(c.donec)
	}
}

func tripFn(cs *coms) func(error) {
	return func(err error) {
		if err != nil {
			cs.Errorf("TRIPPED: %s", err)
			cs.Error("i'm melting! melting! oh, what a world! what a world!-")

			os.Exit(1)
		}
	}
}
