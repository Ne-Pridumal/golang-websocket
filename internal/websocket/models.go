package ws

import (
	"nhooyr.io/websocket"
)

type User struct {
	Username  string
	Msgs      chan []byte
	Conn      *websocket.Conn
	CloseSlow func()
}

type Message struct {
	Type    string `json:"type"`
	User    string `json:"user,omitempty"`
	Content string `json:"content"`
}

type ContactList struct {
	Username     string `json:"username"`
	LastActivity int64  `json:"last_activity"`
}
