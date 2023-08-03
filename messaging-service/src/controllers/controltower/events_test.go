package controltower

import (
	"errors"
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"

	"github.com/stretchr/testify/mock"
)

func (s *ControlTowerSuite) TestProcessTextMessage() {
	tests := []struct {
		test        string
		expectedErr string
		mocks       func()
		event       *requests.TextMessageEvent
	}{
		{
			test:        "GetRoomByRoomUUID fail",
			expectedErr: "database failure",
			mocks: func() {
				s.mockRepoClient.On("GetRoomByRoomUUID", mock.Anything).Return(nil, errors.New("database failure")).Once()
			},
			event: &requests.TextMessageEvent{
				Message: &requests.Message{
					RoomUUID: "room-uuid",
				},
			},
		},
		{
			test:        "room not found",
			expectedErr: "room does not exist",
			mocks: func() {
				s.mockRepoClient.On("GetRoomByRoomUUID", mock.Anything).Return(nil, nil).Once()
			},
			event: &requests.TextMessageEvent{
				Message: &requests.Message{
					RoomUUID: "room-uuid",
				},
			},
		},
		{
			test:        "save message fails",
			expectedErr: "saveMessage failed",
			mocks: func() {
				s.mockRepoClient.On("GetRoomByRoomUUID", mock.Anything).Return(&records.Room{}, nil).Once()
				s.mockRepoClient.On("SaveMessage", mock.Anything).Return(errors.New("saveMessage failed")).Once()
			},
			event: &requests.TextMessageEvent{
				Message: &requests.Message{
					RoomUUID: "room-uuid",
				},
			},
		},
		{
			test:        "publish message fails",
			expectedErr: "publish failed",
			mocks: func() {
				s.mockRepoClient.On("GetRoomByRoomUUID", mock.Anything).Return(&records.Room{}, nil).Once()
				s.mockRepoClient.On("SaveMessage", mock.Anything).Return(nil).Once()
				s.mockRedisClient.On("PublishToRedisChannel", mock.Anything, mock.Anything).Return(errors.New("publish failed")).Once()
			},
			event: &requests.TextMessageEvent{
				Message: &requests.Message{
					RoomUUID: "room-uuid",
				},
			},
		},
		{
			test: "success",
			mocks: func() {
				s.mockRepoClient.On("GetRoomByRoomUUID", mock.Anything).Return(&records.Room{}, nil).Once()
				s.mockRepoClient.On("SaveMessage", mock.Anything).Return(nil).Once()
				s.mockRedisClient.On("PublishToRedisChannel", mock.Anything, mock.Anything).Return(nil).Once()
			},
			event: &requests.TextMessageEvent{
				Message: &requests.Message{
					RoomUUID: "room-uuid",
				},
			},
		},
	}

	for _, t := range tests {
		t.mocks()

		res, err := s.controlTower.ProcessTextMessage(t.event)
		if t.expectedErr != "" {
			s.Error(err)
			s.Contains(err.Error(), t.expectedErr)
			s.Nil(res)
		} else {
			s.NoError(err)
			s.NotNil(res)
			s.NotEmpty(res.UUID)
		}
	}
}
