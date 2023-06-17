package integrestion_testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"messaging-service/types/events"
	"messaging-service/types/records"
	"messaging-service/types/requests"
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
	t.Skip()
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
		msgOut := events.SetClientConnectionEvent{
			FromUUID:  clientUUID,
			EventType: "EVENT_SET_CLIENT_SOCKET",
		}

		err = ws.WriteJSON(msgOut)
		assert.NoError(t, err)

		_, p, err := ws.ReadMessage()
		assert.NoError(t, err)

		msgIn := events.SetClientConnectionEvent{}
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
		createRoomRequest := &requests.CreateRoomRequest{
			FromUUID: tomUUID,
			ToUUID:   jerryUUID,
		}
		openRoom(t, createRoomRequest)
		_, p, err := tomWS.ReadMessage()
		assert.NoError(t, err)

		tomOpenRoomEventResponse := events.OpenRoomEvent{}
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

		jerryOpenRoomEventResponse := events.OpenRoomEvent{}
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
}

func TestRoomAndMessagesPagination(t *testing.T) {
	t.Run("test rooms and messages pagination", func(t *testing.T) {
		aUUID := uuid.New().String()
		bUUID := uuid.New().String()
		cUUID := uuid.New().String()
		dUUID := uuid.New().String()

		aResp, aWebWS := setupClientConnection(t, aUUID)
		bResp, bWebWS := setupClientConnection(t, bUUID)
		cResp, cWebWS := setupClientConnection(t, cUUID)
		dResp, dWebWS := setupClientConnection(t, dUUID)

		aConnectionUUID := aResp.ConnectionUUID
		bConnectionUUID := bResp.ConnectionUUID
		cConnectionUUID := cResp.ConnectionUUID
		dConnectionUUID := dResp.ConnectionUUID

		createRoomRequest1 := &requests.CreateRoomRequest{
			FromUUID: aUUID,
			ToUUID:   bUUID,
		}
		openRoom(t, createRoomRequest1)

		openRoomRes1 := readOpenRoomResponse(t, aWebWS)
		openRoomRes1 = readOpenRoomResponse(t, bWebWS)
		roomUUID1 := openRoomRes1.Room.UUID

		createRoomRequest2 := &requests.CreateRoomRequest{
			FromUUID: aUUID,
			ToUUID:   cUUID,
		}
		openRoom(t, createRoomRequest2)
		openRoomRes2 := readOpenRoomResponse(t, cWebWS)
		openRoomRes2 = readOpenRoomResponse(t, aWebWS)
		roomUUID2 := openRoomRes2.Room.UUID

		// send messages between A and B
		sendMessages(t, aUUID, aConnectionUUID, roomUUID1, aWebWS)
		sendMessages(t, bUUID, bConnectionUUID, roomUUID1, bWebWS)

		recvMessages(t, bWebWS)
		recvMessages(t, aWebWS)

		queryMessages(t, bUUID, roomUUID1, 1)
		queryMessages(t, aUUID, roomUUID1, 2)

		// send messages between A and C
		sendMessages(t, aUUID, aConnectionUUID, roomUUID2, aWebWS)
		sendMessages(t, cUUID, cConnectionUUID, roomUUID2, cWebWS)

		recvMessages(t, aWebWS)
		recvMessages(t, cWebWS)

		queryMessages(t, aUUID, roomUUID2, 2)
		queryMessages(t, cUUID, roomUUID2, 1)

		// create room between A and D
		createRoomReq3 := &requests.CreateRoomRequest{
			FromUUID: aUUID,
			ToUUID:   dUUID,
		}
		openRoom(t, createRoomReq3)
		openRoomRes3 := readOpenRoomResponse(t, dWebWS)
		openRoomRes3 = readOpenRoomResponse(t, aWebWS)
		roomUUID3 := openRoomRes3.Room.UUID

		// send messages between A and D
		sendMessages(t, aUUID, aConnectionUUID, roomUUID3, aWebWS)
		sendMessages(t, dUUID, dConnectionUUID, roomUUID3, dWebWS)

		recvMessages(t, aWebWS)
		recvMessages(t, dWebWS)

		queryMessages(t, aUUID, roomUUID3, 3)
		queryMessages(t, dUUID, roomUUID3, 1)

		// create room between B and C
		openRoomReq4 := &requests.CreateRoomRequest{
			FromUUID: bUUID,
			ToUUID:   cUUID,
		}

		openRoom(t, openRoomReq4)
		openRoomRes4 := readOpenRoomResponse(t, bWebWS)
		openRoomRes4 = readOpenRoomResponse(t, cWebWS)
		roomUUID4 := openRoomRes4.Room.UUID

		// send messages between B and C
		sendMessages(t, bUUID, bConnectionUUID, roomUUID4, bWebWS)
		sendMessages(t, cUUID, cConnectionUUID, roomUUID4, cWebWS)

		recvMessages(t, bWebWS)
		recvMessages(t, cWebWS)

		queryMessages(t, bUUID, roomUUID4, 2)
		queryMessages(t, cUUID, roomUUID4, 2)

		// create room between B and D
		openRoomRequest5 := &requests.CreateRoomRequest{
			FromUUID: bUUID,
			ToUUID:   dUUID,
		}
		openRoom(t, openRoomRequest5)
		openRoomRes5 := readOpenRoomResponse(t, dWebWS)
		openRoomRes5 = readOpenRoomResponse(t, bWebWS)

		// the mobiel device will get the open room msg as well
		roomUUID5 := openRoomRes5.Room.UUID

		// send messages between B and D
		sendMessages(t, bUUID, bConnectionUUID, roomUUID5, bWebWS)
		sendMessages(t, dUUID, dConnectionUUID, roomUUID5, dWebWS)

		recvMessages(t, bWebWS)
		recvMessages(t, dWebWS)

		queryMessages(t, bUUID, roomUUID5, 3)
		queryMessages(t, dUUID, roomUUID5, 2)

	})
}

func TestAllConnectionsRcvMessages(t *testing.T) {
	t.Run("test all connections get msgs", func(t *testing.T) {
		aUUID := uuid.New().String()
		bUUID := uuid.New().String()

		aWebResp, aWebWS := setupClientConnection(t, aUUID)
		bWebResp, bWebWS := setupClientConnection(t, bUUID)
		_, bMobileWS := setupClientConnection(t, bUUID)

		aWebConnUUID := aWebResp.ConnectionUUID
		bWebConnUUID := bWebResp.ConnectionUUID

		openRoomEvent := &requests.CreateRoomRequest{
			FromUUID: aUUID,
			ToUUID:   bUUID,
		}
		openRoom(t, openRoomEvent)

		openRoomRes := readOpenRoomResponse(t, aWebWS)
		readOpenRoomResponse(t, bWebWS)
		readOpenRoomResponse(t, bMobileWS)
		roomUUID := openRoomRes.Room.UUID

		sendMessages(t, bUUID, aWebConnUUID, roomUUID, aWebWS)
		sendMessages(t, bUUID, bWebConnUUID, roomUUID, bWebWS)

		recvMessages(t, bWebWS)
		recvMessages(t, aWebWS)

		// need to recv double the msgs
		recvMessages(t, bMobileWS)
		recvMessages(t, bMobileWS)
		queryMessages(t, bUUID, roomUUID, 1)
		queryMessages(t, aUUID, roomUUID, 1)

		// add new connection
		_, aMobileWS := setupClientConnection(t, aUUID)

		sendMessages(t, bUUID, aWebConnUUID, roomUUID, aWebWS)
		sendMessages(t, bUUID, bWebConnUUID, roomUUID, bWebWS)

		recvMessages(t, bWebWS)
		recvMessages(t, aWebWS)

		// need to recv double the msgs
		recvMessages(t, bMobileWS)
		recvMessages(t, bMobileWS)

		// need to recv double the msgs
		recvMessages(t, aMobileWS)
		recvMessages(t, aMobileWS)

	})
}

func TestDeleteRoom(t *testing.T) {
	t.Run("test delete a room", func(t *testing.T) {
		aUUID := uuid.New().String()
		bUUID := uuid.New().String()
		cUUID := uuid.New().String()

		setupClientConnection(t, aUUID)
		_, bWebWS := setupClientConnection(t, bUUID)
		_, cWebWS := setupClientConnection(t, cUUID)

		openRoomEvent := &requests.CreateRoomRequest{
			FromUUID: aUUID,
			ToUUID:   bUUID,
		}
		openRoom(t, openRoomEvent)
		openRoomRes := readOpenRoomResponse(t, bWebWS)
		roomUUID1 := openRoomRes.Room.UUID

		openRoomEvent = &requests.CreateRoomRequest{
			FromUUID: aUUID,
			ToUUID:   cUUID,
		}
		openRoom(t, openRoomEvent)

		openRoomRes = readOpenRoomResponse(t, cWebWS)
		roomUUID2 := openRoomRes.Room.UUID

		res, err := getRoomsByUserUUID(aUUID, 0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res.Rooms))

		res, err = getRoomsByUserUUID(bUUID, 0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 1, len(res.Rooms))

		res, err = getRoomsByUserUUID(cUUID, 0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 1, len(res.Rooms))

		deleteRoom(t, &requests.DeleteRoomRequest{
			RoomUUID: roomUUID1,
		})

		res, err = getRoomsByUserUUID(aUUID, 0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 1, len(res.Rooms))

		res, err = getRoomsByUserUUID(bUUID, 0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 0, len(res.Rooms))

		res, err = getRoomsByUserUUID(cUUID, 0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 1, len(res.Rooms))

		deleteRoom(t, &requests.DeleteRoomRequest{
			RoomUUID: roomUUID2,
		})

		res, err = getRoomsByUserUUID(aUUID, 0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 0, len(res.Rooms))

		res, err = getRoomsByUserUUID(bUUID, 0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 0, len(res.Rooms))

		res, err = getRoomsByUserUUID(cUUID, 0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 0, len(res.Rooms))
	})
}

func queryMessages(t *testing.T, userUUID string, roomUUID string, expectedRooms int) {
	totalMessages := []*records.ChatMessage{}
	res, err := getRoomsByUserUUID(userUUID, 0)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expectedRooms, len(res.Rooms))

	// ensure it contains the room uuid
	assert.True(t, containsRoomUUID(res.Rooms, roomUUID))

	resp, err := getMessagesByRoomUUID(t, roomUUID, 0)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(resp.Messages))
	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, len(totalMessages), 20)

	resp, err = getMessagesByRoomUUID(t, roomUUID, len(totalMessages))
	assert.NoError(t, err)
	assert.Equal(t, 20, len(resp.Messages))
	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, len(totalMessages), 40)

	resp, err = getMessagesByRoomUUID(t, roomUUID, len(totalMessages))
	assert.NoError(t, err)
	assert.Equal(t, 10, len(resp.Messages))
	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, len(totalMessages), 50)

	for i := 1; i < len(totalMessages); i++ {
		prev := totalMessages[i-1]
		cur := totalMessages[i]

		assert.Greater(t, prev.ID, cur.ID)
	}
}

func sendMessages(t *testing.T, fromUserUUID string, connectionUUID string, roomUUID string, conn *websocket.Conn) {
	for i := 0; i < 25; i++ {
		msgText := fmt.Sprintf("Message %d", i)
		msgEventOut := &events.ChatMessageEvent{
			FromUserUUID:       fromUserUUID,
			FromConnectionUUID: connectionUUID,
			RoomUUID:           roomUUID,
			EventType:          "EVENT_CHAT_TEXT",
			MessageText:        msgText,
		}
		sendTextMessage(t, conn, msgEventOut)
	}
}

func recvMessages(t *testing.T, conn *websocket.Conn) {
	for i := 0; i < 25; i++ {
		_, p, err := conn.ReadMessage()
		assert.NoError(t, err)
		resp := &events.ChatMessageEvent{}
		err = json.Unmarshal(p, resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.EventType)
		assert.Equal(t, "EVENT_CHAT_TEXT", resp.EventType)
		assert.NotNil(t, resp.FromUserUUID)
		assert.NotNil(t, resp.FromConnectionUUID)
		assert.NotNil(t, resp.RoomUUID)
		assert.NotNil(t, resp.MessageText)
	}
}

func containsRoomUUID(s []*records.ChatRoom, str string) bool {
	for _, v := range s {
		if v.UUID == str {
			return true
		}
	}

	return false
}

func readOpenRoomResponse(t *testing.T, conn *websocket.Conn) *events.OpenRoomEvent {
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)
	resp := &events.OpenRoomEvent{}
	err = json.Unmarshal(p, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.EventType)
	assert.Equal(t, resp.EventType, "EVENT_OPEN_ROOM")
	assert.NotNil(t, resp.FromUUID)
	assert.NotNil(t, resp.ToUUID)
	assert.NotNil(t, resp.Room)
	assert.NotNil(t, resp.Room.UUID)
	assert.Equal(t, 2, len(resp.Room.Participants))

	return resp
}

func sendTextMessage(t *testing.T, ws *websocket.Conn, msgEvent *events.ChatMessageEvent) {
	err := ws.WriteJSON(msgEvent)
	assert.NoError(t, err)
}

func getMessagesByRoomUUID(t *testing.T, roomUUID string, offset int) (*requests.GetMessagesByRoomUUIDResponse, error) {
	url := fmt.Sprintf("http://localhost:9090/get-messages-by-room-uuid?roomUuid=%s&offset=%d", roomUUID, offset)
	resp, err := http.Get(url)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotNil(t, b)

	result := &requests.GetMessagesByRoomUUIDResponse{}
	err = json.Unmarshal(b, result)
	assert.NoError(t, err)
	return result, err
}

func getRoomsByUserUUID(userUUID string, offset int) (*requests.GetRoomsByUserUUIDResponse, error) {
	url := fmt.Sprintf("http://localhost:9090/get-rooms-by-user-uuid?userUuid=%s&offset=%d", userUUID, offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &requests.GetRoomsByUserUUIDResponse{}
	err = json.Unmarshal(b, result)
	return result, err
}

// set up a client connection
func setupClientConnection(t *testing.T, userUUID string) (*events.SetClientConnectionEvent, *websocket.Conn) {
	conn, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
	assert.NoError(t, err)

	msgOut := events.SetClientConnectionEvent{
		FromUUID:  userUUID,
		EventType: "EVENT_SET_CLIENT_SOCKET",
	}

	err = conn.WriteJSON(msgOut)
	assert.NoError(t, err)

	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)

	rsp := &events.SetClientConnectionEvent{}
	err = json.Unmarshal(p, &rsp)
	assert.NoError(t, err)
	assert.NotNil(t, rsp.ConnectionUUID)
	assert.NotNil(t, rsp.FromUUID)
	return rsp, conn
}

func openRoom(t *testing.T, openRoomEvent *requests.CreateRoomRequest) {
	postBody, err := json.Marshal(openRoomEvent)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	_, err = http.Post("http://localhost:9090/create-room", "application/json", reqBody)
	assert.NoError(t, err)
}

func deleteRoom(t *testing.T, deleteRoomRequest *requests.DeleteRoomRequest) {
	postBody, err := json.Marshal(deleteRoomRequest)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	_, err = http.Post("http://localhost:9090/delete-room", "application/json", reqBody)
	assert.NoError(t, err)
}
