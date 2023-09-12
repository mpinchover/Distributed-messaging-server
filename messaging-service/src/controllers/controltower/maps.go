package controltower

import (
	"errors"
	"messaging-service/src/types/connections"
)

func (c *ControlTowerCtrlr) GetUserConnection(userUUID string) *connections.UserConnection {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	userConn, ok := c.UserConnections[userUUID]
	if !ok {
		return nil
	}
	return userConn
}

func (c *ControlTowerCtrlr) SetUserConnection(cn *connections.UserConnection) error {
	if cn.UUID == "" {
		return errors.New("no user uuid provided")
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.UserConnections[cn.UUID] = cn
	return nil
}

func (c *ControlTowerCtrlr) DeleteUserFromServer(userUUID string) error {
	if userUUID == "" {
		return errors.New("no user uuid provided")
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	delete(c.UserConnections, userUUID)
	return nil
}

func (c *ControlTowerCtrlr) DeleteDeviceFromServer(userUUID string, deviceUUID string) error {
	if userUUID == "" {
		return errors.New("no user uuid provided")
	}

	if deviceUUID == "" {
		return errors.New("no device uuid provided")
	}

	userConn := c.GetUserConnection(userUUID)
	if userConn == nil {
		return nil
	}

	c.Mu.Lock()
	delete(userConn.Devices, deviceUUID)
	c.Mu.Unlock()

	c.SetUserConnection(userConn)
	return nil
}

func (c *ControlTowerCtrlr) SetUserDevice(userUUID string, deviceUUID string, device *connections.Device) error {
	userConn := c.GetUserConnection(userUUID)
	if userConn == nil {
		return errors.New("user conn not found")
	}

	if userConn.Devices == nil {
		userConn.Devices = map[string]*connections.Device{}
	}

	c.Mu.Lock()
	userConn.Devices[deviceUUID] = device
	c.Mu.Unlock()

	err := c.SetUserConnection(userConn)
	if err != nil {
		return err
	}

	return nil
}

func (c *ControlTowerCtrlr) GetUserDevice(userUUID, deviceUUID string) *connections.Device {
	userConn := c.GetUserConnection(userUUID)
	if userConn == nil {
		return nil
	}

	if userConn.Devices == nil {
		return nil
	}

	c.Mu.RLock()
	defer c.Mu.RUnlock()

	return userConn.Devices[deviceUUID]
}

func (c *ControlTowerCtrlr) GetChannelFromServer(chUUID string) *connections.Channel {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	ch, ok := c.Channels[chUUID]
	if !ok {
		return nil
	}
	return ch
}

func (c *ControlTowerCtrlr) GetAllChannelsOnServerForUser(userUUID string) []*connections.Channel {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	channels := []*connections.Channel{}
	for _, ch := range c.Channels {
		if ch.Users == nil {
			continue
		}

		_, ok := ch.Users[userUUID]
		if !ok {
			continue
		}

		channels = append(channels, ch)
	}
	return channels
}

func (c *ControlTowerCtrlr) SetChannelOnServer(chUUID string, ch *connections.Channel) error {
	if chUUID == "" {
		return errors.New("no ch uuid provided")
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.Channels[chUUID] = ch
	return nil
}

func (c *ControlTowerCtrlr) AddUserToChannel(userUUID, chUUID string) error {
	channel := c.GetChannelFromServer(chUUID)
	if channel == nil {
		return errors.New("channel not found")
	}
	if channel.Users == nil {
		channel.Users = map[string]bool{}
	}

	c.Mu.Lock()
	channel.Users[userUUID] = true
	c.Mu.Unlock()

	return c.SetChannelOnServer(chUUID, channel)
}

func (c *ControlTowerCtrlr) DeleteUserFromChannel(userUUID, chUUID string) error {
	channel := c.GetChannelFromServer(chUUID)
	if channel == nil {
		return errors.New("channel not found")
	}

	c.Mu.Lock()
	delete(channel.Users, userUUID)
	c.Mu.Unlock()

	return nil
}

func (c *ControlTowerCtrlr) DeleteChannelFromServer(chUUID string) error {
	if chUUID == "" {
		return errors.New("ch uuid not provided")
	}

	ch := c.GetChannelFromServer(chUUID)
	if ch == nil {
		return nil
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	delete(c.Channels, chUUID)
	return nil
}
