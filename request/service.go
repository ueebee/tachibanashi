package request

import (
	"context"
	"errors"

	"github.com/ueebee/tachibanashi/auth"
)

type Client interface {
	DoJSON(ctx context.Context, method, path string, req, resp any) error
	VirtualURLs() auth.VirtualURLs
}

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{client: client}
}

func (s *Service) requestURL() (string, error) {
	if s.client == nil {
		return "", errors.New("tachibanashi: request client is nil")
	}
	urls := s.client.VirtualURLs()
	if urls.Request == "" {
		return "", errors.New("tachibanashi: virtual request URL not set")
	}
	return urls.Request, nil
}
