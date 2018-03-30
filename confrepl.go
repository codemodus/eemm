package main

import (
	"flag"
	"fmt"
)

type replAccountsConf struct {
	DstAcctpass []string
	SrcAcctpass []string
}

func (a *replAccountsConf) normalize(index, sub int) error {
	// if dst is not exist, is empty, or first index is empty, return error
	if len(a.DstAcctpass) == 0 || a.DstAcctpass[0] == "" {
		return fmt.Errorf("accounts %d-%d must define a destination", index, sub)
	}

	// if dst second index is not exist, set to empty
	if len(a.DstAcctpass) < 2 {
		a.DstAcctpass = []string{a.DstAcctpass[0], ""}
	}

	// if src is not exists, or is empty, inherit from dst
	if len(a.SrcAcctpass) == 0 {
		a.SrcAcctpass = a.DstAcctpass
	}

	// if src first index is empty, inherit from dst
	if a.SrcAcctpass[0] == "" {
		a.SrcAcctpass[0] = a.DstAcctpass[0]
	}

	// if src second index is not exist, or is empty, inherit from dst
	if len(a.SrcAcctpass) < 2 || a.SrcAcctpass[1] == "" {
		a.SrcAcctpass = []string{a.SrcAcctpass[0], a.DstAcctpass[1]}
	}

	return nil
}

type replServersConf struct {
	DstSrvrport   []string
	DstPathprefix string
	SrcSrvrport   []string
	SrcPathprefix string

	Accounts []replAccountsConf
}

func (g *replServersConf) normalize(index int) error {
	// if dst is not exist, is empty, or first index is empty, return error
	if len(g.DstSrvrport) == 0 || g.DstSrvrport[0] == "" {
		return fmt.Errorf("servers %d must define a destination", index)
	}

	// if dst second index is not exist, or is empty, set to default
	if len(g.DstSrvrport) < 2 || g.DstSrvrport[1] == "" {
		g.DstSrvrport = []string{g.DstSrvrport[0], "993"}
	}

	// if src is not exists, or is empty, inherit from dst
	if len(g.SrcSrvrport) == 0 {
		g.SrcSrvrport = g.DstSrvrport
	}

	// if src first index is empty, inherit from dst
	if g.SrcSrvrport[0] == "" {
		g.SrcSrvrport[0] = g.DstSrvrport[0]
	}

	// if src second index is not exist, or is empty, inherit from dst
	if len(g.SrcSrvrport) < 2 || g.SrcSrvrport[1] == "" {
		g.SrcSrvrport = []string{g.SrcSrvrport[0], g.DstSrvrport[1]}
	}

	return nil
}

type replConf struct {
	fs  *flag.FlagSet
	run bool

	Servers []replServersConf
}

func makeReplConf() replConf {
	return replConf{
		fs: flag.NewFlagSet("replicate", flag.ContinueOnError),
	}
}

func (c *replConf) AttachFlags() {}

func (c *replConf) Normalize() error {
	for k := range c.Servers {
		if err := c.Servers[k].normalize(k); err != nil {
			return err
		}

		for i := range c.Servers[k].Accounts {
			if err := c.Servers[k].Accounts[i].normalize(k, i); err != nil {
				return err
			}
		}
	}

	return nil
}
