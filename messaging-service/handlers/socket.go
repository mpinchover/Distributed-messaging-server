package handlers

import (
	"encoding/json"
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

func (h *Handler) handleIncomingSocketEvents(conn *websocket.Conn) error {

	for {
		// read in a message
		_, p, err := conn.ReadMessage()

		if err != nil {
			break
		}

		// if err != nil && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		// 	break
		// }

		// if err != nil {
		// 	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		// 		break
		// 	}
		// }

		msgType, err := utils.GetEventType(string(p))
		if err != nil {
			panic(err)
		}

		if msgType == enums.EVENT_SET_CLIENT_SOCKET.String() {

			// TODO – have a new event that doesn't include the connectionUUID
			msg := &requests.SetClientConnectionEvent{}
			err := json.Unmarshal(p, msg)
			if err != nil {
				return err
			}
			resp, err := h.ControlTowerCtrlr.SetupClientConnectionV2(conn, msg)
			if err != nil {
				return err
			}
			err = conn.WriteJSON(resp)
			if err != nil {
				return err
			}
		}

		if msgType == enums.EVENT_TEXT_MESSAGE.String() {
			msg := &requests.TextMessageEvent{}
			err := json.Unmarshal(p, msg)
			if err != nil {
				return err
			}
			// break this up into processTextMessage and SaveTextMessage
			_, err = h.ControlTowerCtrlr.ProcessTextMessage(msg)
			if err != nil {
				err = conn.WriteJSON([]byte("could not send text message"))
				if err != nil {
					log.Println("error sending err msg")
				}
			}
		}

		if msgType == enums.EVENT_SEEN_MESSAGE.String() {
			msg := &requests.SeenMessageEvent{}
			err := json.Unmarshal(p, msg)
			if err != nil {
				return err
			}
			h.ControlTowerCtrlr.SaveSeenBy(msg)
		}
	}

	return nil
}
