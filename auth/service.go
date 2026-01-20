package auth

import (
	"context"

	terrors "github.com/ueebee/tachibanashi/errors"
)

type Client interface {
	DoJSON(ctx context.Context, method, path string, req, resp any) error
	TokenStore() TokenStore
}

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{client: client}
}

type Credentials struct {
	LoginID  string
	Password string
}

type VirtualURLs struct {
	Auth    string
	Request string
	Price   string
	Master  string
	Event   string
}

type LoginResponse struct {
	ResultCode string
	ResultText string
	URLs       VirtualURLs
}

func (s *Service) Login(ctx context.Context, creds Credentials) (*LoginResponse, error) {
	return nil, terrors.ErrNotImplemented
}

func (s *Service) Logout(ctx context.Context) error {
	return terrors.ErrNotImplemented
}

func (s *Service) VirtualURL(ctx context.Context) (*VirtualURLs, error) {
	return nil, terrors.ErrNotImplemented
}
