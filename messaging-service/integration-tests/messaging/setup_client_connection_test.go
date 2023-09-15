package integrationtests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"
	"net/http"

	"github.com/google/uuid"
)

func (s *IntegrationTestSuite) TestAPIPing() {

	requestURL := fmt.Sprintf("http://%s:9090/ping", ServerHost)
	res, err := http.Get(requestURL)
	s.NoError(err)

	bytes, err := ioutil.ReadAll(res.Body)
	s.NoError(err)

	resp := struct {
		Message string
	}{}

	err = json.Unmarshal(bytes, &resp)
	s.NoError(err)
	s.Equal("pong", resp.Message)
}

func (s *IntegrationTestSuite) TestOpenSocket() {

	// get token
	newUser := uuid.New().String()
	token := s.GetValidToken(newUser)
	s.NotEmpty(s.T(), token)

	// set up the client
	setupClientConnEvent := &requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		UserUUID:  newUser,
		Token:     token,
	}
	setupClientConnResp, conn := s.CreateClientConnection(setupClientConnEvent)

	s.NotNil(setupClientConnResp)
	s.NotEmpty(setupClientConnResp.DeviceUUID)
	s.NotEmpty(setupClientConnResp.UserUUID)
	s.Equal(setupClientConnResp.UserUUID, newUser)

	pingHandler := conn.PingHandler()
	err := pingHandler("PING")
	s.NoError(err)

	_, p, err := conn.ReadMessage()
	s.NoError(err)
	s.Equal("PONG", string(p))
}

func (s *IntegrationTestSuite) TestSocketConnection() {

	tomUUID := uuid.New().String()
	tomToken := s.GetValidToken(tomUUID)
	s.NotEmpty(tomToken)

	jerryUUID := uuid.New().String()
	jerryToken := s.GetValidToken(jerryUUID)
	s.NotEmpty(jerryToken)

	apiKey := s.GetValidAPIKey()

	clientTom, tom := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		UserUUID:  tomUUID,
		Token:     tomToken,
	})

	clientJerry, jerry := s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		UserUUID:  jerryUUID,
		Token:     jerryToken,
	})

	openRoomEvent := &requests.CreateRoomRequest{
		Members: []*requests.Member{
			{
				UserUUID: clientTom.UserUUID,
			},
			{
				UserUUID: clientJerry.UserUUID,
			},
		},
	}

	room := s.OpenRoom(openRoomEvent, apiKey)
	s.NotNil(room)
	s.NotNil(room.Room)
	s.NotEmpty(room.Room.UUID)

	// tom leaves
	tom.Close()
	// tom comes back
	clientTom, tom = s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		UserUUID:  tomUUID,
		Token:     tomToken,
	})

	// tom goes away again
	tom.Close()
	jerry.Close()

	tomUUID = uuid.New().String()
	tomToken = s.GetValidToken(tomUUID)
	s.NotEmpty(tomToken)

	jerryUUID = uuid.New().String()
	jerryToken = s.GetValidToken(jerryUUID)
	s.NotEmpty(jerryToken)

	_, tom = s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		UserUUID:  clientTom.UserUUID,
		Token:     tomToken,
	})

	_, jerry = s.CreateClientConnection(&requests.SetClientConnectionEvent{
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
		UserUUID:  clientJerry.UserUUID,
		Token:     jerryToken,
	})

	openRoomEvent = &requests.CreateRoomRequest{
		Members: []*requests.Member{
			{
				UserUUID: clientTom.UserUUID,
			},
			{
				UserUUID: clientJerry.UserUUID,
			},
		},
	}

	s.OpenRoom(openRoomEvent, apiKey)

	jerry.Close()
	tom.Close()
}
