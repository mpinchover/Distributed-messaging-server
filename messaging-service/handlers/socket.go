package handlers

import (
	"log"
	"messaging-service/types/enums"
	"messaging-service/types/requests"
	"messaging-service/utils"
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

	err = h.handleIncomingSocketEvents(conn)
	if err != nil {
		log.Println(err)
	}

}

func sendClientError(conn *websocket.Conn, err error) error {
	errResp := requests.ErrorResponse{
		Message: err.Error(),
	}
	conn.WriteJSON(errResp)
	return err
}

func (h *Handler) handleIncomingSocketEvents(conn *websocket.Conn) error {

	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// add in token authenticator

		// if err != nil && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		// 	break
		// }

		// if err != nil {
		// 	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		// 		break
		// 	}
		// }

		// TODO – error message for websockets, don't just panic
		msgType, err := utils.GetEventType(string(p))
		if err != nil {
			errResp := requests.ErrorResponse{
				Message: err.Error(),
			}
			conn.WriteJSON(errResp)
		}

		msgToken, err := utils.GetEventToken(string(p))
		if err != nil {
			sendClientError(conn, err)
		}

		var authErr error
		if msgType == enums.EVENT_SET_CLIENT_SOCKET.String() {
			_, authErr = h.AuthController.VerifyJWT(msgToken, true)
		} else {
			_, authErr = h.AuthController.VerifyJWT(msgToken, false)
		}

		if authErr != nil {
			return sendClientError(conn, err)
		}

		if msgType == enums.EVENT_SET_CLIENT_SOCKET.String() {
			err := h.handleSetClientSocket(conn, p)
			if err != nil {
				return sendClientError(conn, err)
			}

		}

		if msgType == enums.EVENT_TEXT_MESSAGE.String() {
			err := h.handleClientEventTextMessage(conn, p)
			if err != nil {
				sendClientError(conn, err)
			}
		}

		if msgType == enums.EVENT_SEEN_MESSAGE.String() {
			err := h.handleClientEventSeenMessage(conn, p)
			if err != nil {
				sendClientError(conn, err)
			}
		}
	}

	return nil
}
