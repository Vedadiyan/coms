package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vedadiyan/coms/cluster/client"
	cluster "github.com/vedadiyan/coms/cluster/proto"
	"github.com/vedadiyan/coms/cluster/server"
	"github.com/vedadiyan/coms/cluster/state"
)

func main() {
	fmt.Println("my id:", state.GetId())
	selfHost := os.Args[1]
	selfHostValues := strings.Split(selfHost, ":")
	selfHostPort, err := strconv.Atoi(selfHostValues[1])
	if err != nil {
		panic(err)
	}
	self := cluster.Node{}
	self.Id = state.GetId()
	self.Host = selfHostValues[0]
	self.Port = int32(selfHostPort)
	state.JoinSelf(&self)
	go server.New(os.Args[1])
	if len(os.Args) > 2 {
		log.Println("routing...")
		routes := cluster.NodeList{}
		routes.Id = state.GetId()
		routes.Nodes = append(routes.Nodes, &self)
		for i := 2; i < len(os.Args); i++ {
			host := os.Args[i]
			hostValues := strings.Split(host, ":")
			port, err := strconv.Atoi(hostValues[1])
			if err != nil {
				log.Println(err)
				continue
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
		server.Route(&routes)
	}
	fmt.Scanln()
}
