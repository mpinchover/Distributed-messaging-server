package integrationtests

import (
	"bytes"
	"encoding/json"
	"log"
	"messaging-service/types/enums"
	"messaging-service/types/requests"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoom(t *testing.T) {
	t.Run("create room", func(t *testing.T) {
		log.Printf("Running test %s", t.Name())
		tomUUID := uuid.New().String()
		jerryUUID := uuid.New().String()

		_, tomWS := setupClientConnection(t, tomUUID)
		_, jerryWS := setupClientConnection(t, jerryUUID)

		// create a room
		createRoomRequest := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: tomUUID,
				},
				{
					UserUUID: jerryUUID,
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
			assert.NotEmpty(t, m.UUID)
			assert.NotEmpty(t, m.UserUUID)
		}

		// ensure the room is the same room
		assert.Equal(t, jerryOpenRoomEventResponse.Room.UUID, tomOpenRoomEventResponse.Room.UUID)

	})
}
