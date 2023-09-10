package handlers

import (
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"
	"messaging-service/src/utils"
	"net/http"

	"github.com/gorilla/websocket"
)

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

	ws := &requests.Websocket{
		Conn: conn,
	}
	// handle breaking the connection
	defer func() {
		conn.Close()
	}()

	conn.SetPongHandler(func(appData string) error {
		err := conn.WriteMessage(1, []byte("PONG"))
		if err != nil {
			panic(err)
		}
		return nil
	})

	err = h.handleIncomingSocketEvents(ws)
	if err != nil {
		if ws.UserUUID != nil && ws.DeviceUUID != nil {
			h.ControlTowerCtrlr.RemoveClientDeviceFromServer(*ws.UserUUID, *ws.DeviceUUID)
		}
	}
}

func sendClientError(ws *requests.Websocket, err error) error {
	errResp := requests.ErrorResponse{
		Message: err.Error(),
	}
	ws.Conn.WriteJSON(errResp)
	return err
}

func (h *Handler) handleIncomingSocketEvents(ws *requests.Websocket) error {

	for {
		// read in a message
		_, p, err := ws.Conn.ReadMessage()
		if err != nil {
			return err
		}

		// TODO – error message for websockets, don't just panic
		msgType, err := utils.GetEventType(string(p))
		if err != nil {
			errResp := requests.ErrorResponse{
				Message: err.Error(),
			}
			ws.Conn.WriteJSON(errResp)
		}

		msgToken, err := utils.GetEventToken(string(p))
		if err != nil {
			sendClientError(ws, err)
		}

		var authErr error
		if msgType == enums.EVENT_SET_CLIENT_SOCKET.String() {
			_, authErr = utils.VerifyJWT(msgToken, true)
		} else {
			_, authErr = utils.VerifyJWT(msgToken, false)
		}

		if authErr != nil {
			return sendClientError(ws, err)
		}

		if msgType == enums.EVENT_SET_CLIENT_SOCKET.String() {
			err := h.handleSetClientSocket(ws, p)
			if err != nil {
				return sendClientError(ws, err)
			}

		}

		if msgType == enums.EVENT_TEXT_MESSAGE.String() {
			err := h.handleClientEventTextMessage(ws.Conn, p)
			if err != nil {
				sendClientError(ws, err)
			}
		}

		if msgType == enums.EVENT_SEEN_MESSAGE.String() {
			err := h.handleClientEventSeenMessage(ws.Conn, p)
			if err != nil {
				sendClientError(ws, err)
			}
		}
	}

	return nil
}
