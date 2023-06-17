package handlers

import (
	"encoding/json"
	"messaging-service/types/requests"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
)

var decoder = schema.NewDecoder()

func (h *Handler) GetRoomsByUserUUID(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// validation
	req := &requests.GetRoomsByUserUUIDRequest{}
	err := decoder.Decode(req, r.URL.Query())
	if err != nil {
		return nil, err
	}
	return h.getRoomsByUserUUID(req)
}

func (h *Handler) GetMessagesByRoomUUID(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// validation
	req := &requests.GetMessagesByRoomUUIDRequest{}
	err := decoder.Decode(req, r.URL.Query())
	if err != nil {
		return nil, err
	}

	return h.getMessagesByRoomUUID(req)
}

func (h *Handler) SetupWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	go h.ControlTowerCtrlr.SetupClientConnection(conn)
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// validation
	req := &requests.CreateRoomRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	return h.createRoom(req)
}

func (h *Handler) DeleteRoom(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// validation
	req := &requests.DeleteRoomRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	return h.deleteRoom(req)
}
