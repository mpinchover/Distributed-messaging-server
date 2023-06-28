package integrationtests

import (
	"encoding/json"
	"messaging-service/types/enums"
	"messaging-service/types/requests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRoomAndMessages(t *testing.T) {
	t.Skip()
	t.Run("delete room and messages", func(t *testing.T) {
		tomUUID := uuid.New().String()
		jerryUUID := uuid.New().String()

		tResp, tomWS := setupClientConnection(t, tomUUID)
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

		openRoom(t, createRoomRequest)

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

		roomUUID := tomOpenRoomEventResponse.Room.UUID
		sendMessages(t, tomUUID, tResp.ConnectionUUID, roomUUID, tomWS)
		recvMessages(t, jerryWS)

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
	})
}
