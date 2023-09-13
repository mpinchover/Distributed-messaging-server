package integrationtests

import (
	"messaging-service/src/types/enums"
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"

	"github.com/google/uuid"
)

func (s *IntegrationTestSuite) TestRoomAndMessagesPagination() {

	a := uuid.New().String()
	b := uuid.New().String()
	c := uuid.New().String()
	d := uuid.New().String()

	aValidToken := s.GetValidToken(a)
	bValidToken := s.GetValidToken(b)
	cValidToken := s.GetValidToken(c)
	dValidToken := s.GetValidToken(d)

	apiKey := s.GetValidAPIKey()

	aClient, aConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     aValidToken,
		UserUUID:  a,
	})

	bClient, bConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     bValidToken,
		UserUUID:  b,
	})

	cClient, cConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     cValidToken,
		UserUUID:  c,
	})

	dClient, dConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     dValidToken,
		UserUUID:  d,
	})

	createRoomRequest1 := &requests.CreateRoomRequest{
		Members: []*records.Member{
			{
				UserUUID: a,
			},
			{
				UserUUID: b,
			},
		},
	}

	s.OpenRoom(createRoomRequest1, apiKey)

	openRoomRes1 := s.ReadOpenRoomResponse(aConn, 2)
	s.ReadOpenRoomResponse(bConn, 2)
	roomUUID1 := openRoomRes1.Room.UUID

	createRoomRequest2 := &requests.CreateRoomRequest{
		Members: []*records.Member{
			{
				UserUUID: a,
			},
			{
				UserUUID: c,
			},
		},
	}

	s.OpenRoom(createRoomRequest2, apiKey)
	openRoomRes2 := s.ReadOpenRoomResponse(cConn, 2)
	s.ReadOpenRoomResponse(aConn, 2)
	roomUUID2 := openRoomRes2.Room.UUID

	s.SendMessages(aClient.UserUUID, aClient.DeviceUUID, roomUUID1, aConn, aValidToken)
	s.SendMessages(bClient.UserUUID, bClient.DeviceUUID, roomUUID1, bConn, bValidToken)

	s.RecvMessages(bConn)
	s.RecvMessages(aConn)

	aRooms := s.MakeGetRoomsByUserUUIDRequest(a, 0, apiKey)
	s.Len(aRooms.Rooms, 2)
	bRooms := s.MakeGetRoomsByUserUUIDRequest(b, 0, apiKey)
	s.Len(bRooms.Rooms, 1)

	messagesResp := s.MakeGetMessagesByRoomUUIDRequest(roomUUID1, apiKey, 0)
	s.Len(messagesResp.Messages, 20)
	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID1, apiKey, 20)
	s.Len(messagesResp.Messages, 20)
	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID1, apiKey, 40)
	s.Len(messagesResp.Messages, 10)

	// send messages between A and C
	s.SendMessages(aClient.UserUUID, aClient.DeviceUUID, roomUUID2, aConn, aValidToken)
	s.SendMessages(cClient.UserUUID, cClient.DeviceUUID, roomUUID2, cConn, cValidToken)

	s.RecvMessages(cConn)
	s.RecvMessages(aConn)

	aRooms = s.MakeGetRoomsByUserUUIDRequest(a, 0, apiKey)
	s.Len(aRooms.Rooms, 2)
	cRooms := s.MakeGetRoomsByUserUUIDRequest(c, 0, apiKey)
	s.Len(cRooms.Rooms, 1)

	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID2, apiKey, 0)
	s.Len(messagesResp.Messages, 20)
	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID2, apiKey, 20)
	s.Len(messagesResp.Messages, 20)
	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID2, apiKey, 40)
	s.Len(messagesResp.Messages, 10)

	// create room between A and D
	createRoomReq3 := &requests.CreateRoomRequest{
		Members: []*records.Member{
			{
				UserUUID: a,
			},
			{
				UserUUID: d,
			},
		},
	}

	s.OpenRoom(createRoomReq3, apiKey)

	openRoomRes3 := s.ReadOpenRoomResponse(dConn, 2)
	s.ReadOpenRoomResponse(aConn, 2)
	roomUUID3 := openRoomRes3.Room.UUID

	// send messages between A and D
	s.SendMessages(aClient.UserUUID, aClient.DeviceUUID, roomUUID3, aConn, aValidToken)
	s.SendMessages(dClient.UserUUID, dClient.DeviceUUID, roomUUID3, dConn, dValidToken)

	s.RecvMessages(aConn)
	s.RecvMessages(dConn)

	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID3, apiKey, 0)
	s.Len(messagesResp.Messages, 20)
	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID3, apiKey, 20)
	s.Len(messagesResp.Messages, 20)
	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID3, apiKey, 40)
	s.Len(messagesResp.Messages, 10)

	// create room between B and C
	openRoomReq4 := &requests.CreateRoomRequest{
		Members: []*records.Member{
			{
				UserUUID: b,
			},
			{
				UserUUID: c,
			},
		},
	}

	s.OpenRoom(openRoomReq4, apiKey)
	openRoomRes4 := s.ReadOpenRoomResponse(bConn, 2)
	s.ReadOpenRoomResponse(cConn, 2)
	roomUUID4 := openRoomRes4.Room.UUID

	// send messages between B and C
	s.SendMessages(bClient.UserUUID, bClient.DeviceUUID, roomUUID4, bConn, bValidToken)
	s.SendMessages(cClient.UserUUID, cClient.DeviceUUID, roomUUID4, cConn, cValidToken)

	s.RecvMessages(bConn)
	s.RecvMessages(cConn)

	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID4, apiKey, 0)
	s.Len(messagesResp.Messages, 20)
	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID4, apiKey, 20)
	s.Len(messagesResp.Messages, 20)
	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID4, apiKey, 40)
	s.Len(messagesResp.Messages, 10)

	// create room between B and D
	openRoomRequest5 := &requests.CreateRoomRequest{
		Members: []*records.Member{
			{
				UserUUID: bClient.UserUUID,
			},
			{
				UserUUID: dClient.UserUUID,
			},
		},
	}

	s.OpenRoom(openRoomRequest5, apiKey)

	openRoomRes5 := s.ReadOpenRoomResponse(dConn, 2)
	s.ReadOpenRoomResponse(bConn, 2)
	roomUUID5 := openRoomRes5.Room.UUID

	// send messages between B and D
	s.SendMessages(bClient.UserUUID, bClient.DeviceUUID, roomUUID5, bConn, bValidToken)
	s.SendMessages(dClient.UserUUID, dClient.DeviceUUID, roomUUID5, dConn, dValidToken)

	s.RecvMessages(bConn)
	s.RecvMessages(dConn)

	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID5, apiKey, 0)
	s.Len(messagesResp.Messages, 20)
	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID5, apiKey, 20)
	s.Len(messagesResp.Messages, 20)
	messagesResp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID5, apiKey, 40)
	s.Len(messagesResp.Messages, 10)

}

func (s *IntegrationTestSuite) TestAllConnectionsRcvMessages() {
	a := uuid.New().String()
	b := uuid.New().String()
	token := s.GetValidToken(a)
	apiKey := s.GetValidAPIKey()

	aClient, aConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     token,
		UserUUID:  a,
	})

	bClient, bConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     token,
		UserUUID:  b,
	})

	_, bMobileConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     token,
		UserUUID:  b,
	})

	openRoomEvent := &requests.CreateRoomRequest{
		Members: []*records.Member{
			{
				UserUUID: a,
			},
			{
				UserUUID: b,
			},
		},
	}

	s.OpenRoom(openRoomEvent, apiKey)

	openRoomRes := s.ReadOpenRoomResponse(aConn, 2)
	s.ReadOpenRoomResponse(bConn, 2)
	s.ReadOpenRoomResponse(bMobileConn, 2)
	roomUUID := openRoomRes.Room.UUID

	s.SendMessages(aClient.UserUUID, aClient.DeviceUUID, roomUUID, aConn, token)
	s.SendMessages(bClient.UserUUID, bClient.DeviceUUID, roomUUID, bConn, token)

	s.RecvMessages(bConn)
	s.RecvMessages(aConn)

	// need to recv double the msgs
	s.RecvMessages(bMobileConn)
	s.RecvMessages(bMobileConn)

	resp := s.MakeGetMessagesByRoomUUIDRequest(roomUUID, apiKey, 0)
	s.Len(resp.Messages, 20)
	resp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID, apiKey, 20)
	s.Len(resp.Messages, 20)
	resp = s.MakeGetMessagesByRoomUUIDRequest(roomUUID, apiKey, 40)
	s.Len(resp.Messages, 10)

	// add new connection
	_, aMobileConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     token,
		UserUUID:  a,
	})

	s.SendMessages(aClient.UserUUID, aClient.DeviceUUID, roomUUID, aConn, token)
	s.SendMessages(bClient.UserUUID, bClient.DeviceUUID, roomUUID, bConn, token)

	s.RecvMessages(bConn)
	s.RecvMessages(aConn)

	// need to recv double the msgs
	s.RecvMessages(bMobileConn)
	s.RecvMessages(bMobileConn)

	// need to recv double the msgs
	s.RecvMessages(aMobileConn)
	s.RecvMessages(aMobileConn)
}
