package controltower

import "messaging-service/types/requests"

func (c *ControlTowerController) GetClientConnectionFromServer(userUUID string) *requests.Connection {
	c.MapLock.Lock()
	defer c.MapLock.Unlock()

	connection, ok := c.Connections[userUUID]
	if !ok {
		return nil
	}
	return connection
}

func (c *ControlTowerController) GetChannelFromServer(roomUUID string) *requests.ServerChannel {
	c.MapLock.Lock()
	defer c.MapLock.Unlock()

	room, ok := c.ServerChannels[roomUUID]
	if !ok {
		return nil
	}
	return room
}
