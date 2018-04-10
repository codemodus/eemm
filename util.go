package main

import (
	"errors"
	"os"
	"runtime"
)

var (
	errShutdown = errors.New("shutting down")
)

type coms struct {
	donec chan struct{}
}

func newComs() *coms {
	return &coms{
		donec: make(chan struct{}),
	}
}

func (c *coms) term() error {
	select {
	case <-c.donec:
		return errShutdown
	default:
		return nil
	}
}

func (c *coms) close() {
	select {
	case <-c.donec:
	default:
		close(c.donec)
	}
}

func tripFn(c *coms, l Logger) func(error) {
	return func(err error) {
		if err != nil {
			c.close()

			l.Errorf("TRIPPED: %s", err)
			l.Error("i'm melting! melting! oh, what a world! what a world!-")

			os.Exit(1)
		}
	}
}

func runWidth(reserved int) int {
	width := runtime.NumCPU() - reserved
	if width < 1 {
		width = 1
	}

	runtime.GOMAXPROCS(width)

	return width
}
