package handlers

import (
	"encoding/json"
	"messaging-service/src/types/requests"

	"github.com/gorilla/websocket"
)

func (h *Handler) handleClientEventSeenMessage(conn *websocket.Conn, p []byte) error {
	msg := &requests.SeenMessageEvent{}
	err := json.Unmarshal(p, msg)
	if err != nil {
		return err
	}
	return h.ControlTowerCtrlr.SaveSeenBy(msg)
}

func (h *Handler) handleClientEventTextMessage(conn *websocket.Conn, p []byte) error {
	msg := &requests.TextMessageEvent{}
	err := json.Unmarshal(p, msg)
	if err != nil {
		return err
	}
	// break this up into processTextMessage and SaveTextMessage
	_, err = h.ControlTowerCtrlr.ProcessTextMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) handleSetClientSocket(ws *requests.Websocket, p []byte) error {
	// TODO – have a new event that doesn't include the connectionUUID
	msg := &requests.SetClientConnectionEvent{}
	err := json.Unmarshal(p, msg)
	if err != nil {
		return err
	}
	resp, err := h.ControlTowerCtrlr.SetupClientConnectionV2(ws.Conn, msg)
	if err != nil {
		return err
	}
	ws.DeviceUUID = &resp.DeviceUUID
	ws.UserUUID = &resp.UserUUID
	return ws.Conn.WriteJSON(resp)
}
