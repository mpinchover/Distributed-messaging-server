package requests

type GetRoomsByUserUUIDRequest struct {
	UserUUID string `schema:"userUuid"`
	Offset   int    `schema:"offset"`
}

type GetRoomsByUserUUIDResponse struct {
	Rooms []*Room `json:"rooms"`
}

type GetMessagesByRoomUUIDRequest struct {
	RoomUUID string `schema:"roomUuid"`
	Offset   int    `schema:"offset"`
}

type GetMessagesByRoomUUIDResponse struct {
	Messages []*Message `json:"messages"`
}

type CreateRoomRequest struct {
	Members []*Member `json:"participants"`
}

type CreateRoomResponse struct {
	Room *Room `json:"room"`
}

type DeleteRoomRequest struct {
	RoomUUID string `json:"roomUuid"`
	UserUUID string `json:"userUuid"`
}

type DeleteRoomResponse struct {
}

type LeaveRoomRequest struct {
	UserUUID string `json:"userUuid"`
	RoomUUID string `json:"roomUuid"`
}

type LeaveRoomResponse struct {
}
