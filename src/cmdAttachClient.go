package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

type attachClientCmd struct {
	Trino      attachCmdTrino `command:"trino" subcommands-optional:"true" description:"Attach to trino shell"`
	ClientName TypeClientName `short:"n" long:"name" description:"Client group name" default:"client"`
	Machine    TypeMachines   `short:"l" long:"node" description:"Machine to attach to (or comma-separated list, when using '-- ...'). Example: 'attach shell --node=all -- /some/command' will execute command on all nodes" default:"1"`
	Detach     bool           `long:"detach" description:"detach the process stdin - will not kill process on CTRL+C, disables parallel"`
	Parallel   bool           `long:"parallel" description:"enable parallel execution across all machines"`
	Tail       []string       `description:"List containing command parameters to execute, ex: [\"ls\",\"/opt\"]"`
	Help       attachCmdHelp  `command:"help" subcommands-optional:"true" description:"Print help"`
}

func (c *attachClientCmd) Execute(args []string) error {
	if earlyProcess(args) {
		return nil
	}
	return c.run(args)
}

func (c *attachClientCmd) run(args []string) (err error) {
	b.WorkOnClients()
	var nodes []int
	err = c.Machine.ExpandNodes(string(c.ClientName))
	if err != nil {
		return err
	}
	if c.Machine == "all" {
		nodes, err = b.NodeListInCluster(string(c.ClientName))
		if err != nil {
			return err
		}
	} else {
		for _, node := range strings.Split(c.Machine.String(), ",") {
			nodeInt, err := strconv.Atoi(node)
			if err != nil {
				return err
			}
			nodes = append(nodes, nodeInt)
		}
	}
	if len(nodes) > 1 && len(args) == 0 {
		return fmt.Errorf("%s", "When using more than 1 node in node-attach, you must specify the command to run. For example: 'node-attach -l 1,2,3 -- /command/to/run'")
	}

	if c.Detach {
		out, err := b.RunCommands(c.ClientName.String(), [][]string{args}, nodes)
		if err != nil {
			log.Print(err)
		}
		for _, o := range out {
			fmt.Println(string(o))
		}
	}

	if !c.Parallel {
		for _, node := range nodes {
			if len(nodes) > 1 {
				fmt.Printf(" ======== %s:%d ========\n", string(c.ClientName), node)
			}
			erra := b.AttachAndRun(string(c.ClientName), node, args)
			if erra != nil {
				if err == nil {
					err = erra
				} else {
					err = fmt.Errorf("%s\n%s", err.Error(), erra.Error())
				}
			}
		}
		return err
	}

	wg := new(sync.WaitGroup)
	for _, node := range nodes {
		wg.Add(1)
		go c.runbg(wg, node, args)
	}
	wg.Wait()
	return nil
}

func (c *attachClientCmd) runbg(wg *sync.WaitGroup, node int, args []string) {
	defer wg.Done()
	err := b.AttachAndRun(string(c.ClientName), node, args)
	if err != nil {
		log.Printf(" ---- Node %d ERROR: %s", node, err)
	}
}
