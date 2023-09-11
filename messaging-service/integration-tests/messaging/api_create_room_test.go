package integrationtests

import (
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"

	"github.com/google/uuid"
)

func (s *IntegrationTestSuite) TestCreateRoom() {
	a := uuid.New().String()
	b := uuid.New().String()

	aToken := s.GetValidToken(a)
	bToken := s.GetValidToken(b)

	apiKey := s.GetValidAPIKey()

	_, aConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     aToken,
		UserUUID:  a,
	})

	_, bConn := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		Token:     bToken,
		UserUUID:  b,
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
		},
	}

	s.OpenRoom(createRoomRequest, apiKey)

	aResp := s.ReadOpenRoomResponse(aConn, 2)
	bResp := s.ReadOpenRoomResponse(bConn, 2)
	s.Equal(aResp.Room.UUID, bResp.Room.UUID)
}
