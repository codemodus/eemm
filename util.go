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

func runWidth(reserved, concurrency int) int {
	maxProcs := runtime.NumCPU() - reserved
	if maxProcs < 1 {
		maxProcs = 1
	}

	runtime.GOMAXPROCS(maxProcs)

	if concurrency > maxProcs {
		concurrency = maxProcs
	}

	return concurrency
}

func pathNameMatches(ref, val string) bool {
	if len(val) < len(ref) {
		return false
	}

	n := -1
	for i := 0; i < len(ref); i++ {
		n++
		if n >= len(val)-1 {
			break
		}

		if ref[i] == '*' {
			if val[n] == '/' {
				n--
				continue
			}

			i--
			continue
		}

		if ref[i] != val[n] {
			return false
		}
	}

	return n == len(val)-1
}
