package main

import (
	"fmt"
	"sort"
)

type clusterListCmd struct {
	Json bool    `short:"j" long:"json" description:"Provide output in json format"`
	IP   bool    `short:"i" long:"ip" description:"print only the IP of the client machines (disables JSON output)"`
	Help helpCmd `command:"help" subcommands-optional:"true" description:"Print help"`
}

func (c *clusterListCmd) Execute(args []string) error {
	if earlyProcess(args) {
		return nil
	}
	if c.IP {
		clusters, err := b.ClusterList()
		if err != nil {
			return err
		}
		sort.Strings(clusters)
		for _, cluster := range clusters {
			nodesI, err := b.GetNodeIpMap(cluster, true)
			if err != nil {
				return err
			}
			nodesE, err := b.GetNodeIpMap(cluster, false)
			if err != nil {
				return err
			}
			nodesIS := make([]int, 0, len(nodesI))
			for k := range nodesI {
				nodesIS = append(nodesIS, k)
			}
			sort.Ints(nodesIS)
			nodesES := make([]int, 0, len(nodesE))
			for k := range nodesE {
				nodesES = append(nodesES, k)
			}
			sort.Ints(nodesES)
			for _, no := range nodesIS {
				ip := nodesI[no]
				extIP := nodesE[no]
				fmt.Printf("cluster=%s node=%d int_ip=%s ext_ip=%s\n", cluster, no, ip, extIP)
			}
			for _, no := range nodesES {
				ip := nodesE[no]
				if nodesI[no] == "" {
					fmt.Printf("cluster=%s node=%d int_ip= ext_ip=%s\n", cluster, no, ip)
				}
			}
		}
		return nil
	}
	f, e := b.ClusterListFull(c.Json)
	if e != nil {
		return e
	}
	fmt.Println(f)
	return nil
}
