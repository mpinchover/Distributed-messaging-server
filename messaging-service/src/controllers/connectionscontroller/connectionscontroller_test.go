package connectionscontroller

import (
	"messaging-service/src/types/requests"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
)

type ConnectionsControllerSuite struct {
	suite.Suite
	cxnsController *ConnectionsController
}

// this function executes before the test suite begins execution
func (s *ConnectionsControllerSuite) SetupSuite() {
	var mu sync.Mutex
	s.cxnsController = &ConnectionsController{
		Mu:   &mu,
		Cxns: map[string]*requests.Connection{},
	}
}

func TestChannelsControllerSuite(t *testing.T) {
	suite.Run(t, new(ConnectionsControllerSuite))
}

func (s *ConnectionsControllerSuite) TestAddConnection() {
	cn1 := &requests.Connection{
		UserUUID: "uuid-1",
		Connections: map[string]*websocket.Conn{
			"client-1": &websocket.Conn{},
		},
	}
	cn2 := &requests.Connection{
		UserUUID: "uuid-2",
		Connections: map[string]*websocket.Conn{
			"client-2": &websocket.Conn{},
			"client-3": &websocket.Conn{},
		},
	}

	s.cxnsController.AddConnection(cn1)
	s.cxnsController.AddConnection(cn2)
	cn := s.cxnsController.GetConnection("uuid-3")
	s.Nil(cn)

	cn = s.cxnsController.GetConnection("uuid-1")
	s.NotNil(cn)
	s.Equal("uuid-1", cn.UserUUID)
	s.Len(cn.Connections, 1)

	cn = s.cxnsController.GetConnection("uuid-2")
	s.NotNil(cn)
	s.Equal("uuid-2", cn.UserUUID)
	s.Len(cn.Connections, 2)

	cn = s.cxnsController.GetConnection("uuid-1")
	s.cxnsController.AddClient(cn, "client-2", &websocket.Conn{})
	cn = s.cxnsController.GetConnection("uuid-1")
	s.Equal("uuid-1", cn.UserUUID)
	s.Len(cn.Connections, 2)

	s.cxnsController.DelClient("uuid-1", "client-1")
	cn = s.cxnsController.GetConnection("uuid-1")
	s.Equal("uuid-1", cn.UserUUID)
	s.Len(cn.Connections, 1)

	s.cxnsController.DelClient("uuid-1", "client-2")
	cn = s.cxnsController.GetConnection("uuid-1")
	s.Nil(cn)
}
