package main

import (
	"log"
	"os"
	"os/signal"

	flags "github.com/jessevdk/go-flags"
	"github.com/simia-tech/edkvs"
)

type options struct {
	Address       string   `short:"a" long:"address" description:"listener address"`
	NodeAddresses []string `short:"n" long:"node" description:"address of target node. multiple specifications possible"`
}

var (
	opts   options
	parser = flags.NewParser(&opts, flags.Default)
)

func main() {
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	store := edkvs.NewStore()
	node, err := edkvs.NewNode(store, "tcp", opts.Address)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("node is listing at %s", node.Addr())

	for _, nodeAddress := range opts.NodeAddresses {
		node.AddTarget("tcp", nodeAddress)
		log.Printf("node connected to %s", nodeAddress)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch

	if err := node.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("node teared down")
}
