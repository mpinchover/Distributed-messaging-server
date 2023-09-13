package requests

import "messaging-service/src/types/records"

type DeleteRoomEvent struct {
	EventType string `json:"eventType"`
	RoomUUID  string `json:"roomUuid"`
}

type LeaveRoomEvent struct {
	EventType string `json:"eventType"`
	RoomUUID  string `json:"roomUuid"`
	UserUUID  string `json:"userUuid"`
	Token     string `json:"token"`
}

// sennd to clients room has been opened
type OpenRoomEvent struct {
	EventType string        `json:"eventType"`
	Room      *records.Room `json:"room"`
}

// subscrve the sever to a room
type SubscribeToRoomEvent struct {
	EventType string   `json:"eventType"`
	Channel   string   `json:"channel"`
	Members   []string `json:"members"`
}

type SetClientConnectionEvent struct {
	EventType  string `json:"eventType"`
	UserUUID   string `json:"userUuid"`
	DeviceUUID string `json:"deviceUuid"`
	Token      string `json:"token"`
}

type TextMessageEvent struct {
	EventType  string           `json:"eventType"`
	FromUUID   string           `json:"fromUuid"`
	DeviceUUID string           `json:"deviceUuid"`
	Message    *records.Message `json:"message"`
	Token      string           `json:"token"`
}

type RoomsByUserUUIDEvent struct {
	EventType string          `json:"eventType"`
	UserUUID  string          `schema:"userUuid" validate:"required"`
	Offset    int             `schema:"offset"`
	Key       string          `schema:"key,-"`
	Rooms     []*records.Room `json:"rooms"`
	Token     string          `json:"token"`
}

type MessagesByRoomUUIDEvent struct {
	EventType string             `json:"eventType"`
	UserUUID  string             `schema:"userUuid" validate:"required"`
	RoomUUID  string             `schema:"roomUuid" validate:"required"`
	Offset    int                `schema:"offset"`
	Messages  []*records.Message `json:"messages"` // maybe make everything the actual record?
	Token     string             `json:"token"`
}

// the recpt has read the message
// client will have the user uuid stored. If the message is opened
// by not owner user uuid, send out the event
type SeenMessageEvent struct {
	EventType   string `json:"eventType"`
	MessageUUID string `json:"messageUuid"`
	UserUUID    string `json:"userUuid"`
	RoomUUID    string `json:"roomUuid"`
	Token       string `json:"token"`
}

type DeleteMessageEvent struct {
	EventType   string `json:"eventType"`
	MessageUUID string `json:"messageUuid"`
	UserUUID    string `json:"userUuid"`
	RoomUUID    string `json:"roomUuid"`
	Token       string `json:"token"`
}
