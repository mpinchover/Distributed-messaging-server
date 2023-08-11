package channelscontroller

import (
	"messaging-service/src/serrors"
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

func (s *ChannelsController) AddUserToChannel(chUUID string, userUUID string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	ch, ok := s.Chnls[chUUID]
	if !ok {
		serrors.InternalErrorf("channel not found", nil)
	}

	ch.MembersOnServer[userUUID] = true
	return nil
}

func (s *ChannelsController) DeleteUser(roomUUID string, userUUID string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	serverChannel, ok := s.Chnls[roomUUID]
	if !ok {
		serrors.InternalErrorf("channel not found", nil)
	}

	if serverChannel != nil && serverChannel.MembersOnServer != nil {
		delete(serverChannel.MembersOnServer, userUUID)
		if len(serverChannel.MembersOnServer) == 0 {
			delete(s.Chnls, roomUUID)
		}
	}
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

func (s *ChannelsController) GetChannelsByUserUUID(userUUID string) []*requests.ServerChannel {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	res := []*requests.ServerChannel{}

	for _, ch := range s.Chnls {
		if ch.MembersOnServer != nil && ch.MembersOnServer[userUUID] {
			res = append(res, ch)
		}
	}
	return res
}
