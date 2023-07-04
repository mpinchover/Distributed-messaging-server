package integrationtests

import (
	"encoding/json"
	"log"
	"messaging-service/types/enums"
	"messaging-service/types/requests"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRoom(t *testing.T) {
	t.Run("test delete a room", func(t *testing.T) {
		log.Printf("Running test %s", t.Name())
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
				},
				{
					UserUUID: bUUID,
				},
			},
		}
		err := openRoom(openRoomEvent)
		assert.NoError(t, err)
		readOpenRoomResponse(t, aWS, 2)
		openRoomRes := readOpenRoomResponse(t, bWS, 2)
		roomUUID1 := openRoomRes.Room.UUID

		openRoomEvent = &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
				},
				{
					UserUUID: cUUID,
				},
			},
		}
		err = openRoom(openRoomEvent)
		assert.NoError(t, err)

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

func TestDeleteRoomAndMessages(t *testing.T) {
	t.Run("delete room and messages", func(t *testing.T) {
		log.Printf("Running test %s", t.Name())
		tomUUID := uuid.New().String()
		jerryUUID := uuid.New().String()
		aliceUUID := uuid.New().String()
		benUUID := uuid.New().String()

		tResp, tomWS := setupClientConnection(t, tomUUID)
		_, jerryWS := setupClientConnection(t, jerryUUID)
		_, aliceWS := setupClientConnection(t, aliceUUID)
		_, benWS1 := setupClientConnection(t, benUUID)
		_, benWS2 := setupClientConnection(t, benUUID)
		_, benWS3 := setupClientConnection(t, benUUID)

		// create a room
		createRoomRequest := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: tomUUID,
				},
				{
					UserUUID: jerryUUID,
				},
				{
					UserUUID: aliceUUID,
				},
				{
					UserUUID: benUUID,
				},
			},
		}

		err := openRoom(createRoomRequest)
		assert.NoError(t, err)

		// // get open room response over socket
		tomOpenRoomEventResponse := readOpenRoomResponse(t, tomWS, 4)

		readOpenRoomResponse(t, jerryWS, 4)
		readOpenRoomResponse(t, aliceWS, 4)
		readOpenRoomResponse(t, benWS1, 4)
		readOpenRoomResponse(t, benWS2, 4)
		readOpenRoomResponse(t, benWS3, 4)

		roomUUID := tomOpenRoomEventResponse.Room.UUID
		sendMessages(t, tomUUID, tResp.ConnectionUUID, roomUUID, tomWS)
		recvMessages(t, jerryWS)
		recvMessages(t, aliceWS)
		recvMessages(t, benWS1)
		recvMessages(t, benWS2)
		recvMessages(t, benWS3)

		res, err := getMessagesByRoomUUID(t, roomUUID, 0)
		assert.NoError(t, err)

		assert.Len(t, res.Messages, 20)

		deleteRoomRequest := &requests.DeleteRoomRequest{
			RoomUUID: roomUUID,
		}
		deleteRoom(t, deleteRoomRequest)
		res, err = getMessagesByRoomUUID(t, roomUUID, 0)
		assert.NoError(t, err)
		assert.Len(t, res.Messages, 0)

		// ensure everyone got the deletedRoom event
		e := &requests.DeleteRoomEvent{}
		recvDeletedRoomMsg(t, tomWS, e)
		recvDeletedRoomMsg(t, jerryWS, e)
		recvDeletedRoomMsg(t, aliceWS, e)
		recvDeletedRoomMsg(t, benWS1, e)
		recvDeletedRoomMsg(t, benWS2, e)
		recvDeletedRoomMsg(t, benWS3, e)
	})
}

func recvDeletedRoomMsg(t *testing.T, conn *websocket.Conn, resp *requests.DeleteRoomEvent) {
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)
	err = json.Unmarshal(p, resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.EventType)
	assert.Equal(t, enums.EVENT_DELETE_ROOM.String(), resp.EventType)
	assert.NotEmpty(t, resp.RoomUUID)
}
