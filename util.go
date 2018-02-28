package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type coms struct {
	ic chan string
	ec chan error
	dc chan struct{}
	wg *sync.WaitGroup
}

func newComs() *coms {
	return &coms{
		ic: make(chan string),
		ec: make(chan error),
		dc: make(chan struct{}),
		wg: &sync.WaitGroup{},
	}
}

func tripFn(cs *coms) func(error) {
	return func(err error) {
		if err != nil {
			cs.ec <- fmt.Errorf("TRIPPED: %s", err)
			cs.ec <- fmt.Errorf("i'm melting! melting! oh, what a world! what a world!")
			select {
			case <-cs.dc:
			default:
				close(cs.dc)
			}

			cs.wg.Wait()

			os.Exit(1)
		}
	}
}

func log(cs *coms) {
	l := logrus.New()

	cs.wg.Add(1)

	go func() {
		defer cs.wg.Done()

		for {
			select {
			case i := <-cs.ic:
				l.Info(i)
			case e := <-cs.ec:
				l.Error(e)
			case <-cs.dc:
				return
			}
		}
	}()
}
