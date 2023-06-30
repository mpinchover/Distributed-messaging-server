package integrationtests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"messaging-service/types/enums"
	"messaging-service/types/requests"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

const (
	SocketURL = "ws://localhost:9090/ws"
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

		ws, _, err := websocket.DefaultDialer.Dial(SocketURL, nil)
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

		recvMessages(t, bWebWS)
		recvMessages(t, aWebWS)

		time.Sleep(1 * time.Second)
		queryMessages(t, bUUID, roomUUID1, 1)
		queryMessages(t, aUUID, roomUUID1, 2)

		// send messages between A and C
		sendMessages(t, aUUID, aConnectionUUID, roomUUID2, aWebWS)
		sendMessages(t, cUUID, cConnectionUUID, roomUUID2, cWebWS)

		recvMessages(t, aWebWS)
		recvMessages(t, cWebWS)

		time.Sleep(1 * time.Second)
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

		time.Sleep(1 * time.Second)
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

		time.Sleep(1 * time.Second)
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

		time.Sleep(100 * time.Millisecond)
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

		_, aWS := setupClientConnection(t, aUUID)
		_, bWS := setupClientConnection(t, bUUID)
		_, cWS := setupClientConnection(t, cUUID)

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
		readOpenRoomResponse(t, aWS, 2)
		openRoomRes := readOpenRoomResponse(t, bWS, 2)
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

		readOpenRoomResponse(t, aWS, 2)
		openRoomRes = readOpenRoomResponse(t, cWS, 2)
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

		// ensure delete event is recd
		resp := &requests.DeleteRoomEvent{}
		err = readEvent(aWS, resp)
		assert.NoError(t, err)
		assert.Equal(t, enums.EVENT_DELETE_ROOM.String(), resp.EventType)
		assert.Equal(t, roomUUID1, resp.RoomUUID)

		resp = &requests.DeleteRoomEvent{}
		err = readEvent(bWS, resp)
		assert.NoError(t, err)
		assert.Equal(t, enums.EVENT_DELETE_ROOM.String(), resp.EventType)
		assert.Equal(t, roomUUID1, resp.RoomUUID)

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

		// ensure delete event is recd
		resp = &requests.DeleteRoomEvent{}
		err = readEvent(aWS, resp)
		assert.NoError(t, err)
		assert.Equal(t, enums.EVENT_DELETE_ROOM.String(), resp.EventType)
		assert.Equal(t, roomUUID2, resp.RoomUUID)

		resp = &requests.DeleteRoomEvent{}
		err = readEvent(cWS, resp)
		assert.NoError(t, err)
		assert.Equal(t, enums.EVENT_DELETE_ROOM.String(), resp.EventType)
		assert.Equal(t, roomUUID2, resp.RoomUUID)
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

		// read the message from leaving the room
		resp := &requests.LeaveRoomEvent{}
		err = readEvent(aWebWS, resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, cUUID, resp.UserUUID)
		assert.Equal(t, roomUUID, resp.RoomUUID)
		assert.Equal(t, enums.EVENT_LEAVE_ROOM.String(), resp.EventType)

		resp = &requests.LeaveRoomEvent{}
		err = readEvent(bWebWS, resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, cUUID, resp.UserUUID)
		assert.Equal(t, roomUUID, resp.RoomUUID)
		assert.Equal(t, enums.EVENT_LEAVE_ROOM.String(), resp.EventType)

		resp = &requests.LeaveRoomEvent{}
		err = readEvent(dWebWS, resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, cUUID, resp.UserUUID)
		assert.Equal(t, roomUUID, resp.RoomUUID)
		assert.Equal(t, enums.EVENT_LEAVE_ROOM.String(), resp.EventType)

		resp = &requests.LeaveRoomEvent{}
		err = readEvent(dMobileWS, resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, cUUID, resp.UserUUID)
		assert.Equal(t, roomUUID, resp.RoomUUID)
		assert.Equal(t, enums.EVENT_LEAVE_ROOM.String(), resp.EventType)
	})
}

func readEvent(conn *websocket.Conn, v interface{}) error {
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	_, p, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	err = json.Unmarshal(p, v)
	return err
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
	assert.Equal(t, 20, len(totalMessages))

	resp, err = getMessagesByRoomUUID(t, roomUUID, len(totalMessages))
	assert.NoError(t, err)
	assert.Equal(t, 20, len(resp.Messages))
	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, 40, len(totalMessages))

	resp, err = getMessagesByRoomUUID(t, roomUUID, len(totalMessages))
	assert.NoError(t, err)
	assert.Equal(t, 10, len(resp.Messages))
	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, 50, len(totalMessages))

	// jump by 15 because the msgs are being sent too fast.
	for i := 15; i < len(totalMessages); i++ {
		prev := totalMessages[i-1]
		cur := totalMessages[i]
		assert.True(t, prev.CreatedAt >= cur.CreatedAt)
	}

}
