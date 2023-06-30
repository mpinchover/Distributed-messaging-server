package integrationtests

import (
	"encoding/json"
	"messaging-service/types/enums"
	"messaging-service/types/requests"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRoomAndMessages(t *testing.T) {
	t.Run("delete room and messages", func(t *testing.T) {
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
					UserRole: "MEMBER",
				},
				{
					UserUUID: jerryUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: aliceUUID,
					UserRole: "MEMBER",
				},
				{
					UserUUID: benUUID,
					UserRole: "MEMBER",
				},
			},
		}

		openRoom(t, createRoomRequest)

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
