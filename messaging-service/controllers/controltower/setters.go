package controltower

import (
	"errors"
	"messaging-service/types/requests"

	"github.com/gorilla/websocket"
)

func (c *ControlTowerController) SetRoomToServer(room *requests.ServerChannel) {
	c.MapLock.Lock()
	defer c.MapLock.Unlock()

	roomUUID := room.UUID
	c.ServerChannels[roomUUID] = room
}

func (c *ControlTowerController) RemoveUserFromServerChannel(roomUUID string, userUUID string) error {
	c.MapLock.Lock()
	defer c.MapLock.Unlock()

	serverChannel, ok := c.ServerChannels[roomUUID]
	if !ok {
		return errors.New("server channel does not exist")
	}

	delete(serverChannel.MembersOnServer, userUUID)
	return nil
}

func (c *ControlTowerController) DeleteServerChannel(roomUUID string) {
	c.MapLock.Lock()
	defer c.MapLock.Unlock()

	delete(c.ServerChannels, roomUUID)
}

func (c *ControlTowerController) AddServerChannel(ch *requests.ServerChannel) *requests.ServerChannel {
	c.MapLock.Lock()
	defer c.MapLock.Unlock()

	c.ServerChannels[ch.UUID] = ch
	return ch
}

func (c *ControlTowerController) AddUserConnection(connection *requests.Connection) {
	c.MapLock.Lock()
	defer c.MapLock.Unlock()

	c.Connections[connection.UserUUID] = connection
}

func (c *ControlTowerController) AddClientConnection(connection *requests.Connection,
	connectionUUID string,
	conn *websocket.Conn) {
	c.MapLock.Lock()
	defer c.MapLock.Unlock()

	connection.Connections[connectionUUID] = conn
}
