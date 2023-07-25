package requests

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
	UserUUID       string `json:"userUuid"`
	ConnectionUUID string `json:"connectionUuid"`
	Token          string `json:"token"`
}

type TextMessageEvent struct {
	EventType      string   `json:"eventType"`
	FromUUID       string   `json:"fromUuid"`
	ConnectionUUID string   `json:"connectionUuid"`
	Message        *Message `json:"message"`
	Token          string   `json:"token"`
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
