package entities

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Message struct {
	EventType string `json:"eventType"`
}

type SetClientConnectionMessage struct {
	FromUUID string `json:"fromUuid"`
}

type Channel struct {
	Subscriber *redis.PubSub

	// just the participants that are on this server
	MembersOnServer map[string]bool
	UUID            string
}

type ChatRoom struct {
	UUID         string   `json:"uuid"`
	Participants []string `json:"participants"`
}

type Connection struct {
	Conn *websocket.Conn
	UUID string
}

// UnmarshalBinary decodes the struct into a User
func (m *Message) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}
	return nil
}

// MarshalBinary encodes the struct into a binary blob
func (m Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}
