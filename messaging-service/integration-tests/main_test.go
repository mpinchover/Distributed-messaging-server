package integrestion_testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"messaging-service/types/enums"
	"messaging-service/types/requests"
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

func TestOpenSocket(t *testing.T) {
	t.Run("test set open socket info", func(t *testing.T) {

		clientUUID := uuid.New().String()
		setupClientConnection(t, clientUUID)

	})
}

func TestCreateRoom(t *testing.T) {
	t.Run("create room", func(t *testing.T) {
		tomUUID := uuid.New().String()
		jerryUUID := uuid.New().String()

		_, tomWS := setupClientConnection(t, tomUUID)
		_, jerryWS := setupClientConnection(t, jerryUUID)

		// create a room
		createRoomRequest := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: tomUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: jerryUUID,
					UserRole: "MEMBER",
				},
			},
		}

		postBody, err := json.Marshal(createRoomRequest)
		assert.NoError(t, err)
		reqBody := bytes.NewBuffer(postBody)

		resp, err := http.Post("http://localhost:9090/create-room", "application/json", reqBody)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.StatusCode >= 200 && resp.StatusCode <= 299)

		_, p, err := tomWS.ReadMessage()
		assert.NoError(t, err)

		// // get open room response over socket
		tomOpenRoomEventResponse := &requests.OpenRoomEvent{}
		err = json.Unmarshal(p, tomOpenRoomEventResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, tomOpenRoomEventResponse.EventType)
		assert.Equal(t, tomOpenRoomEventResponse.EventType, enums.EVENT_OPEN_ROOM.String())
		assert.NotNil(t, tomOpenRoomEventResponse.Room)
		assert.NotEmpty(t, tomOpenRoomEventResponse.Room.UUID)
		assert.Equal(t, 2, len(tomOpenRoomEventResponse.Room.Members))

		for _, m := range tomOpenRoomEventResponse.Room.Members {
			assert.Equal(t, "MEMBER", m.UserRole)
			assert.NotEmpty(t, m.UUID)
			assert.NotEmpty(t, m.UserUUID)
		}

		_, p, err = jerryWS.ReadMessage()
		assert.NoError(t, err)

		jerryOpenRoomEventResponse := requests.OpenRoomEvent{}
		err = json.Unmarshal(p, &jerryOpenRoomEventResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, jerryOpenRoomEventResponse.EventType)
		assert.Equal(t, jerryOpenRoomEventResponse.EventType, enums.EVENT_OPEN_ROOM.String())
		assert.NotNil(t, jerryOpenRoomEventResponse.Room)
		assert.NotEmpty(t, jerryOpenRoomEventResponse.Room.UUID)
		assert.Equal(t, 2, len(jerryOpenRoomEventResponse.Room.Members))

		for _, m := range jerryOpenRoomEventResponse.Room.Members {
			assert.Equal(t, "MEMBER", m.UserRole)
			assert.NotEmpty(t, m.UUID)
			assert.NotEmpty(t, m.UserUUID)
		}

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
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: bUUID,
					UserRole: "MEMBER",
				},
			},
		}
		openRoom(t, createRoomRequest1)

		openRoomRes1 := readOpenRoomResponse(t, aWebWS, 2)
		openRoomRes1 = readOpenRoomResponse(t, bWebWS, 2)
		roomUUID1 := openRoomRes1.Room.UUID

		createRoomRequest2 := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: cUUID,
					UserRole: "MEMBER",
				},
			},
		}
		openRoom(t, createRoomRequest2)
		openRoomRes2 := readOpenRoomResponse(t, cWebWS, 2)
		openRoomRes2 = readOpenRoomResponse(t, aWebWS, 2)

		roomUUID2 := openRoomRes2.Room.UUID

		// send messages between A and B
		sendMessages(t, aUUID, aConnectionUUID, roomUUID1, aWebWS)
		sendMessages(t, bUUID, bConnectionUUID, roomUUID1, bWebWS)

		// hanging here
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
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: dUUID,
					UserRole: "MEMBER",
				},
			},
		}
		openRoom(t, createRoomReq3)
		openRoomRes3 := readOpenRoomResponse(t, dWebWS, 2)
		openRoomRes3 = readOpenRoomResponse(t, aWebWS, 2)
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
			Members: []*requests.Member{
				{
					UserUUID: bUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: cUUID,
					UserRole: "MEMBER",
				},
			},
		}

		openRoom(t, openRoomReq4)
		openRoomRes4 := readOpenRoomResponse(t, bWebWS, 2)
		openRoomRes4 = readOpenRoomResponse(t, cWebWS, 2)
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
			Members: []*requests.Member{
				{
					UserUUID: bUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: dUUID,
					UserRole: "MEMBER",
				},
			},
		}
		openRoom(t, openRoomRequest5)
		openRoomRes5 := readOpenRoomResponse(t, dWebWS, 2)
		openRoomRes5 = readOpenRoomResponse(t, bWebWS, 2)

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

// Need to get the room id first and pass it to the text message id
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
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: bUUID,
					UserRole: "MEMBER",
				},
			},
		}
		openRoom(t, openRoomEvent)

		openRoomRes := readOpenRoomResponse(t, aWebWS, 2)
		readOpenRoomResponse(t, bWebWS, 2)
		readOpenRoomResponse(t, bMobileWS, 2)
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
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: bUUID,
					UserRole: "MEMBER",
				},
			},
		}
		openRoom(t, openRoomEvent)
		openRoomRes := readOpenRoomResponse(t, bWebWS, 2)
		roomUUID1 := openRoomRes.Room.UUID

		openRoomEvent = &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: cUUID,
					UserRole: "MEMBER",
				},
			},
		}
		openRoom(t, openRoomEvent)

		openRoomRes = readOpenRoomResponse(t, cWebWS, 2)
		roomUUID2 := openRoomRes.Room.UUID

		res, err := getRoomsByUserUUID(aUUID, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 2, len(res.Rooms))
		assert.Equal(t, 2, len(res.Rooms[0].Members))
		assert.Equal(t, 2, len(res.Rooms[1].Members))

		res, err = getRoomsByUserUUID(bUUID, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 1, len(res.Rooms))
		assert.Equal(t, 2, len(res.Rooms[0].Members))

		res, err = getRoomsByUserUUID(cUUID, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 1, len(res.Rooms))
		assert.Equal(t, 2, len(res.Rooms[0].Members))

		deleteRoom(t, &requests.DeleteRoomRequest{
			RoomUUID: roomUUID1,
		})

		res, err = getRoomsByUserUUID(aUUID, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 1, len(res.Rooms))
		assert.Equal(t, 2, len(res.Rooms[0].Members))

		res, err = getRoomsByUserUUID(bUUID, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 0, len(res.Rooms))

		res, err = getRoomsByUserUUID(cUUID, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 1, len(res.Rooms))
		assert.Equal(t, 2, len(res.Rooms[0].Members))

		deleteRoom(t, &requests.DeleteRoomRequest{
			RoomUUID: roomUUID2,
		})

		res, err = getRoomsByUserUUID(aUUID, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 0, len(res.Rooms))

		res, err = getRoomsByUserUUID(bUUID, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 0, len(res.Rooms))

		res, err = getRoomsByUserUUID(cUUID, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 0, len(res.Rooms))
	})
}

func TestLeaveRoom(t *testing.T) {
	t.Run("test leave room", func(t *testing.T) {
		aUUID := uuid.New().String()
		bUUID := uuid.New().String()
		cUUID := uuid.New().String()
		dUUID := uuid.New().String()

		_, aWebWS := setupClientConnection(t, aUUID)
		_, bWebWS := setupClientConnection(t, bUUID)
		_, cWebWS := setupClientConnection(t, cUUID)
		_, dWebWS := setupClientConnection(t, dUUID)
		_, dMobileWS := setupClientConnection(t, dUUID)

		openRoomEvent := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: bUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: cUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: dUUID,
					UserRole: "MEMBER",
				},
			},
		}
		openRoom(t, openRoomEvent)
		readOpenRoomResponse(t, aWebWS, 4)
		readOpenRoomResponse(t, bWebWS, 4)
		readOpenRoomResponse(t, cWebWS, 4)
		readOpenRoomResponse(t, dWebWS, 4)
		openRoomRes := readOpenRoomResponse(t, dMobileWS, 4)
		roomUUID := openRoomRes.Room.UUID

		res, err := getRoomsByUserUUID(aUUID, 0)
		assert.NoError(t, err)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 4)

		res, err = getRoomsByUserUUID(bUUID, 0)
		assert.NoError(t, err)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 4)

		res, err = getRoomsByUserUUID(cUUID, 0)
		assert.NoError(t, err)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 4)

		res, err = getRoomsByUserUUID(dUUID, 0)
		assert.NoError(t, err)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 4)

		leaveRoomReq := &requests.LeaveRoomRequest{
			UserUUID: cUUID,
			RoomUUID: roomUUID,
		}

		leaveRoom(t, leaveRoomReq)

		// // c should now be 0
		res, err = getRoomsByUserUUID(cUUID, 0)
		assert.NoError(t, err)
		assert.Len(t, res.Rooms, 0)

		// everyone else should still be 1
		res, err = getRoomsByUserUUID(aUUID, 0)
		assert.NoError(t, err)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 3)

		res, err = getRoomsByUserUUID(bUUID, 0)
		assert.NoError(t, err)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 3)

		res, err = getRoomsByUserUUID(dUUID, 0)
		assert.NoError(t, err)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 3)

		// TODO - send out a message that this person has left the chat
		// so you shouuld have an INFO type of msg
		// verify that the message has been sent out
	})
}

func queryMessages(t *testing.T, userUUID string, roomUUID string, expectedRooms int) {
	totalMessages := []*requests.Message{}
	res, err := getRoomsByUserUUID(userUUID, 0)

	assert.NoError(t, err)
	assert.NotEmpty(t, res)
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

	// jump by 15 because the msgs are being sent too fast.
	for i := 15; i < len(totalMessages); i++ {
		prev := totalMessages[i-1]
		cur := totalMessages[i]
		assert.True(t, prev.CreatedAt >= cur.CreatedAt)
	}
}

func sendMessages(t *testing.T, fromUserUUID string, connectionUUID string, roomUUID string, conn *websocket.Conn) {
	for i := 0; i < 25; i++ {
		msgText := fmt.Sprintf("Message %d", i)
		msgEventOut := &requests.TextMessageEvent{
			FromUUID:       fromUserUUID,
			ConnectionUUID: connectionUUID,
			RoomUUID:       roomUUID,
			EventType:      enums.EVENT_TEXT_MESSAGE.String(),
			MessageText:    msgText,
		}
		sendTextMessage(t, conn, msgEventOut)
	}
}

func recvMessages(t *testing.T, conn *websocket.Conn) {
	for i := 0; i < 25; i++ {
		// conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		_, p, err := conn.ReadMessage()

		assert.NoError(t, err)
		resp := &requests.TextMessageEvent{}
		err = json.Unmarshal(p, resp)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.EventType)
		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
		assert.NotEmpty(t, resp.FromUUID)
		assert.NotEmpty(t, resp.ConnectionUUID)
		assert.NotEmpty(t, resp.RoomUUID)
		assert.NotEmpty(t, resp.MessageText)
	}
}

func containsRoomUUID(s []*requests.Room, str string) bool {
	for _, v := range s {
		if v.UUID == str {
			return true
		}
	}

	return false
}

func readOpenRoomResponse(t *testing.T, conn *websocket.Conn, expectedMembers int) *requests.OpenRoomEvent {
	// TODO - ensure correct users are in the room
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)
	resp := &requests.OpenRoomEvent{}
	err = json.Unmarshal(p, resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.EventType)
	assert.Equal(t, resp.EventType, enums.EVENT_OPEN_ROOM.String())
	assert.NotEmpty(t, resp.Room)
	assert.NotEmpty(t, resp.Room.UUID)
	assert.Equal(t, expectedMembers, len(resp.Room.Members))

	return resp
}

func sendTextMessage(t *testing.T, ws *websocket.Conn, msgEvent *requests.TextMessageEvent) {
	err := ws.WriteJSON(msgEvent)
	assert.NoError(t, err)
}

func getMessagesByRoomUUID(t *testing.T, roomUUID string, offset int) (*requests.GetMessagesByRoomUUIDResponse, error) {
	url := fmt.Sprintf("http://localhost:9090/get-messages-by-room-uuid?roomUuid=%s&offset=%d", roomUUID, offset)
	resp, err := http.Get(url)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

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
func setupClientConnection(t *testing.T, userUUID string) (*requests.SetClientConnectionEvent, *websocket.Conn) {
	conn, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
	assert.NoError(t, err)

	msgOut := requests.SetClientConnectionEvent{
		FromUUID:  userUUID,
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
	}

	err = conn.WriteJSON(msgOut)
	assert.NoError(t, err)

	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)

	rsp := &requests.SetClientConnectionEvent{}
	err = json.Unmarshal(p, &rsp)
	assert.NoError(t, err)
	assert.NotEmpty(t, rsp.ConnectionUUID)
	assert.NotEmpty(t, rsp.FromUUID)
	return rsp, conn
}

func openRoom(t *testing.T, openRoomEvent *requests.CreateRoomRequest) {
	postBody, err := json.Marshal(openRoomEvent)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:9090/create-room", "application/json", reqBody)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode <= 299)

}

func deleteRoom(t *testing.T, deleteRoomRequest *requests.DeleteRoomRequest) {
	postBody, err := json.Marshal(deleteRoomRequest)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	_, err = http.Post("http://localhost:9090/delete-room", "application/json", reqBody)
	assert.NoError(t, err)
}

func leaveRoom(t *testing.T, req *requests.LeaveRoomRequest) {
	postBody, err := json.Marshal(req)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	_, err = http.Post("http://localhost:9090/leave-room", "application/json", reqBody)
	assert.NoError(t, err)
}
