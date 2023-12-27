package c2

import (
	"github.com/pygrum/Empress/transport"
	"net/http"
	"time"
)

const maxDuration time.Duration = 1<<63 - 1

type Client struct {
	Address      string
	HTTPAddress  string
	ClientInfo   *transport.Registration
	HttpClient   *http.Client
	Task         chan *transport.Request // A channel that receives tasks from the longPoll routine
	router       *Router
	cookieJar    []*http.Cookie // for whatever reason net/http/cookiejar was misbehaving
	nextResponse *transport.Response
}

func NewClient(addr, httpAddr string) (*Client, error) {
	return &Client{
		ClientInfo: &transport.Registration{},
		HttpClient: &http.Client{
			Timeout: maxDuration, // Polling client. Timeout has no effect if we get a connection refused
		},
		Address:     addr,
		HTTPAddress: httpAddr,
		Task:        make(chan *transport.Request),
	}, nil
}

func (c *Client) SetRouter(r *Router) {
	c.router = r
}
func (c *Client) Router() *Router {
	return c.router
}

func (c *Client) SetResponse(r *transport.Response) {
	c.nextResponse = r
}

func (c *Client) Response() *transport.Response {
	return c.nextResponse
}
