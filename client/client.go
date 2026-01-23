package client

import (
	"net/http"
	"sync"

	"github.com/ueebee/tachibanashi/auth"
	"github.com/ueebee/tachibanashi/event"
	"github.com/ueebee/tachibanashi/master"
	"github.com/ueebee/tachibanashi/price"
	"github.com/ueebee/tachibanashi/request"
)

type Client struct {
	cfg    Config
	http   *http.Client
	token  auth.TokenStore
	urlsMu sync.RWMutex
	urls   auth.VirtualURLs

	eventMu     sync.Mutex
	eventActive bool
	eventParams event.Params
	eventEno    int64
}

func New(cfg Config, opts ...Option) (*Client, error) {
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = BaseURLProd
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = DefaultUserAgent
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: cfg.Timeout}
	}
	if cfg.TokenStore == nil {
		cfg.TokenStore = auth.NewMemoryTokenStore()
	}

	return &Client{
		cfg:         cfg,
		http:        cfg.HTTPClient,
		token:       cfg.TokenStore,
		eventParams: cfg.EventParams,
	}, nil
}

func (c *Client) Auth() *auth.Service {
	return auth.NewService(c)
}

func (c *Client) Request() *request.Service {
	return request.NewService(c)
}

func (c *Client) Price() *price.Service {
	return price.NewService(c)
}

func (c *Client) Master() *master.Service {
	return master.NewService(c)
}

func (c *Client) Event() *event.Service {
	return event.NewService(c)
}

func (c *Client) BaseURL() string {
	return c.cfg.BaseURL
}

func (c *Client) TokenStore() auth.TokenStore {
	return c.token
}

func (c *Client) SetVirtualURLs(urls auth.VirtualURLs) {
	c.urlsMu.Lock()
	defer c.urlsMu.Unlock()
	c.urls = urls
}

func (c *Client) VirtualURLs() auth.VirtualURLs {
	c.urlsMu.RLock()
	defer c.urlsMu.RUnlock()
	return c.urls
}
