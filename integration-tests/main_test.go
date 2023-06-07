package integrestion_testing

import (
	"bytes"
	"chatapi/types"
	"chatapi/utils"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

const (
	socketURL = "ws://localhost:9090/ws"
)

func TestMockEndpoint(t *testing.T) {
	t.Run("mock test", func(t *testing.T) {
		assert.Equal(t, 123, 123, "they should be equal")
	})
	t.Run("mock test", func(t *testing.T) {
		assert.NotEqual(t, 124, 123, "they should be equal")
	})
}

func TestConnectWebsocket(t *testing.T) {
	t.Run("test opening websocket", func(t *testing.T) {

		ws, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
		assert.NoError(t, err)

		pingHandler := ws.PingHandler()
		err = pingHandler("PING")
		assert.NoError(t, err)

		_, p, err := ws.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, "PONG", string(p))
	})
}

func TestSetClientSocketInfo(t *testing.T) {
	t.Run("test set open socket info on client", func(t *testing.T) {

		ws, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
		assert.NoError(t, err)

		clientUUID := uuid.New().String()
		msgOut := types.SetClientConnectionEvent{
			FromUUID:  utils.ToStrPtr(clientUUID),
			EventType: utils.ToStrPtr("EVENT_SET_CLIENT_SOCKET"),
		}

		err = ws.WriteJSON(msgOut)
		assert.NoError(t, err)

		_, p, err := ws.ReadMessage()
		assert.NoError(t, err)

		msgIn := types.SetClientConnectionEvent{}
		err = json.Unmarshal(p, &msgIn)
		assert.NoError(t, err)

		assert.NotNil(t, msgIn.ConnectionUUID)
		assert.NotNil(t, msgIn.FromUUID)
	})
	t.Run("test create a room", func(t *testing.T) {
		tomUUID := uuid.New().String()
		jerryUUID := uuid.New().String()

		_, tomWS := setupClientConnection(t, tomUUID)
		_, jerryWS := setupClientConnection(t, jerryUUID)

		// create a room
		openRoomEvent := &types.OpenRoomRequest{
			FromUUID: utils.ToStrPtr(tomUUID),
			ToUUID:   utils.ToStrPtr(jerryUUID),
		}
		openRoom(t, openRoomEvent)

		_, p, err := tomWS.ReadMessage()
		assert.NoError(t, err)

		tomOpenRoomEventResponse := types.OpenRoomEvent{}
		err = json.Unmarshal(p, &tomOpenRoomEventResponse)
		assert.NoError(t, err)
		assert.NotNil(t, tomOpenRoomEventResponse.EventType)
		assert.NotNil(t, tomOpenRoomEventResponse.FromUUID)
		assert.NotNil(t, tomOpenRoomEventResponse.ToUUID)
		assert.NotNil(t, tomOpenRoomEventResponse.Room)
		assert.NotNil(t, tomOpenRoomEventResponse.Room.UUID)
		assert.Equal(t, 2, len(tomOpenRoomEventResponse.Room.Participants))

		// ensure size of channels
		_, p, err = jerryWS.ReadMessage()
		assert.NoError(t, err)

		jerryOpenRoomEventResponse := types.OpenRoomEvent{}
		err = json.Unmarshal(p, &jerryOpenRoomEventResponse)
		assert.NoError(t, err)
		assert.NotNil(t, jerryOpenRoomEventResponse.EventType)
		assert.NotNil(t, jerryOpenRoomEventResponse.FromUUID)
		assert.NotNil(t, jerryOpenRoomEventResponse.ToUUID)
		assert.NotNil(t, jerryOpenRoomEventResponse.Room)
		assert.NotNil(t, jerryOpenRoomEventResponse.Room.UUID)
		assert.Equal(t, 2, len(jerryOpenRoomEventResponse.Room.Participants))

		// ensure the room is the same room
		assert.Equal(t, jerryOpenRoomEventResponse.Room.UUID, tomOpenRoomEventResponse.Room.UUID)
	})

	t.Run("test send messages across a room between two connections", func(t *testing.T) {
		// set up ws connections
		tomUUID := uuid.New().String()
		jerryUUID := uuid.New().String()

		tomClient, tomWS := setupClientConnection(t, tomUUID)
		jerryClient, jerryWS := setupClientConnection(t, jerryUUID)

		// create a room
		openRoomEvent := &types.OpenRoomRequest{
			FromUUID: utils.ToStrPtr(tomUUID),
			ToUUID:   utils.ToStrPtr(jerryUUID),
		}
		openRoom(t, openRoomEvent)

		tOpenRoomResp := readOpenRoomResponse(t, tomWS)
		jOpenRoomResp := readOpenRoomResponse(t, jerryWS)

		// send first text message
		msgEventOut := &types.ChatMessageEvent{
			FromUserUUID:       &tomUUID,
			FromConnectionUUID: tomClient.ConnectionUUID,
			RoomUUID:           tOpenRoomResp.Room.UUID,
			EventType:          utils.ToStrPtr("EVENT_CHAT_TEXT"),
			MessageText:        utils.ToStrPtr("Message 1"),
		}
		sendTextMessage(t, tomWS, msgEventOut)

		// read the first text message
		msgEventIn := readTextMessage(t, jerryWS)
		assert.Equal(t, msgEventOut.MessageText, msgEventIn.MessageText)

		// send the second text message
		msgEventOut = &types.ChatMessageEvent{
			FromUserUUID:       &jerryUUID,
			FromConnectionUUID: jerryClient.ConnectionUUID,
			RoomUUID:           jOpenRoomResp.Room.UUID,
			EventType:          utils.ToStrPtr("EVENT_CHAT_TEXT"),
			MessageText:        utils.ToStrPtr("Message 2"),
		}
		sendTextMessage(t, jerryWS, msgEventOut)

		// read the second text message
		msgEventIn = readTextMessage(t, tomWS)
		assert.Equal(t, msgEventOut.MessageText, msgEventIn.MessageText)
	})

	// EVENT_CHAT_TEXT
}

func readOpenRoomResponse(t *testing.T, conn *websocket.Conn) *types.OpenRoomEvent {
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)

	resp := &types.OpenRoomEvent{}
	err = json.Unmarshal(p, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.EventType)
	assert.NotNil(t, resp.FromUUID)
	assert.NotNil(t, resp.ToUUID)
	assert.NotNil(t, resp.Room)
	assert.NotNil(t, resp.Room.UUID)
	assert.Equal(t, 2, len(resp.Room.Participants))

	return resp
}

func readTextMessage(t *testing.T, conn *websocket.Conn) *types.ChatMessageEvent {
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)

	msg := &types.ChatMessageEvent{}
	err = json.Unmarshal(p, msg)
	assert.NoError(t, err)
	assert.NotNil(t, msg.FromConnectionUUID)
	assert.NotNil(t, msg.FromUserUUID)
	assert.NotNil(t, msg.MessageText)
	assert.NotNil(t, msg.EventType)
	assert.NotNil(t, msg.RoomUUID)
	return msg
}

func sendTextMessage(t *testing.T, ws *websocket.Conn, msgEvent *types.ChatMessageEvent) {
	err := ws.WriteJSON(msgEvent)
	assert.NoError(t, err)

}

// set up a client connection
func setupClientConnection(t *testing.T, userUUID string) (*types.SetClientConnectionEvent, *websocket.Conn) {
	conn, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
	assert.NoError(t, err)

	msgOut := types.SetClientConnectionEvent{
		FromUUID:  utils.ToStrPtr(userUUID),
		EventType: utils.ToStrPtr("EVENT_SET_CLIENT_SOCKET"),
	}

	err = conn.WriteJSON(msgOut)
	assert.NoError(t, err)

	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)

	rsp := &types.SetClientConnectionEvent{}
	err = json.Unmarshal(p, &rsp)
	assert.NoError(t, err)
	assert.NotNil(t, rsp.ConnectionUUID)
	assert.NotNil(t, rsp.FromUUID)
	return rsp, conn
}

func openRoom(t *testing.T, openRoomEvent *types.OpenRoomRequest) {
	postBody, err := json.Marshal(openRoomEvent)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	_, err = http.Post("http://localhost:9090/create-room", "application/json", reqBody)
	assert.NoError(t, err)
}
