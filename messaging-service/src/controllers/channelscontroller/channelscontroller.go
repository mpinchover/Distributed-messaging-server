package channelscontroller

import (
	"errors"
	"messaging-service/src/types/requests"
	"sync"
)

type ChannelsController struct {
	Mu    *sync.Mutex
	Chnls map[string]*requests.ServerChannel
}

func New() *ChannelsController {
	var mu sync.Mutex
	return &ChannelsController{
		Mu:    &mu,
		Chnls: map[string]*requests.ServerChannel{},
	}
}

func (s *ChannelsController) GetChannel(chUUID string) *requests.ServerChannel {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	ch, ok := s.Chnls[chUUID]
	if !ok {
		return nil
	}
	return ch
}

func (s *ChannelsController) SetChannel(ch *requests.ServerChannel) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	chUUID := ch.UUID
	s.Chnls[chUUID] = ch
}

func (s *ChannelsController) DeleteUser(roomUUID string, userUUID string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	serverChannel, ok := s.Chnls[roomUUID]
	if !ok {
		return errors.New("server channel does not exist")
	}

	delete(serverChannel.MembersOnServer, userUUID)
	return nil
}

func (s *ChannelsController) DeleteChannel(roomUUID string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	delete(s.Chnls, roomUUID)
}

func (s *ChannelsController) AddChannel(ch *requests.ServerChannel) *requests.ServerChannel {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Chnls[ch.UUID] = ch
	return ch
}
