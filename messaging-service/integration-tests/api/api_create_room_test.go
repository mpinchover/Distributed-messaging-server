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

func TestCreateRoom(t *testing.T) {
	// t.Skip()
	t.Run("create room", func(t *testing.T) {
		log.Printf("Running test %s", t.Name())
		validMessagingToken, validAPIKey := common.GetValidToken(t)

		tomClient, tomConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		jerryClient, jerryConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		// create a room
		createRoomRequest := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: tomClient.UserUUID,
				},
				{
					UserUUID: jerryClient.UserUUID,
				},
			},
		}

		common.OpenRoom(t, createRoomRequest, validAPIKey)

		tOpenRoomResponse := common.ReadOpenRoomResponse(t, tomConn, 2)
		jOpenRoomResponse := common.ReadOpenRoomResponse(t, jerryConn, 2)
		assert.Equal(t, tOpenRoomResponse.Room.UUID, jOpenRoomResponse.Room.UUID)

	})
}
