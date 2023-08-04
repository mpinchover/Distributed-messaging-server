package channelscontroller

import (
	"messaging-service/src/types/requests"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ChannelsControllerSuite struct {
	suite.Suite
	channelsCtrl *ChannelsController
}

// this function executes before the test suite begins execution
func (s *ChannelsControllerSuite) SetupSuite() {
	var mu sync.Mutex
	s.channelsCtrl = &ChannelsController{
		Mu:    &mu,
		Chnls: map[string]*requests.ServerChannel{},
	}
}

func TestChannelsControllerSuite(t *testing.T) {
	suite.Run(t, new(ChannelsControllerSuite))
}

func (s *ChannelsControllerSuite) TestGetAndSetAndDeleteChannel() {
	ch1 := &requests.ServerChannel{
		MembersOnServer: map[string]bool{
			"user-1": true,
			"user-2": true,
		},
		UUID: "room-uuid-ch-1",
	}
	ch2 := &requests.ServerChannel{
		MembersOnServer: map[string]bool{
			"user-3": true,
			"user-4": true,
		},
		UUID: "room-uuid-ch-2",
	}

	s.channelsCtrl.SetChannel(ch1)
	s.channelsCtrl.SetChannel(ch2)
	ch := s.channelsCtrl.GetChannel("room-uuid-ch-1")
	s.NotNil(ch)
	s.Equal("room-uuid-ch-1", ch.UUID)
	s.Len(ch.MembersOnServer, 2)

	ch = s.channelsCtrl.GetChannel("null")
	s.Nil(ch)
	ch = s.channelsCtrl.GetChannel("room-uuid-ch-2")
	s.NotNil(ch)
	s.Equal("room-uuid-ch-2", ch.UUID)
	s.Len(ch.MembersOnServer, 2)

	ch = s.channelsCtrl.GetChannel("room-uuid-ch-2")
	ch.MembersOnServer["user-2"] = true
	numChannels := s.channelsCtrl.GetChannelsByUserUUID("user-2")
	s.Len(numChannels, 2)

	numChannels = s.channelsCtrl.GetChannelsByUserUUID("user-1")
	s.Len(numChannels, 1)

	numChannels = s.channelsCtrl.GetChannelsByUserUUID("user-3")
	s.Len(numChannels, 1)

	s.channelsCtrl.DeleteChannel("room-uuid-ch-1")
	ch = s.channelsCtrl.GetChannel("room-uuid-ch-1")
	s.Nil(ch)

	numChannels = s.channelsCtrl.GetChannelsByUserUUID("user-1")
	s.Len(numChannels, 0)

	ch = s.channelsCtrl.GetChannel("room-uuid-ch-2")
	s.NotNil(ch)
	s.channelsCtrl.DeleteChannel("room-uuid-ch-2")
	ch = s.channelsCtrl.GetChannel("room-uuid-ch-2")
	s.Nil(ch)

	numChannels = s.channelsCtrl.GetChannelsByUserUUID("user-2")
	s.Len(numChannels, 0)
}

func (s *ChannelsControllerSuite) TestAddChannel() {
	ch1 := &requests.ServerChannel{
		MembersOnServer: map[string]bool{
			"user-1": true,
			"user-2": true,
		},
		UUID: "room-uuid-ch-1",
	}
	s.channelsCtrl.AddChannel(ch1)
	ch := s.channelsCtrl.GetChannel("room-uuid-ch-1")
	s.NotNil(ch)
	s.Equal("room-uuid-ch-1", ch.UUID)
	s.Len(ch.MembersOnServer, 2)

	ch.MembersOnServer["user-3"] = true
	s.channelsCtrl.SetChannel(ch)
	ch = s.channelsCtrl.GetChannel("room-uuid-ch-1")
	s.NotNil(ch)
	s.Equal("room-uuid-ch-1", ch.UUID)
	s.Len(ch.MembersOnServer, 3)
}

func (s *ChannelsControllerSuite) TestDeleteUser() {
	ch1 := &requests.ServerChannel{
		MembersOnServer: map[string]bool{
			"user-1": true,
			"user-2": true,
		},
		UUID: "room-uuid-ch-1",
	}
	s.channelsCtrl.AddChannel(ch1)
	s.channelsCtrl.DeleteUser("room-uuid-ch-1", "user-1")
	ch := s.channelsCtrl.GetChannel("room-uuid-ch-1")
	s.NotNil(ch)
	s.Equal("room-uuid-ch-1", ch.UUID)
	s.Len(ch.MembersOnServer, 1)
}
