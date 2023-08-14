package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/vedadiyan/coms/cluster/client"
	cluster "github.com/vedadiyan/coms/cluster/proto"
	"github.com/vedadiyan/coms/cluster/server"
	"github.com/vedadiyan/coms/cluster/state"
	"github.com/vedadiyan/coms/socket"
)

func Bootstrap(port int, clusterHost *string, clusterUrls []string) error {
	if clusterHost != nil {
		log.Println("creating cluster")
		selfHostValues := strings.Split(*clusterHost, ":")
		selfHostPort, err := strconv.Atoi(selfHostValues[1])
		if err != nil {
			return err
		}
		self := cluster.Node{}
		self.Id = state.GetId()
		self.Host = selfHostValues[0]
		self.Port = int32(selfHostPort)
		state.JoinSelf(&self)
		go server.New(*clusterHost)
		if len(clusterUrls) > 0 {
			log.Println("soliciting routes...")
			routes := cluster.NodeList{}
			routes.Id = state.GetId()
			routes.Nodes = append(routes.Nodes, &self)
			for _, host := range clusterUrls {
				hostValues := strings.Split(host, ":")
				port, err := strconv.Atoi(hostValues[1])
				if err != nil {
					return err
				}
				node := &cluster.Node{
					Host: hostValues[0],
					Port: int32(port),
				}
				client, err := client.New(node)
				if err != nil {
					log.Println(err)
					continue
				}
				id, err := client.GetId(context.TODO(), &cluster.Void{})
				if err != nil {
					log.Println(err)
					continue
				}
				node.Id = id.Id
				routes.Nodes = append(routes.Nodes, node)
			}
			server.Solicit(&routes)
		}
		log.Println("cluster started")
	}
	socket.New(fmt.Sprintf("127.0.0.1:%d", port), "/comms")
	return nil
}
