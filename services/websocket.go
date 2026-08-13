package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WebSocketService manages live WebSocket connections. Each connection runs a
// read goroutine that streams inbound frames to the frontend via the single
// "ws:event" application event, keyed by a caller-supplied connID.
type WebSocketService struct {
	mu    sync.Mutex
	conns map[string]*wsConn
}

type wsConn struct {
	conn   *websocket.Conn
	cancel context.CancelFunc
}

// WsEvent is the payload emitted on the "ws:event" channel for every connection
// lifecycle change and inbound/outbound message.
type WsEvent struct {
	ConnID string `json:"connID"`
	// Type is one of: "open" | "message" | "sent" | "close" | "error".
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	Code int    `json:"code,omitempty"`
}

func (s *WebSocketService) emit(ev WsEvent) {
	if app := application.Get(); app != nil {
		app.Event.Emit("ws:event", ev)
	}
}

// Connect dials the WebSocket at url and begins streaming inbound frames. Any
// existing connection under the same connID is closed first. Optional
// subprotocols are comma-separated.
func (s *WebSocketService) Connect(connID, url string, headers []Pair, subprotocols string) error {
	s.mu.Lock()
	if s.conns == nil {
		s.conns = make(map[string]*wsConn)
	}
	if existing, ok := s.conns[connID]; ok {
		existing.cancel()
		existing.conn.Close(websocket.StatusNormalClosure, "reconnecting")
		delete(s.conns, connID)
	}
	s.mu.Unlock()

	hdr := http.Header{}
	for _, h := range headers {
		if h.Enabled && h.Name != "" {
			hdr.Set(h.Name, h.Value)
		}
	}

	var subs []string
	for _, p := range strings.Split(subprotocols, ",") {
		if p = strings.TrimSpace(p); p != "" {
			subs = append(subs, p)
		}
	}

	// Bounded dial; the read loop uses a separate long-lived context.
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer dialCancel()

	conn, _, err := websocket.Dial(dialCtx, url, &websocket.DialOptions{
		HTTPHeader:   hdr,
		Subprotocols: subs,
	})
	if err != nil {
		s.emit(WsEvent{ConnID: connID, Type: "error", Data: err.Error()})
		return err
	}
	conn.SetReadLimit(-1) // don't silently drop large frames in a dev tool

	ctx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.conns[connID] = &wsConn{conn: conn, cancel: cancel}
	s.mu.Unlock()

	s.emit(WsEvent{ConnID: connID, Type: "open"})
	go s.readLoop(ctx, connID, conn)
	return nil
}

func (s *WebSocketService) readLoop(ctx context.Context, connID string, conn *websocket.Conn) {
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			// A cancelled context means we initiated the close ourselves; the
			// Close/drop path already emitted the appropriate event.
			if ctx.Err() != nil {
				return
			}
			if code := websocket.CloseStatus(err); code != -1 {
				s.emit(WsEvent{ConnID: connID, Type: "close", Code: int(code), Data: err.Error()})
			} else {
				s.emit(WsEvent{ConnID: connID, Type: "error", Data: err.Error()})
			}
			s.drop(connID)
			return
		}
		s.emit(WsEvent{ConnID: connID, Type: "message", Data: string(data)})
	}
}

// Send writes a UTF-8 text frame to the connection.
func (s *WebSocketService) Send(connID, message string) error {
	s.mu.Lock()
	wc, ok := s.conns[connID]
	s.mu.Unlock()
	if !ok {
		return fmt.Errorf("no active connection: %s", connID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := wc.conn.Write(ctx, websocket.MessageText, []byte(message)); err != nil {
		s.emit(WsEvent{ConnID: connID, Type: "error", Data: err.Error()})
		return err
	}
	s.emit(WsEvent{ConnID: connID, Type: "sent", Data: message})
	return nil
}

// Close ends the connection with a normal-closure handshake.
func (s *WebSocketService) Close(connID string) error {
	s.mu.Lock()
	wc, ok := s.conns[connID]
	if ok {
		delete(s.conns, connID)
	}
	s.mu.Unlock()
	if !ok {
		return nil
	}
	wc.cancel()
	err := wc.conn.Close(websocket.StatusNormalClosure, "client closed")
	s.emit(WsEvent{ConnID: connID, Type: "close", Code: int(websocket.StatusNormalClosure)})
	return err
}

// drop tears down a connection whose read loop has already ended.
func (s *WebSocketService) drop(connID string) {
	s.mu.Lock()
	if wc, ok := s.conns[connID]; ok {
		wc.cancel()
		delete(s.conns, connID)
	}
	s.mu.Unlock()
}
