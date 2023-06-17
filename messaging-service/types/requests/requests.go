package requests

import (
	"messaging-service/types/records"
)

type GetRoomsByUserUUIDRequest struct {
	UserUUID string `schema:"userUuid"`
	Offset   int    `schema:"offset"`
}

type GetRoomsByUserUUIDResponse struct {
	Rooms []*records.ChatRoom `json:"rooms"`
}

type GetMessagesByRoomUUIDRequest struct {
	RoomUUID string `schema:"roomUuid"`
	Offset   int    `schema:"offset"`
}

type GetMessagesByRoomUUIDResponse struct {
	Messages []*records.ChatMessage `json:"messages"`
}

type CreateRoomRequest struct {
	FromUUID string `json:"fromUuid"`
	ToUUID   string `json:"toUuid"`
}

type CreateRoomResponse struct {
	Room *records.ChatRoom `json:"room"`
}

type DeleteRoomRequest struct {
	RoomUUID string `json:"roomUuid"`
}

type DeleteRoomResponse struct {
}
