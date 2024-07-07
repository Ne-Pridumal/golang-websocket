package ws

import (
	"context"
	"errors"
	sl "golang-websocket-chat/lib/logger/slog"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"nhooyr.io/websocket"
)

type WsServer struct {
	// subscriberMessageBuffer controls the max number
	// of messages that can be queued for a subscriber
	// before it is kicked.
	// Defaults to 16.
	subscriberMessageBuffer int
	// logf controls where logs are sent.
	// Defaults to log.Printf.
	logger        *slog.Logger
	serverMux     *chi.Mux
	subscribersMu sync.Mutex
	subscribers   map[*subscriber]struct{}
}

func New(log *slog.Logger) *WsServer {
	const op = "ws.New"

	router := chi.NewRouter()

	s := &WsServer{
		logger:                  log,
		subscriberMessageBuffer: 16,
		serverMux:               router,
	}

	return s
}

func (s *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serverMux.ServeHTTP(w, r)
}

func (s *WsServer) configureRoutes() {
	s.serverMux.Group(func(r chi.Router) {
		s.serverMux.HandleFunc("/", s.subscribeHandler)
	})
}

func (s *WsServer) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	const op = "ws.subscribeHandler"
	err := s.subscribe(r.Context(), w, r)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		s.logger.Error(op, sl.Err(err))
		return
	}
}

func (s *WsServer) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var mu sync.Mutex
	var c *websocket.Conn
	var closed bool
	sub := &subscriber{
		msgs: make(chan []byte, s.subscriberMessageBuffer),
		closeSlow: func() {
			mu.Lock()
			defer mu.Unlock()
			closed = true
			if c != nil {
				c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
			}
		},
	}
	s.addSubscriber(sub)
	defer s.deleteSubscriber(sub)

	c2, err := websocket.Accept(w, r, nil)

	if err != nil {
		return err
	}
	mu.Lock()
	if closed {
		mu.Unlock()
		return net.ErrClosed
	}

	c = c2
	mu.Unlock()
	defer c.CloseNow()

	ctx = c.CloseRead(ctx)

	for {
		select {
		case msg := <-sub.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *WsServer) sendMessage() {

}

func (s *WsServer) addSubscriber(sub *subscriber) {
	s.subscribersMu.Lock()
	s.subscribers[sub] = struct{}{}
	s.subscribersMu.Unlock()
}

func (s *WsServer) deleteSubscriber(sub *subscriber) {
	s.subscribersMu.Lock()
	delete(s.subscribers, sub)
	s.subscribersMu.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
