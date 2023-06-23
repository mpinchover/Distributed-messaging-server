package requests

type GetRoomsByUserUUIDRequest struct {
	UserUUID string `schema:"userUuid" validate:"required"`
	Offset   int    `schema:"offset"`
}

type GetRoomsByUserUUIDResponse struct {
	Rooms []*Room `json:"rooms"`
}

type GetMessagesByRoomUUIDRequest struct {
	RoomUUID string `schema:"roomUuid" validate:"required"`
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
	RoomUUID string `json:"roomUuid" validate:"required"`
	UserUUID string `json:"userUuid" validate:"required"`
}

type DeleteRoomResponse struct {
}

type LeaveRoomRequest struct {
	UserUUID string `json:"userUuid" validate:"required"`
	RoomUUID string `json:"roomUuid" validate:"required"`
}

type LeaveRoomResponse struct {
}
