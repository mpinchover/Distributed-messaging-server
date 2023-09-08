package connections

import "github.com/gorilla/websocket"

// map user uuid -> devices

// TODO - possibly make this a map
type ChatConnections map[string][]*Device
type Device struct {
	WS *websocket.Conn
}

type UserConnection struct {
	UUID    string
	Devices map[string]*Device
}

// room uuid -> participants in the room
type Channels map[string][]string
