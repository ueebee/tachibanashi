package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestReadMessageTimeoutReturns(t *testing.T) {
	done := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		<-done
	}))
	defer server.Close()
	defer close(done)

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	errCh := make(chan error, 1)
	panicCh := make(chan any, 1)
	c := &wsConn{}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicCh <- r
				return
			}
		}()
		_, err := c.readMessage(context.Background(), conn)
		errCh <- err
	}()

	select {
	case r := <-panicCh:
		t.Fatalf("readMessage panicked: %v", r)
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected timeout error")
		}
	case <-time.After(7 * time.Second):
		t.Fatal("readMessage did not return within timeout")
	}
}
