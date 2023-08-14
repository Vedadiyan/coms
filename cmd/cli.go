package main

import (
	"strings"

	flaggy "github.com/vedadiyan/flaggy/pkg"
)

type Options struct {
	Port    int     `long:"--port" short:"-p" help:"Websocket port"`
	Cluster *string `long:"--cluster" short:"-c" help:"Cluster URL for solicited routes"`
	Routes  *string `long:"--routes" short:"-r" help:"Routes to solicit and connect"`
}

func (options Options) Run() error {
	if options.Port == 0 {
		flaggy.PrintHelp()
		return nil
	}
	var routes []string
	if options.Routes != nil {
		routes = strings.Split(*options.Routes, ",")
	}
	return Bootstrap(options.Port, options.Cluster, routes)
}
