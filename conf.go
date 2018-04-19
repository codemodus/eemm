package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/codemodus/config"
)

var (
	errFlagParse = errors.New("failed to parse flags")
)

type mainConf struct {
	fs      *flag.FlagSet
	verbose bool
	rsrvd   int
	conc    int
}

func makeMainConf() mainConf {
	return mainConf{
		fs:    flag.NewFlagSet("main", flag.ContinueOnError),
		rsrvd: 2,
		conc:  8,
	}
}

func (c *mainConf) AttachFlags() {
	c.fs.BoolVar(&c.verbose, "v", c.verbose, "enable logging")
	c.fs.IntVar(&c.rsrvd, "rcpus", c.rsrvd, "reserved cpu count")
	c.fs.IntVar(&c.conc, "conc", c.conc, "concurrency count; limited to unreserved cpus")
}

func (c *mainConf) Normalize() error { return nil }

// Conf ...
type Conf struct {
	sync.Mutex
	Main mainConf `toml:"Main" json:"Main"`
	Repl replConf `toml:"Repl" json:"Repl"`
}

// NewConf ...
func NewConf(fpath string) (*Conf, error) {
	c := &Conf{
		Main: makeMainConf(),
		Repl: makeReplConf(),
	}

	if err := config.Init(c, fpath); err != nil {
		return nil, err
	}

	return c, nil
}

// InitPost ...
func (c *Conf) InitPost() error {
	c.Main.AttachFlags()
	c.Repl.AttachFlags()

	if err := c.Main.fs.Parse(os.Args[1:]); err != nil {
		return errFlagParse
	}

	if len(c.Main.fs.Args()) == 0 {
		return nil
	}

	switch cmd := c.Main.fs.Args()[0]; cmd {
	case c.Repl.fs.Name():
		if err := c.Repl.fs.Parse(nextArgs(os.Args, cmd)); err != nil {
			return errFlagParse
		}

		c.Repl.run = true

	default:
		fmt.Fprintf(
			c.Main.fs.Output(),
			"%q is not a valid subcommand: [%s]\n",
			cmd, c.Repl.fs.Name(),
		)

		return errFlagParse

	}

	return c.Main.Normalize()
}

func nextArgs(vals []string, val string) []string {
	for k, v := range vals {
		if v == val {
			return vals[k+1:]
		}
	}

	return vals
}
