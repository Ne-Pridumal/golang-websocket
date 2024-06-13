package ws

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

type ChatRoom struct {
	ID       string
	Name     string
	Users    map[*User]struct{}
	Messages map[*Message]struct{}
	roomMu   sync.Mutex
}

func NewChatRoom(id string, name string) ChatRoom {
	return ChatRoom{
		ID:   id,
		Name: name,
	}
}

func (s *ChatRoom) RoomSubscribeHandler(w http.ResponseWriter, r *http.Request) {
	err := s.RoomSubscribe(r.Context(), w, r)
	if errors.Is(err, context.Canceled) {
		log.Println(err)
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
}

func (cr *ChatRoom) RoomSubscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var mu sync.Mutex
	var c *websocket.Conn
	var closed bool
	user := &User{
		Msgs: make(chan []byte, 1024),
		Conn: c,
		CloseSlow: func() {
			mu.Lock()
			defer mu.Unlock()
			closed = true
			if c != nil {
				c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
			}
		},
	}
	cr.addUser(user)
	defer cr.removeUser(user)

	c2, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}
	mu.Lock()
	if closed {
		mu.Unlock()
		return net.ErrClosed
	}
	c = c2
	user.Conn = c
	mu.Unlock()
	defer c.CloseNow()

	for {
		_, m, err := c.Read(ctx)
		if err != nil {
			return err
		}
		cr.sendMessage(ctx, m)
	}
}

func (cr *ChatRoom) addUser(user *User) {
	cr.roomMu.Lock()
	defer cr.roomMu.Unlock()
	cr.Users[user] = struct{}{}
}

func (cr *ChatRoom) removeUser(user *User) {
	cr.roomMu.Lock()
	defer cr.roomMu.Unlock()
	delete(cr.Users, user)
}

func (cr *ChatRoom) writeMessage(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Write(ctx, websocket.MessageText, msg)
}

func (cr *ChatRoom) sendMessage(ctx context.Context, msg []byte) {
	cr.roomMu.Lock()
	defer cr.roomMu.Unlock()
	for u := range cr.Users {
		select {
		case u.Msgs <- msg:
			err := cr.writeMessage(ctx, time.Second*5, u.Conn, <-u.Msgs)
			if err != nil {
				go u.CloseSlow()
			}
		default:
			go u.CloseSlow()
		}
	}
}
