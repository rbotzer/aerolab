package main

import (
	"errors"

	"github.com/jessevdk/go-flags"
)

type clientCreateBaseCmd struct {
	ClientName    TypeClientName         `short:"n" long:"group-name" description:"Client group name" default:"client"`
	ClientCount   int                    `short:"c" long:"count" description:"Number of clients" default:"1"`
	NoSetHostname bool                   `short:"H" long:"no-set-hostname" description:"by default, hostname of each machine will be set, use this to prevent hostname change"`
	StartScript   flags.Filename         `short:"X" long:"start-script" description:"optionally specify a script to be installed which will run when the client machine starts"`
	Aws           clusterCreateCmdAws    `no-flag:"true"`
	Docker        clusterCreateCmdDocker `no-flag:"true"`
	osSelectorCmd
	Help helpCmd `command:"help" subcommands-optional:"true" description:"Print help"`
}

func (c *clientCreateBaseCmd) Execute(args []string) error {
	if earlyProcess(args) {
		return nil
	}
	_, err := c.createBase(args)
	return err
}

func (c *clientCreateBaseCmd) createBase(args []string) (machines []int, err error) {
	// TODO remember to work out if the command is GROW or CREATE and act accordingly
	b.WorkOnClients()
	return nil, errors.New("NOT IMPLEMENTED YET")
}