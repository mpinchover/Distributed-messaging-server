package integrationtests

import (
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"

	"github.com/google/uuid"
)

func (s *IntegrationTestSuite) TestDeleteRoom() {

	a := uuid.New().String()
	b := uuid.New().String()
	c := uuid.New().String()

	atoken := s.GetValidToken(a)
	btoken := s.GetValidToken(b)
	ctoken := s.GetValidToken(c)

	apiKey := s.GetValidAPIKey()

	_, aConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     atoken,
		UserUUID:  a,
	})

	_, bConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     btoken,
		UserUUID:  b,
	})

	_, cConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     ctoken,
		UserUUID:  c,
	})
	openRoomEvent := &requests.CreateRoomRequest{
		Members: []*requests.Member{
			{
				UserUUID: a,
			},
			{
				UserUUID: b,
			},
		},
	}

	s.OpenRoom(openRoomEvent, apiKey)
	s.ReadOpenRoomResponse(aConn, 2)
	openRoomRes := s.ReadOpenRoomResponse(bConn, 2)
	roomUUID1 := openRoomRes.Room.UUID

	openRoomEvent = &requests.CreateRoomRequest{
		Members: []*requests.Member{
			{
				UserUUID: a,
			},
			{
				UserUUID: c,
			},
		},
	}

	s.OpenRoom(openRoomEvent, apiKey)
	s.ReadOpenRoomResponse(aConn, 2)
	openRoomRes = s.ReadOpenRoomResponse(cConn, 2)
	roomUUID2 := openRoomRes.Room.UUID

	res := s.MakeGetRoomsByUserUUIDRequest(a, 0, apiKey)
	s.Equal(2, len(res.Rooms))
	s.Equal(2, len(res.Rooms[0].Members))
	s.Equal(2, len(res.Rooms[1].Members))

	res = s.MakeGetRoomsByUserUUIDRequest(b, 0, apiKey)
	s.Equal(1, len(res.Rooms))
	s.Equal(2, len(res.Rooms[0].Members))

	res = s.MakeGetRoomsByUserUUIDRequest(c, 0, apiKey)
	s.Equal(1, len(res.Rooms))
	s.Equal(2, len(res.Rooms[0].Members))

	s.DeleteRoom(&requests.DeleteRoomRequest{
		RoomUUID: roomUUID1,
	}, apiKey)

	res = s.MakeGetRoomsByUserUUIDRequest(a, 0, apiKey)
	s.Equal(1, len(res.Rooms))
	s.Equal(2, len(res.Rooms[0].Members))

	res = s.MakeGetRoomsByUserUUIDRequest(b, 0, apiKey)
	s.NotEmpty(res)
	s.Equal(0, len(res.Rooms))

	res = s.MakeGetRoomsByUserUUIDRequest(c, 0, apiKey)
	s.NotEmpty(res)
	s.Equal(1, len(res.Rooms))
	s.Equal(2, len(res.Rooms[0].Members))

	// ensure delete event is recd
	resp := &requests.DeleteRoomEvent{}
	s.ReadEvent(aConn, resp)
	s.Equal(enums.EVENT_DELETE_ROOM.String(), resp.EventType)
	s.Equal(roomUUID1, resp.RoomUUID)

	resp = &requests.DeleteRoomEvent{}
	s.ReadEvent(bConn, resp)
	s.Equal(enums.EVENT_DELETE_ROOM.String(), resp.EventType)
	s.Equal(roomUUID1, resp.RoomUUID)

	s.DeleteRoom(&requests.DeleteRoomRequest{
		RoomUUID: roomUUID2,
	}, apiKey)

	res = s.MakeGetRoomsByUserUUIDRequest(a, 0, apiKey)
	s.NotEmpty(res)
	s.Equal(0, len(res.Rooms))

	res = s.MakeGetRoomsByUserUUIDRequest(b, 0, apiKey)
	s.NotEmpty(res)
	s.Equal(0, len(res.Rooms))

	res = s.MakeGetRoomsByUserUUIDRequest(c, 0, apiKey)
	s.NotEmpty(res)
	s.Equal(0, len(res.Rooms))

	// ensure delete event is recd
	resp = &requests.DeleteRoomEvent{}
	s.ReadEvent(aConn, resp)
	s.Equal(enums.EVENT_DELETE_ROOM.String(), resp.EventType)
	s.Equal(roomUUID2, resp.RoomUUID)

	resp = &requests.DeleteRoomEvent{}
	s.ReadEvent(cConn, resp)
	s.Equal(enums.EVENT_DELETE_ROOM.String(), resp.EventType)
	s.Equal(roomUUID2, resp.RoomUUID)

}

func (s *IntegrationTestSuite) TestDeleteRoomAndMessages() {

	apiKey := s.GetValidAPIKey()

	a := uuid.New().String()
	b := uuid.New().String()
	c := uuid.New().String()
	d := uuid.New().String()

	atoken := s.GetValidToken(a)
	btoken := s.GetValidToken(b)
	ctoken := s.GetValidToken(c)
	dtoken_1 := s.GetValidToken(d)
	dtoken_2 := s.GetValidToken(d)
	dtoken_3 := s.GetValidToken(d)

	aClient, aConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     atoken,
		UserUUID:  a,
	})

	_, bConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     btoken,
		UserUUID:  b,
	})

	_, cConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     ctoken,
		UserUUID:  c,
	})

	_, dConnOne := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     dtoken_1,
		UserUUID:  d,
	})

	_, dConnTwo := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     dtoken_2,
		UserUUID:  d,
	})

	_, dConnThree := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     dtoken_3,
		UserUUID:  d,
	})

	// create a room
	createRoomRequest := &requests.CreateRoomRequest{
		Members: []*requests.Member{
			{
				UserUUID: a,
			},
			{
				UserUUID: b,
			},
			{
				UserUUID: c,
			},
			{
				UserUUID: d,
			},
		},
	}

	s.OpenRoom(createRoomRequest, apiKey)

	// get open room response over socket
	resp := s.ReadOpenRoomResponse(aConn, 4)

	s.ReadOpenRoomResponse(bConn, 4)
	s.ReadOpenRoomResponse(cConn, 4)
	s.ReadOpenRoomResponse(dConnOne, 4)
	s.ReadOpenRoomResponse(dConnTwo, 4)
	s.ReadOpenRoomResponse(dConnThree, 4)

	roomUUID := resp.Room.UUID

	s.SendMessages(a, aClient.DeviceUUID, roomUUID, aConn, atoken)
	s.RecvMessages(bConn)
	s.RecvMessages(cConn)
	s.RecvMessages(dConnOne)
	s.RecvMessages(dConnTwo)
	s.RecvMessages(dConnThree)

	res := s.MakeGetMessagesByRoomUUIDRequest(roomUUID, apiKey, 0)
	s.Len(res.Messages, 20)

	deleteRoomRequest := &requests.DeleteRoomRequest{
		RoomUUID: roomUUID,
	}

	s.DeleteRoom(deleteRoomRequest, apiKey)
	res = s.MakeGetMessagesByRoomUUIDRequest(roomUUID, apiKey, 0)
	s.Len(res.Messages, 0)

	// ensure everyone got the deletedRoom event
	e := &requests.DeleteRoomEvent{}
	s.RecvDeletedRoomMsg(aConn, e)
	s.RecvDeletedRoomMsg(bConn, e)
	s.RecvDeletedRoomMsg(cConn, e)
	s.RecvDeletedRoomMsg(dConnOne, e)
	s.RecvDeletedRoomMsg(dConnTwo, e)
	s.RecvDeletedRoomMsg(dConnThree, e)
}
