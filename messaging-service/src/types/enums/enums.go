package enums

type MessageType int64

const (
	EVENT_TEXT_MESSAGE MessageType = iota
	EVENT_CHAT_TEXT_METADATA
	EVENT_OPEN_ROOM         // open a chat room request
	EVENT_SET_CLIENT_SOCKET // set the client socket
	EVENT_DELETE_ROOM       // delete a room
	EVENT_LEAVE_ROOM        // leave a room
	EVENT_SUBSCRIBE_TO_ROOM // subscribe to a room
	EVENT_SEEN_MESSAGE      // recpt saw message
	EVENT_DELETE_MESSAGE    // message was deleted
)

const (

	// server channel for server side events
	CHANNEL_SERVER_EVENTS = "CHANNEL_SERVER_EVENTS"
)

func (m MessageType) String() string {
	switch m {
	case EVENT_DELETE_MESSAGE:
		return "EVENT_DELETE_MESSAGE"
	case EVENT_TEXT_MESSAGE:
		return "EVENT_TEXT_MESSAGE"
	case EVENT_DELETE_ROOM:
		return "EVENT_DELETE_ROOM"
	case EVENT_CHAT_TEXT_METADATA:
		return "EVENT_CHAT_TEXT_METADATA"
	case EVENT_OPEN_ROOM:
		return "EVENT_OPEN_ROOM"
	case EVENT_SET_CLIENT_SOCKET:
		return "EVENT_SET_CLIENT_SOCKET"
	case EVENT_SUBSCRIBE_TO_ROOM:
		return "EVENT_SUBSCRIBE_TO_ROOM"
	case EVENT_SEEN_MESSAGE:
		return "EVENT_SEEN_MESSAGE"
	}
	return "UNKNOWN"
}

type MessageStatus int64

const (
	MESSAGE_STATUS_LIVE MessageStatus = iota
	MESSAGE_STATUS_DELETED
)

func (m MessageStatus) String() string {
	switch m {
	case MESSAGE_STATUS_LIVE:
		return "MESSAGE_STATUS_LIVE"
	case MESSAGE_STATUS_DELETED:
		return "MESSAGE_STATUS_DELETED"
	}
	return "UNKNOWN"
}

type AbortCode int64

const (
	ABORT_CODE_NEED_MORE_LIKED_QUESTIONS AbortCode = iota
	ABORT_CODE_NO_MATCHES
	ABORT_CODE_NO_OVERLAPPING_QUESTIONS
)

func (m AbortCode) String() string {
	switch m {
	case ABORT_CODE_NEED_MORE_LIKED_QUESTIONS:
		return "ABORT_CODE_NEED_MORE_LIKED_QUESTIONS"
	case ABORT_CODE_NO_MATCHES:
		return "ABORT_CODE_NO_MATCHES"
	case ABORT_CODE_NO_OVERLAPPING_QUESTIONS:
		return "ABORT_CODE_NO_OVERLAPPING_QUESTIONS"
	}
	return "UNKNOWN"
}
