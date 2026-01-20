package price

import "context"

type Client interface {
	DoJSON(ctx context.Context, method, path string, req, resp any) error
}

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{client: client}
}
