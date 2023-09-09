package apitests

// func TestCreateRoom(t *testing.T) {
// 	// t.Skip()
// 	t.Run("create room", func(t *testing.T) {
// 		t.Parallel()
// 		t.Logf("Runningg test %s at %d", t.Name(), time.Now().UnixNano())
// 		validMessagingToken, validAPIKey := common.GetValidToken(t)

// 		tomClient, tomConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     validMessagingToken,
// 			UserUUID:  uuid.New().String() + "_7",
// 		})

// 		jerryClient, jerryConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     validMessagingToken,
// 			UserUUID:  uuid.New().String() + "_8",
// 		})

// 		// create a room
// 		createRoomRequest := &requests.CreateRoomRequest{
// 			Members: []*requests.Member{
// 				{
// 					UserUUID: tomClient.UserUUID,
// 				},
// 				{
// 					UserUUID: jerryClient.UserUUID,
// 				},
// 			},
// 		}

// 		common.OpenRoom(t, createRoomRequest, validAPIKey)

// 		tOpenRoomResponse := common.ReadOpenRoomResponse(t, tomConn, 2)
// 		jOpenRoomResponse := common.ReadOpenRoomResponse(t, jerryConn, 2)

// 		assert.Equal(t, tOpenRoomResponse.Room.UUID, jOpenRoomResponse.Room.UUID)

// 	})
// }
