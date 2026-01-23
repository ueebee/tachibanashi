package client

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	terrors "github.com/ueebee/tachibanashi/errors"
	"github.com/ueebee/tachibanashi/event"
	"nhooyr.io/websocket"
)

const (
	eventReconnectBase = time.Second
	eventReconnectMax  = 30 * time.Second
)

type wsConn struct {
	parent    *Client
	connMu    sync.Mutex
	conn      *websocket.Conn
	closeOnce sync.Once
	closed    chan struct{}
}

func (c *Client) DialEvent(ctx context.Context) (event.Conn, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	c.eventMu.Lock()
	if c.eventActive {
		c.eventMu.Unlock()
		return nil, errors.New("tachibanashi: event session already active")
	}
	c.eventActive = true
	c.eventMu.Unlock()

	ws := &wsConn{
		parent: c,
		closed: make(chan struct{}),
	}

	if err := ws.reconnect(ctx); err != nil {
		c.clearEventActive()
		return nil, err
	}

	return ws, nil
}

func (c *Client) clearEventActive() {
	c.eventMu.Lock()
	c.eventActive = false
	c.eventMu.Unlock()
}

func (c *Client) eventURL() (string, error) {
	urls := c.VirtualURLs()
	if urls.EventWS == "" {
		return "", errors.New("tachibanashi: virtual event websocket URL not set")
	}

	c.eventMu.Lock()
	params := c.eventParams
	lastEno := c.eventEno
	c.eventMu.Unlock()

	if lastEno > 0 && params.Eno == 0 {
		params.Eno = lastEno
	}

	return event.BuildWSURL(urls.EventWS, params)
}

func (c *Client) updateEventEno(value int64) {
	if value <= 0 {
		return
	}
	c.eventMu.Lock()
	if value > c.eventEno {
		c.eventEno = value
	}
	c.eventMu.Unlock()
}

func (c *wsConn) Recv(ctx context.Context) (event.Event, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	for {
		if c.isClosed() {
			return nil, errors.New("tachibanashi: event connection closed")
		}

		conn := c.current()
		if conn == nil {
			if err := c.reconnect(ctx); err != nil {
				return nil, err
			}
			continue
		}

		_, data, err := conn.Read(ctx)
		if err == nil {
			ev, err := event.DecodeEvent(data)
			if err != nil {
				if isUnsupportedCommand(err) {
					return event.Unknown{Raw: data}, nil
				}
				return nil, err
			}
			c.parent.updateEventEno(parseEventEno(ev))
			return ev, nil
		}

		if c.isClosed() {
			return nil, err
		}
		c.dropConn(conn)
		if err := c.reconnect(ctx); err != nil {
			return nil, err
		}
	}
}

func (c *wsConn) Close() error {
	var err error
	c.closeOnce.Do(func() {
		close(c.closed)
		conn := c.current()
		if conn != nil {
			err = conn.Close(websocket.StatusNormalClosure, "")
			c.dropConn(conn)
		}
		c.parent.clearEventActive()
	})
	return err
}

func (c *wsConn) reconnect(ctx context.Context) error {
	delay := eventReconnectBase
	for {
		if c.isClosed() {
			return errors.New("tachibanashi: event connection closed")
		}
		if err := ctx.Err(); err != nil {
			return err
		}

		url, err := c.parent.eventURL()
		if err != nil {
			return err
		}
		conn, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
			HTTPClient: c.parent.http,
		})
		if err == nil {
			c.setConn(conn)
			return nil
		}

		c.parent.logf("event reconnect failed: %v", err)
		if !sleep(ctx, delay) {
			return ctx.Err()
		}
		delay = minDuration(delay*2, eventReconnectMax)
	}
}

func (c *wsConn) current() *websocket.Conn {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	return c.conn
}

func (c *wsConn) setConn(conn *websocket.Conn) {
	c.connMu.Lock()
	c.conn = conn
	c.connMu.Unlock()
}

func (c *wsConn) dropConn(conn *websocket.Conn) {
	c.connMu.Lock()
	if c.conn == conn {
		c.conn = nil
	}
	c.connMu.Unlock()
}

func (c *wsConn) isClosed() bool {
	select {
	case <-c.closed:
		return true
	default:
		return false
	}
}

func (c *Client) logf(format string, args ...any) {
	if c.cfg.Logger == nil {
		return
	}
	c.cfg.Logger.Printf(format, args...)
}

func isUnsupportedCommand(err error) bool {
	var vErr *terrors.ValidationError
	if !errors.As(err, &vErr) {
		return false
	}
	if vErr.Field != "p_cmd" {
		return false
	}
	return strings.EqualFold(vErr.Reason, "unsupported")
}

func parseEventEno(ev event.Event) int64 {
	if ev == nil {
		return 0
	}
	if value, ok := ev.(interface{ Value(string) string }); ok {
		raw := strings.TrimSpace(value.Value("p_ENO"))
		if raw == "" {
			return 0
		}
		if parsed, err := parseInt64(raw); err == nil {
			return parsed
		}
	}
	return 0
}

func parseInt64(value string) (int64, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("tachibanashi: invalid integer %q", value)
	}
	return parsed, nil
}

func sleep(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
