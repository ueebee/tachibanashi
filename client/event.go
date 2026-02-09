package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	terrors "github.com/ueebee/tachibanashi/errors"
	"github.com/ueebee/tachibanashi/event"
)

const (
	eventReconnectBase = time.Second
	eventReconnectMax  = 30 * time.Second
	eventReadTimeout   = 60 * time.Second
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

		data, err := c.readMessage(ctx, conn)
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
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			err = conn.Close()
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
		dialer := c.parent.wsDialer(url)
		conn, _, err := dialer.DialContext(ctx, url, nil)
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

func (c *wsConn) readMessage(ctx context.Context, conn *websocket.Conn) ([]byte, error) {
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if deadline, ok := ctx.Deadline(); ok {
			_ = conn.SetReadDeadline(deadline)
		} else {
			_ = conn.SetReadDeadline(time.Now().Add(eventReadTimeout))
		}

		messageType, data, err := conn.ReadMessage()
		if err != nil {
			return nil, err
		}
		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}
		return data, nil
	}
}

func (c *Client) logf(format string, args ...any) {
	if c.cfg.Logger == nil {
		return
	}
	c.cfg.Logger.Printf(format, args...)
}

func (c *Client) wsDialer(targetURL string) *websocket.Dialer {
	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = c.cfg.Timeout
	if transport, ok := c.http.Transport.(*http.Transport); ok && transport != nil {
		dialer.TLSClientConfig = transport.TLSClientConfig
		dialer.Proxy = transport.Proxy
		if transport.DialContext != nil {
			dialer.NetDialContext = transport.DialContext
		}
		if transport.DialTLSContext != nil {
			dialer.NetDialTLSContext = transport.DialTLSContext
		}
	}
	if dialer.TLSClientConfig == nil {
		// Use the target WebSocket URL's hostname for TLS ServerName,
		// not the base URL. This is important when eventWS uses a different
		// host than the base API URL (e.g., price-kabuka.e-shiten.jp vs kabuka.e-shiten.jp)
		wsURL := strings.TrimSpace(targetURL)
		if wsURL != "" {
			if parsed, err := url.Parse(wsURL); err == nil && (parsed.Scheme == "wss" || parsed.Scheme == "https") {
				dialer.TLSClientConfig = &tls.Config{ServerName: parsed.Hostname()}
			}
		}
	}
	return &dialer
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
