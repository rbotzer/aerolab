package main

import (
	"errors"
	"log"
)

type clusterDestroyCmd struct {
	clusterStopCmd
	Docker clusterDestroyCmdDocker `no-flag:"true"`
	clusterStartStopDestroyCmd
}

type clusterDestroyCmdDocker struct {
	Force bool `short:"f" long:"force" description:"force stop before destroy"`
}

func init() {
	addBackendSwitch("cluster.destroy", "docker", &a.opts.Cluster.Destroy.Docker)
}

func (c *clusterDestroyCmd) Execute(args []string) error {
	if earlyProcess(args) {
		return nil
	}
	log.Println("Running cluster.destroy")
	err := c.Nodes.ExpandNodes(string(c.ClusterName))
	if err != nil {
		return err
	}
	cList, nodes, err := c.getBasicData(string(c.ClusterName), c.Nodes.String())
	if err != nil {
		return err
	}
	var nerr error
	for _, ClusterName := range cList {
		if c.Docker.Force && a.opts.Config.Backend.Type == "docker" {
			b.ClusterStop(ClusterName, nodes[ClusterName])
		}
		err = b.ClusterDestroy(ClusterName, nodes[ClusterName])
		if err != nil {
			if nerr == nil {
				nerr = err
			} else {
				nerr = errors.New(nerr.Error() + "\n" + err.Error())
			}
		}
	}
	if nerr != nil {
		return nerr
	}
	log.Println("Done")
	return nil
}
