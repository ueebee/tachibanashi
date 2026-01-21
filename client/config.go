package client

import (
	"net/http"
	"time"

	"github.com/ueebee/tachibanashi/auth"
)

const (
	BaseURLProd = "https://kabuka.e-shiten.jp/e_api_v4r7/"
	BaseURLDemo = "https://10.62.26.91/e_api_v4r7/"

	DefaultTimeout   = 30 * time.Second
	DefaultUserAgent = "tachibanashi"
)

type Logger interface {
	Printf(format string, args ...any)
}

type Config struct {
	BaseURL    string
	Timeout    time.Duration
	UserAgent  string
	HTTPClient *http.Client
	Logger     Logger
	TokenStore auth.TokenStore
}
