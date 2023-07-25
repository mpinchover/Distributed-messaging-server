package connectionscontroller

import (
	"messaging-service/src/types/requests"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionsController struct {
	Mu   *sync.Mutex
	Cxns map[string]*requests.Connection
}

func New() *ConnectionsController {
	var mu sync.Mutex
	return &ConnectionsController{
		Mu:   &mu,
		Cxns: map[string]*requests.Connection{},
	}
}

func (s *ConnectionsController) GetConnection(userUUID string) *requests.Connection {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	connection, ok := s.Cxns[userUUID]
	if !ok {
		return nil
	}
	return connection
}

// for users
func (s *ConnectionsController) AddConnection(connection *requests.Connection) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Cxns[connection.UserUUID] = connection
}

// for websockets
func (s *ConnectionsController) AddClient(connection *requests.Connection,
	connectionUUID string,
	conn *websocket.Conn) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	connection.Connections[connectionUUID] = conn
}
