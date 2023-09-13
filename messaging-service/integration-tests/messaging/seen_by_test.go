package integrationtests

import (
	"messaging-service/src/types/enums"
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"

	"github.com/google/uuid"
)

func (s *IntegrationTestSuite) TestSeenBy() {

	tomUUID := uuid.New().String()
	jerryUUID := uuid.New().String()
	aliceUUID := uuid.New().String()
	deanUUID := uuid.New().String()

	tomToken := s.GetValidToken(tomUUID)
	aliceToken := s.GetValidToken(aliceUUID)
	jerryToken := s.GetValidToken(jerryUUID)
	deanToken := s.GetValidToken(deanUUID)

	apiKey := s.GetValidAPIKey()
	// validToken := s.GetValidToken(tomUUID)

	tomClient, tomConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     tomToken,
		UserUUID:  tomUUID,
	})

	_, aliceConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     aliceToken,
		UserUUID:  aliceUUID,
	})

	jerryClient, jerryConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     jerryToken,
		UserUUID:  jerryUUID,
	})

	_, deanConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     deanToken,
		UserUUID:  deanUUID,
	})

	_, deanMobileConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     deanToken,
		UserUUID:  deanUUID,
	})

	// create a room
	createRoomRequest := &requests.CreateRoomRequest{
		Members: []*records.Member{
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

	s.OpenRoom(createRoomRequest, apiKey)
	openRoomResponse := s.ReadOpenRoomResponse(tomConn, 4)

	s.ReadOpenRoomResponse(jerryConn, 4)
	s.ReadOpenRoomResponse(aliceConn, 4)
	s.ReadOpenRoomResponse(deanConn, 4)
	s.ReadOpenRoomResponse(deanMobileConn, 4)
	roomUUID := openRoomResponse.Room.UUID

	// send out a message tom -> room
	msgEventOut := &requests.TextMessageEvent{
		FromUUID:   tomClient.UserUUID,
		DeviceUUID: tomClient.DeviceUUID,
		EventType:  enums.EVENT_TEXT_MESSAGE.String(),
		Message: &records.Message{
			MessageText: "TEXT",
			RoomUUID:    roomUUID,
		},
		Token: tomToken,
	}
	s.SendTextMessage(tomConn, msgEventOut)

	// clear out recv msg
	resp := &requests.TextMessageEvent{}
	s.RecvMessage(jerryConn, resp)
	s.Equal(enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
	s.RecvMessage(aliceConn, resp)
	s.Equal(enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
	s.RecvMessage(deanConn, resp)
	s.Equal(enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
	s.RecvMessage(deanMobileConn, resp)
	s.Equal(enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)

	// send seen event jerry -> room
	seenEvent := &requests.SeenMessageEvent{
		EventType:   enums.EVENT_SEEN_MESSAGE.String(),
		MessageUUID: resp.Message.UUID,
		UserUUID:    jerryClient.UserUUID,
		RoomUUID:    roomUUID,
		Token:       jerryToken,
	}

	err := jerryConn.WriteJSON(seenEvent)
	s.NoError(err)

	// everyone should get the seen event message
	s.RecvSeenMessageEvent(tomConn, resp.Message.UUID)
	s.RecvSeenMessageEvent(aliceConn, resp.Message.UUID)
	s.RecvSeenMessageEvent(deanConn, resp.Message.UUID)
	s.RecvSeenMessageEvent(deanMobileConn, resp.Message.UUID)

	msgs := s.MakeGetMessagesByRoomUUIDRequest(roomUUID, apiKey, 0)
	s.NoError(err)
	s.Len(msgs.Messages, 1)
	s.Len(msgs.Messages[0].SeenBy, 1)

	s.Equal(msgs.Messages[0].SeenBy[0].MessageUUID, resp.Message.UUID)
	s.Equal(msgs.Messages[0].SeenBy[0].UserUUID, jerryClient.UserUUID)
}
