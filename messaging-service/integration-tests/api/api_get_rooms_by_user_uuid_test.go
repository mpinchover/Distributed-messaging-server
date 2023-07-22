package apitests

import (
	"log"
	"messaging-service/integration-tests/common"
	"messaging-service/types/requests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetByUUIDWithApiKey(t *testing.T) {
	// t.Skip()
	t.Run("test get rooms with api key", func(t *testing.T) {
		log.Printf("Running %s", t.Name())

		// need to get valid API key as well
		_, validAPIKey := common.GetValidToken(t)

		tom := uuid.New().String()
		jerry := uuid.New().String()

		for i := 0; i < 50; i++ {
			createRoomRequest := &requests.CreateRoomRequest{
				Members: []*requests.Member{
					{
						UserUUID: tom,
					},
					{
						UserUUID: uuid.New().String(),
					},
				},
			}
			common.OpenRoom(t, createRoomRequest, validAPIKey)
		}

		for i := 0; i < 5; i++ {
			createRoomRequest := &requests.CreateRoomRequest{
				Members: []*requests.Member{
					{
						UserUUID: jerry,
					},
					{
						UserUUID: uuid.New().String(),
					},
				},
			}
			common.OpenRoom(t, createRoomRequest, validAPIKey)
		}

		// test tom
		totalRooms := []*requests.Room{}
		roomsResponse := common.GetRoomsByUserUUIDWithApiKey(t, tom, 0, validAPIKey)
		assert.NotNil(t, roomsResponse)
		totalRooms = append(totalRooms, roomsResponse.Rooms...)

		assert.Len(t, roomsResponse.Rooms, 20)
		assert.Len(t, totalRooms, 20)

		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, tom, len(totalRooms), validAPIKey)
		assert.NotNil(t, roomsResponse)
		totalRooms = append(totalRooms, roomsResponse.Rooms...)

		assert.Len(t, roomsResponse.Rooms, 20)
		assert.Len(t, totalRooms, 40)

		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, tom, len(totalRooms), validAPIKey)
		assert.NotNil(t, roomsResponse)
		totalRooms = append(totalRooms, roomsResponse.Rooms...)

		assert.Len(t, roomsResponse.Rooms, 10)
		assert.Len(t, totalRooms, 50)

		// test jerry
		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, jerry, 0, validAPIKey)
		assert.NotNil(t, roomsResponse)

		assert.Len(t, roomsResponse.Rooms, 5)

		// // send messages between A and B
		// // send 150 messages
		// common.SendMessages(t, aClient.UserUUID, aClient.ConnectionUUID, roomUUID, aConn, validMessagingToken)
		// common.SendMessages(t, bClient.UserUUID, bClient.ConnectionUUID, roomUUID, bConn, validMessagingToken)
		// common.SendMessages(t, aClient.UserUUID, aClient.ConnectionUUID, roomUUID, aConn, validMessagingToken)
		// common.SendMessages(t, bClient.UserUUID, bClient.ConnectionUUID, roomUUID, bConn, validMessagingToken)
		// common.SendMessages(t, aClient.UserUUID, aClient.ConnectionUUID, roomUUID, aConn, validMessagingToken)
		// common.SendMessages(t, bClient.UserUUID, bClient.ConnectionUUID, roomUUID, bConn, validMessagingToken)

		// time.Sleep(1 * time.Second)
		// common.GetRoomsByUserUUIDWithApiKey(t, aClient.UserUUID, roomUUID, 2, validAPIKey)
	})

}
