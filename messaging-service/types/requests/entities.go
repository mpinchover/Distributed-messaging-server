package requests

import (
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Member struct {
	UUID     string `json:"uuid"`
	UserUUID string `json:"userUuid"`
	UserRole string `json:"user_role"`
}

type Room struct {
	UUID     string     `json:"uuid"`
	Members  []*Member  `json:"members"`
	Messages []*Message `json:"messages"`
}

type Message struct {
	UUID        string `json:"uuid"`
	FromUUID    string `json:"fromUuid"`
	RoomUUID    string `json:"roomUuid"`
	MessageText string `json:"messageText"`
	CreatedAt   int64  `json:"createdAt"`
}

type Connection struct {
	UserUUID    string
	Connections map[string]*websocket.Conn
}

// server specific room
type ServerChannel struct {
	Subscriber *redis.PubSub

	// just the participants that are on this server
	MembersOnServer map[string]bool
	UUID            string
}
