package integrationtests

import (
	"encoding/json"
	"log"
	"messaging-service/types/enums"
	"messaging-service/types/requests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSeenBy(t *testing.T) {
	t.Run("delete room and messages", func(t *testing.T) {
		log.Printf("Running test %s", t.Name())
		tomUUID := uuid.New().String()
		jerryUUID := uuid.New().String()
		aliceUUID := uuid.New().String()
		deanUUID := uuid.New().String()

		tResp, tomWS := setupClientConnection(t, tomUUID)
		_, jerryWS := setupClientConnection(t, jerryUUID)
		_, aliceWS := setupClientConnection(t, aliceUUID)
		_, deanWS := setupClientConnection(t, deanUUID)
		_, deanMobileWS := setupClientConnection(t, deanUUID)

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
					UserUUID: deanUUID,
				},
			},
		}
		err := openRoom(createRoomRequest)
		assert.NoError(t, err)

		// read message from create room
		_, p, err := tomWS.ReadMessage()
		assert.NoError(t, err)

		// get open room response over socket
		tomOpenRoomEventResponse := &requests.OpenRoomEvent{}
		err = json.Unmarshal(p, tomOpenRoomEventResponse)
		assert.NoError(t, err)

		// clear other WS's
		_, _, err = jerryWS.ReadMessage()
		assert.NoError(t, err)
		_, _, err = aliceWS.ReadMessage()
		assert.NoError(t, err)
		_, _, err = deanWS.ReadMessage()
		assert.NoError(t, err)
		_, _, err = deanMobileWS.ReadMessage()
		assert.NoError(t, err)

		// send out a message tom -> room
		roomUUID := tomOpenRoomEventResponse.Room.UUID
		msgEventOut := &requests.TextMessageEvent{
			FromUUID:       tomUUID,
			ConnectionUUID: tResp.ConnectionUUID,
			EventType:      enums.EVENT_TEXT_MESSAGE.String(),
			Message: &requests.Message{
				MessageText: "TEXT",
				RoomUUID:    roomUUID,
			},
		}
		sendTextMessage(t, tomWS, msgEventOut)

		// everyone should recv the message

		// clear out recv msg
		resp := &requests.TextMessageEvent{}
		recvMessage(t, jerryWS, resp)
		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
		recvMessage(t, aliceWS, resp)
		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
		recvMessage(t, deanWS, resp)
		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
		recvMessage(t, deanMobileWS, resp)
		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)

		// send seen event jerry -> room
		seenEvent := &requests.SeenMessageEvent{
			EventType:   enums.EVENT_SEEN_MESSAGE.String(),
			MessageUUID: resp.Message.UUID,
			UserUUID:    jerryUUID,
			RoomUUID:    roomUUID,
		}

		err = jerryWS.WriteJSON(seenEvent)
		assert.NoError(t, err)

		// everyone should get the seen event message
		recvSeenMessageEvent(t, tomWS, resp.Message.UUID)
		recvSeenMessageEvent(t, aliceWS, resp.Message.UUID)
		recvSeenMessageEvent(t, deanWS, resp.Message.UUID)
		recvSeenMessageEvent(t, deanMobileWS, resp.Message.UUID)

		res, err := getMessagesByRoomUUID(t, roomUUID, 0)
		assert.NoError(t, err)
		assert.Len(t, res.Messages, 1)
		assert.Len(t, res.Messages[0].SeenBy, 1)

		assert.Equal(t, res.Messages[0].SeenBy[0].MessageUUID, resp.Message.UUID)
		assert.Equal(t, res.Messages[0].SeenBy[0].UserUUID, jerryUUID)

	})
}
