package main

import "os"

type clusterAddCmd struct {
	Exporter clusterAddExporterCmd `command:"exporter" subcommands-optional:"true" description:"Install ams exporter in a cluster or clusters"`
	Help     helpCmd               `command:"help" subcommands-optional:"true" description:"Print help"`
}

func (c *clusterAddCmd) Execute(args []string) error {
	a.parser.WriteHelp(os.Stderr)
	os.Exit(1)
	return nil
}
