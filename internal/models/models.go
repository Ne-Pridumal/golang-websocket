package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Room struct {
	bun.BaseModel `bun:"table:rooms"`
	ID            int `bun:"id,pk,autoincrement"`
	Name          string
	Users         []User    `bun:"m2m:rooms_users,join:Room=User"`
	Messages      []Message `bun:"rel:has-many,join:id=room_id"`
}

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int `bun:"id,pk,autoincrement,alias:user_id"`
	Name          string
	Rooms         []Room `bun:"m2m:rooms_users,join:User=Room"`
}

type UserToRoom struct {
	bun.BaseModel `bun:"table:rooms_users"`
	RoomID        int   `bun:"room_id"`
	Room          *Room `bun:"rel:belongs-to,join:room_id=id"`
	UserID        int   `bun:"user_id"`
	User          *User `bun:"rel:belongs-to,join:user_id=id"`
}

type Message struct {
	bun.BaseModel `bun:"table:messages"`
	ID            int `bun:"id,pk,autoincrement"`
	RoomId        int `bun:"room_id,pk"`
	Content       string
	Date          time.Time
}
