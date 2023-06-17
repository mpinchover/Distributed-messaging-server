package events

import "messaging-service/types/records"

type DeleteChatRoomEvent struct {
	EventType string `json:"eventType"`
	RoomUUID  string `json:"roomUuid"`
}

type OpenRoomEvent struct {
	EventType          string            `json:"eventType"`
	FromUUID           string            `json:"fromUuid"`
	FromConnectionUUID string            `json:"fromConnectionUuid"`
	ToUUID             string            `json:"toUuid"`
	Room               *records.ChatRoom `json:"room"`
}

type SetClientConnectionEvent struct {
	EventType      string `json:"eventType"`
	FromUUID       string `json:"fromUuid"`
	ConnectionUUID string `json:"connectionUuid"`
}

type ChatMessageEvent struct {
	FromUserUUID       string `json:"fromUuid"`
	FromConnectionUUID string `json:"fromConnectionUuid"`
	RoomUUID           string `json:"roomUuid"`
	MessageText        string `json:"messageText"`
	EventType          string `json:"eventType"`
}
