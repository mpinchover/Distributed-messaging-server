package requests

type DeleteRoomEvent struct {
	EventType string `json:"eventType"`
	RoomUUID  string `json:"roomUuid"`
	UserUUID  string `json:"userUuid"`
}

type LeaveRoomEvent struct {
	EventType string `json:"eventType"`
	RoomUUID  string `json:"roomUuid"`
	UserUUID  string `json:"userUuid"`
}

// sennd to clients room has been opened
type OpenRoomEvent struct {
	EventType string `json:"eventType"`
	Room      *Room  `json:"room"`
}

// subscrve the sever to a room
type SubscribeToRoomEvent struct {
	EventType string   `json:"eventType"`
	Channel   string   `json:"channel"`
	Members   []string `json:"members"`
}

type SetClientConnectionEvent struct {
	EventType      string `json:"eventType"`
	FromUUID       string `json:"fromUuid"`
	ConnectionUUID string `json:"connectionUuid"`
}

type TextMessageEvent struct {
	EventType      string `json:"eventType"`
	FromUUID       string `json:"fromUuid"`
	ConnectionUUID string `json:"connectionUuid"`
	RoomUUID       string `json:"roomUuid"`
	MessageText    string `json:"messageText"`
	CreatedAt      int64  `json:"createdAt"`
}
