package controltower

import (
	"context"
	"errors"
	"testing"

	mockRedis "messaging-service/mocks/src/redis"
	mockRepo "messaging-service/mocks/src/repo"
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ControlTowerSuite struct {
	suite.Suite

	mockRepoClient  *mockRepo.RepoInterface
	mockRedisClient *mockRedis.RedisInterface
	// mockConnectionsCtrlr *mockConnections.ConnectionsControllerInterface
	controlTower *ControlTowerCtrlr
}

// this function executes before the test suite begins execution
func (s *ControlTowerSuite) SetupSuite() {

	s.mockRepoClient = mockRepo.NewRepoInterface(s.T())
	s.mockRedisClient = mockRedis.NewRedisInterface(s.T())
	// s.mockConnectionsCtrlr = mockConnections.NewConnectionsControllerInterface(s.T())
	s.controlTower = &ControlTowerCtrlr{
		Repo:        s.mockRepoClient,
		RedisClient: s.mockRedisClient,
		// ConnCtrlr:   s.mockConnectionsCtrlr,
	}

}

func TestControlTowerSuite(t *testing.T) {
	suite.Run(t, new(ControlTowerSuite))
}

func (s *ControlTowerSuite) TestCreateRoom() {
	tests := []struct {
		test        string
		members     []*records.Member
		expectedErr string
		mocks       func()
	}{
		{
			test:        "SaveRoom failed",
			expectedErr: "database error",
			mocks: func() {
				s.mockRepoClient.On("SaveRoom", mock.Anything).Return(errors.New("database error")).Once()
			},
		},
		{
			test:        "PublishToRedisChannel failed",
			expectedErr: "redis error",
			members: []*records.Member{
				{
					UUID:     "room-uuid",
					UserUUID: "user-uuid-1",
				},
				{
					UUID:     "room-uuid",
					UserUUID: "user-uuid-2",
				},
			},
			mocks: func() {
				s.mockRepoClient.On("SaveRoom", mock.Anything).Return(nil).Once()
				s.mockRedisClient.On("PublishToRedisChannel", mock.Anything, mock.Anything).Return(errors.New("redis error")).Once()
			},
		},
		{
			test: "success",
			members: []*records.Member{
				{
					UUID:     "room-uuid",
					UserUUID: "user-uuid-1",
				},
				{
					UUID:     "room-uuid",
					UserUUID: "user-uuid-2",
				},
			},
			mocks: func() {
				s.mockRepoClient.On("SaveRoom", mock.Anything).Return(nil).Once()
				s.mockRedisClient.On("PublishToRedisChannel", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, t := range tests {
		t.mocks()

		res, err := s.controlTower.CreateRoom(context.Background(), t.members)
		if t.expectedErr != "" {
			s.Error(err, t.test)
			s.Contains(err.Error(), t.expectedErr)
			s.Nil(res, t.test)
		} else {
			s.NoError(err)
			s.NotNil(res)
			s.Len(res.Members, 2, t.test)
			for _, m := range res.Members {
				s.NotEmpty(m.UUID, t.test)
				s.NotEmpty(m.UserUUID, t.test)
			}
		}
	}
}

func (s *ControlTowerSuite) TestUpdateMessage() {

}

func (s *ControlTowerSuite) TestLeaveRoom() {}

func (s *ControlTowerSuite) TestDeleteRoom() {
	tests := []struct {
		test        string
		expectedErr string
		mocks       func()
	}{
		{
			test:        "GetRoomByRoomUUID failed",
			expectedErr: "internal err",
			mocks: func() {
				s.mockRepoClient.On("GetRoomByRoomUUID", mock.Anything).Return(nil, errors.New("internal err")).Once()
			},
		},
		{
			test:        "GetRoomByRoomUUID failed",
			expectedErr: "room not found",
			mocks: func() {
				s.mockRepoClient.On("GetRoomByRoomUUID", mock.Anything).Return(nil, nil).Once()
			},
		},
		{
			test:        "DeleteRoom failed",
			expectedErr: "delete room failed",
			mocks: func() {
				s.mockRepoClient.On("GetRoomByRoomUUID", mock.Anything).Return(&records.Room{
					UUID: "uuid",
					Members: []*records.Member{
						{
							UUID:     "uuid",
							RoomUUID: "room-uuid",
							RoomID:   10,
							UserUUID: "user-uuid",
						},
						{
							UUID:     "uuid-2",
							RoomUUID: "room-uuid",
							RoomID:   10,
							UserUUID: "user-uuid-2",
						},
					},
				}, nil).Once()
				s.mockRepoClient.On("DeleteRoom", mock.Anything).Return(errors.New("delete room failed")).Once()
			},
		},
		{
			test:        "PublishToRedisChannel failed",
			expectedErr: "redis publish failed",
			mocks: func() {
				s.mockRepoClient.On("GetRoomByRoomUUID", mock.Anything).Return(&records.Room{
					UUID: "uuid",
					Members: []*records.Member{
						{
							UUID:     "uuid",
							RoomUUID: "room-uuid",
							RoomID:   10,
							UserUUID: "user-uuid",
						},
						{
							UUID:     "uuid-2",
							RoomUUID: "room-uuid",
							RoomID:   10,
							UserUUID: "user-uuid-2",
						},
					},
				}, nil).Once()
				s.mockRepoClient.On("DeleteRoom", mock.Anything).Return(nil).Once()
				s.mockRedisClient.On("PublishToRedisChannel", mock.Anything, mock.Anything).Return(errors.New("redis publish failed")).Once()
			},
		},
		{
			test: "success",
			mocks: func() {
				s.mockRepoClient.On("GetRoomByRoomUUID", mock.Anything).Return(&records.Room{
					UUID: "uuid",
					Members: []*records.Member{
						{
							UUID:     "uuid",
							RoomUUID: "room-uuid",
							RoomID:   10,
							UserUUID: "user-uuid",
						},
						{
							UUID:     "uuid-2",
							RoomUUID: "room-uuid",
							RoomID:   10,
							UserUUID: "user-uuid-2",
						},
					},
				}, nil).Once()
				s.mockRepoClient.On("DeleteRoom", mock.Anything).Return(nil).Once()
				s.mockRedisClient.On("PublishToRedisChannel", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, t := range tests {
		t.mocks()

		err := s.controlTower.DeleteRoom(context.Background(), "uuid")
		if t.expectedErr != "" {
			s.Error(err, t.test)
			s.Contains(err.Error(), t.expectedErr, t.test)
		} else {
			s.NoError(err, t.test)
		}
	}
}

func (s *ControlTowerSuite) TestsSetupClientConnectionV2() {
	tests := []struct {
		test  string
		req   *requests.SetClientConnectionEvent
		mocks func()
	}{
		{
			test: "create new connection for new user",
			req: &requests.SetClientConnectionEvent{
				UserUUID: "uuid",
			},
			mocks: func() {
				// s.mockConnectionsCtrlr.On("GetConnection", mock.Anything).Return(nil).Once()
				// s.mockConnectionsCtrlr.On("AddConnection", mock.Anything).Return().Once()
				// s.mockConnectionsCtrlr.On("AddClient", mock.Anything, mock.Anything, mock.Anything).Return().Once()
			},
		},
		{
			test: "create new connection for existing user",
			req: &requests.SetClientConnectionEvent{
				UserUUID: "uuid",
			},
			mocks: func() {
				// s.mockConnectionsCtrlr.On("GetConnection", mock.Anything).Return(&records.Connection{
				// 	Connections: map[string]*websocket.Conn{},
				// }).Once()
				// s.mockConnectionsCtrlr.On("AddClient", mock.Anything, mock.Anything, mock.Anything).Return().Once()
			},
		},
	}

	for _, t := range tests {
		t.mocks()

		res, err := s.controlTower.SetupClientConnectionV2(&requests.Websocket{
			Conn: &websocket.Conn{},
		}, t.req)
		s.NoError(err, t.test)
		s.NotNil(res, t.test)
		s.NotEmpty(res.DeviceUUID, t.test)
		s.NotEmpty(res.UserUUID, t.test)
	}
}

func (s *ControlTowerSuite) TestSaveSeenBy() {
	tests := []struct {
		test        string
		req         *requests.SeenMessageEvent
		expectedErr string
		mocks       func()
	}{
		{
			test:        "GetMessageByUUID fails",
			expectedErr: "internal err",
			req: &requests.SeenMessageEvent{
				MessageUUID: "message-uuid",
			},
			mocks: func() {
				s.mockRepoClient.On("GetMessageByUUID", mock.Anything).Return(nil, errors.New("internal err")).Once()
			},
		},
		{
			test:        "GetMessageByUUID returns nil",
			expectedErr: "message not found",
			req: &requests.SeenMessageEvent{
				MessageUUID: "message-uuid",
			},
			mocks: func() {
				s.mockRepoClient.On("GetMessageByUUID", mock.Anything).Return(nil, nil).Once()
			},
		},
		{
			test:        "SaveSeenBy returns err",
			expectedErr: "SaveSeenBy failed",
			req: &requests.SeenMessageEvent{
				MessageUUID: "message-uuid",
			},
			mocks: func() {
				s.mockRepoClient.On("GetMessageByUUID", mock.Anything).Return(&records.Message{}, nil).Once()
				s.mockRepoClient.On("SaveSeenBy", mock.Anything).Return(errors.New("SaveSeenBy failed")).Once()
			},
		},
		{
			test:        "redis publish returns err",
			expectedErr: "publishing error",
			req: &requests.SeenMessageEvent{
				MessageUUID: "message-uuid",
			},
			mocks: func() {
				s.mockRepoClient.On("GetMessageByUUID", mock.Anything).Return(&records.Message{}, nil).Once()
				s.mockRepoClient.On("SaveSeenBy", mock.Anything).Return(nil).Once()
				s.mockRedisClient.On("PublishToRedisChannel", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("publishing error")).Once()
			},
		},
		{
			test: "success",
			req: &requests.SeenMessageEvent{
				MessageUUID: "message-uuid",
			},
			mocks: func() {
				s.mockRepoClient.On("GetMessageByUUID", mock.Anything).Return(&records.Message{}, nil).Once()
				s.mockRepoClient.On("SaveSeenBy", mock.Anything).Return(nil).Once()
				s.mockRedisClient.On("PublishToRedisChannel", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, t := range tests {
		t.mocks()

		err := s.controlTower.SaveSeenBy(t.req)
		if t.expectedErr != "" {
			s.Error(err, t.test)
			s.Contains(err.Error(), t.expectedErr, t.test)
		} else {
			s.NoError(err, t.test)
		}
	}
}

func (s *ControlTowerSuite) GetRoomsByUserUUID() {
	tests := []struct {
		test        string
		expectedErr string
		mocks       func()
	}{
		{
			test:        "GetRoomsByUserUUID fails",
			expectedErr: "database error",
			mocks: func() {
				s.mockRepoClient.On("GetRoomsByUserUUID", mock.Anything, mock.Anything).Return(nil, errors.New("database error")).Once()
			},
		},
		{
			test: "success",
			mocks: func() {
				s.mockRepoClient.On("GetRoomsByUserUUID", mock.Anything, mock.Anything).Return([]*records.Room{
					{
						UUID:          "room-uuid",
						CreatedAtNano: 1000,
						Messages: []*records.Message{
							{
								SeenBy: []*records.SeenBy{
									{},
									{},
								},
							},
							{
								SeenBy: []*records.SeenBy{
									{},
									{},
								},
							},
						},
						Members: []*records.Member{
							{
								UUID:     "member-uuid",
								RoomUUID: "room uuid",
								UserUUID: "user-uuid",
							},
							{
								UUID:     "member-uuid",
								RoomUUID: "room uuid",
								UserUUID: "user-uuid",
							},
						},
					},
				}, nil).Once()
			},
		},
	}

	for _, t := range tests {
		t.mocks()
		res, err := s.controlTower.GetRoomsByUserUUID(context.Background(), "user-uuid", 0)
		if t.expectedErr != "" {
			s.Error(err, t.test)
			s.Nil(res, t.test)
			s.Contains(err.Error(), t.expectedErr, t.test)
		} else {
			s.NoError(err, t.test)
			s.NotNil(res, t.test)
			s.Len(res, 1, t.test)
			s.Len(res[0].Members, 2, t.test)
			s.Len(res[0].Messages, 2, t.test)
			s.Len(res[0].Messages[0].SeenBy, 2, t.test)
			s.Len(res[0].Messages[1].SeenBy, 2, t.test)
		}
	}
}
