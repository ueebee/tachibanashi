package client

import (
	"net/http"
	"time"

	"github.com/ueebee/tachibanashi/auth"
	"github.com/ueebee/tachibanashi/event"
)

type Option func(*Config)

func WithBaseURL(url string) Option {
	return func(c *Config) {
		c.BaseURL = url
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.Timeout = d
	}
}

func WithUserAgent(agent string) Option {
	return func(c *Config) {
		c.UserAgent = agent
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

func WithTokenStore(store auth.TokenStore) Option {
	return func(c *Config) {
		c.TokenStore = store
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

func WithEventParams(params event.Params) Option {
	return func(c *Config) {
		c.EventParams = params
	}
}
