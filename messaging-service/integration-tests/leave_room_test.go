package integrationtests

import (
	"log"
	"messaging-service/types/enums"
	"messaging-service/types/requests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLeaveRoom(t *testing.T) {
	t.Skip()
	t.Run("test leave room", func(t *testing.T) {
		log.Printf("Running test %s", t.Name())
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
				},
				{
					UserUUID: bUUID,
				},
				{
					UserUUID: cUUID,
				},
				{
					UserUUID: dUUID,
				},
			},
		}
		err := openRoom(openRoomEvent)
		assert.NoError(t, err)
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
