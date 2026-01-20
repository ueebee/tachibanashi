package client

import (
	"context"

	terrors "github.com/ueebee/tachibanashi/errors"
	"github.com/ueebee/tachibanashi/event"
)

func (c *Client) DoJSON(ctx context.Context, method, path string, req, resp any) error {
	return terrors.ErrNotImplemented
}

func (c *Client) DialEvent(ctx context.Context) (event.Conn, error) {
	return nil, terrors.ErrNotImplemented
}
