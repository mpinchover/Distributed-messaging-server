package entities

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Message struct {
	EventType string `json:"eventType"`
}

type ChatMessageEvent struct {
	FromUserUUID       *string `json:"fromUuid"`
	FromConnectionUUID *string `json:"fromConnectionUuid"`
	RoomUUID           *string `json:"roomUuid"`
	MessageText        *string `json:"messageText"`
	EventType          *string `json:"eventType"`
}

type OpenRoomRequestMessage struct {
	FromUUID *string `json:"fromUuid"`
	ToUUID   *string `json:"toUUID"`
}

type OpenRoomEvent struct {
	EventType          *string   `json:"eventType"`
	FromUUID           *string   `json:"fromUuid"`
	FromConnectionUUID *string   `json:"fromConnectionUuid"`
	ToUUID             *string   `json:"toUuid"`
	Room               *ChatRoom `json:"room"`
}

type SetClientConnectionEvent struct {
	EventType      *string `json:"eventType"`
	FromUUID       *string `json:"fromUuid"`
	ConnectionUUID *string `json:"connectionUuid"`
}

type OpenRoomRequest struct {
	FromUUID *string `json:"fromUuid"`
	ToUUID   *string `json:"toUuid"`
}

type SetClientConnectionMessage struct {
	FromUUID *string `json:"fromUuid"`
}

type Channel struct {
	Subscriber *redis.PubSub

	// just the participants that are on this server
	ParticipantsOnServer map[string]bool
	UUID                 *string
}

type ChatRoom struct {
	UUID         *string  `json:"uuid"`
	Participants []string `json:"participants"`
}

type Connection struct {
	Conn *websocket.Conn
	UUID *string
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
