package main

type clientCmd struct {
	Create  clientCreateCmd  `command:"create" subcommands-optional:"true" description:"Create new client machines"`
	Add     clientAddCmd     `command:"add" subcommands-optional:"true" description:"Add features to existing client machines"`
	List    clientListCmd    `command:"list" subcommands-optional:"true" description:"List client machine groups"`
	Start   clientStartCmd   `command:"start" subcommands-optional:"true" description:"Start a client machine group"`
	Stop    clientStopCmd    `command:"stop" subcommands-optional:"true" description:"Stop a client machine group"`
	Grow    clientGrowCmd    `command:"grow" subcommands-optional:"true" description:"Grow a client machine group"`
	Destroy clientDestroyCmd `command:"destroy" subcommands-optional:"true" description:"Destroy client(s)"`
	Help    helpCmd          `command:"help" subcommands-optional:"true" description:"Print help"`
}