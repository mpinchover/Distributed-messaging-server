package integrationtests

import (
	"log"
	"messaging-service/integration-tests/common"
	"messaging-service/types/enums"
	"messaging-service/types/requests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSeenBy(t *testing.T) {
	// t.Skip()
	t.Run("delete room and messages", func(t *testing.T) {

		log.Printf("Running test %s", t.Name())

		validMessagingToken, validAPIKey := common.GetValidToken(t)
		tomClient, tomConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		aliceClient, aliceConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		jerryClient, jerryConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		deanClient, deanConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		_, deanMobileConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  deanClient.UserUUID,
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
				{
					UserUUID: aliceClient.UserUUID,
				},
				{
					UserUUID: deanClient.UserUUID,
				},
			},
		}
		common.OpenRoom(t, createRoomRequest, validAPIKey)

		openRoomResponse := common.ReadOpenRoomResponse(t, tomConn, 4)

		common.ReadOpenRoomResponse(t, jerryConn, 4)
		common.ReadOpenRoomResponse(t, aliceConn, 4)
		common.ReadOpenRoomResponse(t, deanConn, 4)
		common.ReadOpenRoomResponse(t, deanMobileConn, 4)

		// send out a message tom -> room
		roomUUID := openRoomResponse.Room.UUID
		msgEventOut := &requests.TextMessageEvent{
			FromUUID:       tomClient.UserUUID,
			ConnectionUUID: tomClient.ConnectionUUID,
			EventType:      enums.EVENT_TEXT_MESSAGE.String(),
			Message: &requests.Message{
				MessageText: "TEXT",
				RoomUUID:    roomUUID,
			},
			Token: validMessagingToken,
		}
		common.SendTextMessage(t, tomConn, msgEventOut)

		// everyone should recv the message

		// clear out recv msg
		resp := &requests.TextMessageEvent{}
		common.RecvMessage(t, jerryConn, resp)
		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
		common.RecvMessage(t, aliceConn, resp)
		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
		common.RecvMessage(t, deanConn, resp)
		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
		common.RecvMessage(t, deanMobileConn, resp)
		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)

		// send seen event jerry -> room
		seenEvent := &requests.SeenMessageEvent{
			EventType:   enums.EVENT_SEEN_MESSAGE.String(),
			MessageUUID: resp.Message.UUID,
			UserUUID:    jerryClient.UserUUID,
			RoomUUID:    roomUUID,
			Token:       validMessagingToken,
		}

		err := jerryConn.WriteJSON(seenEvent)
		assert.NoError(t, err)

		// everyone should get the seen event message
		common.RecvSeenMessageEvent(t, tomConn, resp.Message.UUID)
		common.RecvSeenMessageEvent(t, aliceConn, resp.Message.UUID)
		common.RecvSeenMessageEvent(t, deanConn, resp.Message.UUID)
		common.RecvSeenMessageEvent(t, deanMobileConn, resp.Message.UUID)

		res, err := common.GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, 0, validMessagingToken)
		assert.NoError(t, err)
		assert.Len(t, res.Messages, 1)
		assert.Len(t, res.Messages[0].SeenBy, 1)

		assert.Equal(t, res.Messages[0].SeenBy[0].MessageUUID, resp.Message.UUID)
		assert.Equal(t, res.Messages[0].SeenBy[0].UserUUID, jerryClient.UserUUID)

	})
}
