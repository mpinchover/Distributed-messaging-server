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

func TestSeenBy(t *testing.T) {
	t.Run("delete room and messages", func(t *testing.T) {
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
		openRoom(t, createRoomRequest)

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
			RoomUUID:       roomUUID,
			EventType:      enums.EVENT_TEXT_MESSAGE.String(),
			MessageText:    "TEXT",
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
			MessageUUID: resp.MessageUUID,
			UserUUID:    jerryUUID,
			RoomUUID:    roomUUID,
		}

		err = jerryWS.WriteJSON(seenEvent)
		assert.NoError(t, err)

		// everyone should get the seen event message
		recvSeenMessageEvent(t, tomWS, resp.MessageUUID)
		recvSeenMessageEvent(t, aliceWS, resp.MessageUUID)
		recvSeenMessageEvent(t, deanWS, resp.MessageUUID)
		recvSeenMessageEvent(t, deanMobileWS, resp.MessageUUID)

		res, err := getMessagesByRoomUUID(t, roomUUID, 0)
		assert.NoError(t, err)
		assert.Len(t, res.Messages, 1)
		assert.Len(t, res.Messages[0].SeenBy, 1)

		assert.Equal(t, res.Messages[0].SeenBy[0].MessageUUID, resp.MessageUUID)
		assert.Equal(t, res.Messages[0].SeenBy[0].UserUUID, jerryUUID)

	})
}

func recvSeenMessageEvent(t *testing.T, conn *websocket.Conn, messageUUID string) {
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)
	seenMessageEvent := &requests.SeenMessageEvent{}
	err = json.Unmarshal(p, seenMessageEvent)
	assert.NoError(t, err)

	assert.NotEmpty(t, seenMessageEvent.EventType)
	assert.Equal(t, enums.EVENT_SEEN_MESSAGE.String(), seenMessageEvent.EventType)
	assert.NotEmpty(t, seenMessageEvent.MessageUUID)
	assert.Equal(t, messageUUID, seenMessageEvent.MessageUUID)
	assert.NotEmpty(t, seenMessageEvent.RoomUUID)
	assert.NotEmpty(t, seenMessageEvent.UserUUID)
}
