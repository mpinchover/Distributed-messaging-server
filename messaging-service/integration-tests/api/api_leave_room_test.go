package apitests

import (
	"log"
	"messaging-service/integration-tests/common"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLeaveRoom(t *testing.T) {
	// t.Skip()
	t.Run("test leave room", func(t *testing.T) {
		log.Printf("Running test %s", t.Name())

		validMessagingToken, validAPIKey := common.GetValidToken(t)

		aClient, aConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		bClient, bConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		cClient, cConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		dClient, dConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		_, dMobileConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  dClient.UserUUID,
		})

		openRoomEvent := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aClient.UserUUID,
				},
				{
					UserUUID: bClient.UserUUID,
				},
				{
					UserUUID: cClient.UserUUID,
				},
				{
					UserUUID: dClient.UserUUID,
				},
			},
		}
		common.OpenRoom(t, openRoomEvent, validAPIKey)
		common.ReadOpenRoomResponse(t, aConn, 4)
		common.ReadOpenRoomResponse(t, bConn, 4)
		common.ReadOpenRoomResponse(t, cConn, 4)
		common.ReadOpenRoomResponse(t, dConn, 4)
		openRoomRes := common.ReadOpenRoomResponse(t, dMobileConn, 4)
		roomUUID := openRoomRes.Room.UUID

		res := common.GetRoomsByUserUUIDByMessagingJWT(t, cClient.UserUUID, 0, validMessagingToken)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 4)

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, bClient.UserUUID, 0, validMessagingToken)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 4)

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, cClient.UserUUID, 0, validMessagingToken)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 4)

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, dClient.UserUUID, 0, validMessagingToken)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 4)

		leaveRoomReq := &requests.LeaveRoomRequest{
			UserUUID: cClient.UserUUID,
			RoomUUID: roomUUID,
		}

		common.LeaveRoom(t, leaveRoomReq, validAPIKey)

		// c should now be 0
		res = common.GetRoomsByUserUUIDByMessagingJWT(t, cClient.UserUUID, 0, validMessagingToken)
		assert.Len(t, res.Rooms, 0)

		// everyone else should still be 1
		res = common.GetRoomsByUserUUIDByMessagingJWT(t, aClient.UserUUID, 0, validMessagingToken)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 3)

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, bClient.UserUUID, 0, validMessagingToken)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 3)

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, dClient.UserUUID, 0, validMessagingToken)
		assert.Len(t, res.Rooms, 1)
		assert.Len(t, res.Rooms[0].Members, 3)

		// read the message from leaving the room
		resp := &requests.LeaveRoomEvent{}
		common.ReadEvent(t, aConn, resp)
		assert.NotNil(t, resp)
		assert.Equal(t, cClient.UserUUID, resp.UserUUID)
		assert.Equal(t, roomUUID, resp.RoomUUID)
		assert.Equal(t, enums.EVENT_LEAVE_ROOM.String(), resp.EventType)

		resp = &requests.LeaveRoomEvent{}
		common.ReadEvent(t, bConn, resp)
		assert.NotNil(t, resp)
		assert.Equal(t, cClient.UserUUID, resp.UserUUID)
		assert.Equal(t, roomUUID, resp.RoomUUID)
		assert.Equal(t, enums.EVENT_LEAVE_ROOM.String(), resp.EventType)

		resp = &requests.LeaveRoomEvent{}
		common.ReadEvent(t, dConn, resp)
		assert.NotNil(t, resp)
		assert.Equal(t, cClient.UserUUID, resp.UserUUID)
		assert.Equal(t, roomUUID, resp.RoomUUID)
		assert.Equal(t, enums.EVENT_LEAVE_ROOM.String(), resp.EventType)

		resp = &requests.LeaveRoomEvent{}
		common.ReadEvent(t, dMobileConn, resp)
		assert.NotNil(t, resp)
		assert.Equal(t, cClient.UserUUID, resp.UserUUID)
		assert.Equal(t, roomUUID, resp.RoomUUID)
		assert.Equal(t, enums.EVENT_LEAVE_ROOM.String(), resp.EventType)
	})
}
