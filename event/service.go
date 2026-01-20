package event

import (
	"context"

	terrors "github.com/ueebee/tachibanashi/errors"
)

type Dialer interface {
	DialEvent(ctx context.Context) (Conn, error)
}

type Conn interface {
	Recv(ctx context.Context) (Event, error)
	Close() error
}

type Event interface {
	Kind() string
}

type Unknown struct {
	Raw []byte
}

func (e Unknown) Kind() string {
	return "unknown"
}

type Service struct {
	dialer Dialer
}

func NewService(dialer Dialer) *Service {
	return &Service{dialer: dialer}
}

func (s *Service) Connect(ctx context.Context) (Conn, error) {
	if s.dialer == nil {
		return nil, terrors.ErrNotImplemented
	}
	return s.dialer.DialEvent(ctx)
}

func (s *Service) Stream(ctx context.Context) (<-chan Event, <-chan error) {
	events := make(chan Event)
	errs := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(errs)

		conn, err := s.Connect(ctx)
		if err != nil {
			errs <- err
			return
		}
		defer conn.Close()

		for {
			event, err := conn.Recv(ctx)
			if err != nil {
				errs <- err
				return
			}

			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case events <- event:
			}
		}
	}()

	return events, errs
}
