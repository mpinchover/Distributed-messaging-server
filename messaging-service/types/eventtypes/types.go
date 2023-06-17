package eventtypes

type MessageType int64

const (
	EVENT_CHAT_TEXT MessageType = iota
	EVENT_CHAT_TEXT_METADATA
	EVENT_OPEN_ROOM         // open a chat room request
	EVENT_SET_CLIENT_SOCKET // set the client socket
	EVENT_DELETE_ROOM       // delete a room
	EVENT_LEAVE_ROOM        // leave a room
)

const (

	// server channel for server side events
	CHANNEL_SERVER_EVENTS = "CHANNEL_SERVER_EVENTS"
)

func (m MessageType) String() string {
	switch m {
	case EVENT_CHAT_TEXT:
		return "EVENT_CHAT_TEXT"
	case EVENT_DELETE_ROOM:
		return "EVENT_DELETE_ROOM"
	case EVENT_CHAT_TEXT_METADATA:
		return "EVENT_CHAT_TEXT_METADATA"
	case EVENT_OPEN_ROOM:
		return "EVENT_OPEN_ROOM"
	case EVENT_SET_CLIENT_SOCKET:
		return "EVENT_SET_CLIENT_SOCKET"
	}
	return "UNKNOWN"
}
