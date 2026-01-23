package client

import (
	"net/http"
	"time"

	"github.com/ueebee/tachibanashi/auth"
	"github.com/ueebee/tachibanashi/event"
)

const (
	BaseURLProd = "https://kabuka.e-shiten.jp/e_api_v4r8/"
	BaseURLDemo = "https://demo-kabuka.e-shiten.jp/e_api_v4r8/"

	DefaultTimeout   = 30 * time.Second
	DefaultUserAgent = "tachibanashi"
)

type Logger interface {
	Printf(format string, args ...any)
}

type Config struct {
	BaseURL     string
	Timeout     time.Duration
	UserAgent   string
	HTTPClient  *http.Client
	Logger      Logger
	TokenStore  auth.TokenStore
	EventParams event.Params
}
