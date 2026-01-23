package master

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/ueebee/tachibanashi/auth"
)

type Client interface {
	DoJSON(ctx context.Context, method, path string, req, resp any) error
	DoStream(ctx context.Context, method, path string, req any) (*http.Response, io.Reader, error)
	VirtualURLs() auth.VirtualURLs
}

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{client: client}
}

func (s *Service) masterURL() (string, error) {
	if s.client == nil {
		return "", errors.New("tachibanashi: master client is nil")
	}
	urls := s.client.VirtualURLs()
	if urls.Master == "" {
		return "", errors.New("tachibanashi: virtual master URL not set")
	}
	return urls.Master, nil
}
