package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/simia-tech/edkvs"
)

type options struct {
	ListenURL             string        `short:"l" long:"listen" default:"tcp://localhost:0" description:"listener address"`
	PeerURLs              []string      `short:"p" long:"peer" description:"address of target node. multiple specifications possible"`
	PeerReconnectInterval time.Duration `short:"r" long:"peer-reconnect-interval" default:"5s" description:"duration after which a failing peer is reconnected"`
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

	edkvs, err := edkvs.NewEDKVS(edkvs.Options{
		ListenURL: opts.ListenURL,
		PeerURLs:  opts.PeerURLs,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("node is listing at %s", edkvs.ListenURL())

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch

	if err := edkvs.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("node shut down")
}
