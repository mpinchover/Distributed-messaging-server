package apitests

// func TestGetMessagesByUserUUIDWithAPIKey(t *testing.T) {
// 	// t.Skip()

// 	t.Run("test get messages with api key", func(t *testing.T) {
// 		t.Parallel()
// 		t.Logf("Runningg test %s at %d", t.Name(), time.Now().UnixNano())

// 		// need to get valid API key as well
// 		validMessagingToken, validAPIKey := common.GetValidToken(t)

// 		tom := uuid.New().String() + "_12"
// 		jerry := uuid.New().String() + "_12"

// 		tomClient, tomConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     validMessagingToken,
// 			UserUUID:  tom,
// 		})

// 		jerryClient, jerryConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     validMessagingToken,
// 			UserUUID:  jerry,
// 		})

// 		createRoomRequest := &requests.CreateRoomRequest{
// 			Members: []*requests.Member{
// 				{
// 					UserUUID: tom,
// 				},
// 				{
// 					UserUUID: jerry,
// 				},
// 			},
// 		}

// 		common.OpenRoom(t, createRoomRequest, validAPIKey)

// 		// test tom
// 		roomsResponse := common.GetRoomsByUserUUIDWithApiKey(t, tom, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 1)

// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, jerry, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 1)

// 		roomUUID := roomsResponse.Rooms[0].UUID

// 		// send messages between A and B
// 		// send 150 messages
// 		common.SendMessages(t, jerry, jerryClient.ConnectionUUID, roomUUID, jerryConn, validMessagingToken)
// 		common.SendMessages(t, tom, tomClient.ConnectionUUID, roomUUID, tomConn, validMessagingToken)
// 		common.SendMessages(t, jerry, jerryClient.ConnectionUUID, roomUUID, jerryConn, validMessagingToken)
// 		common.SendMessages(t, tom, tomClient.ConnectionUUID, roomUUID, tomConn, validMessagingToken)
// 		common.SendMessages(t, jerry, jerryClient.ConnectionUUID, roomUUID, jerryConn, validMessagingToken)
// 		common.SendMessages(t, tom, tomClient.ConnectionUUID, roomUUID, tomConn, validMessagingToken)

// 		time.Sleep(2 * time.Second)

// 		// ensure that all the messages are there
// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, tom, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 1)

// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, jerry, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 1)

// 		totalMessages := []*requests.Message{}
// 		messagesResponse := common.GetMessagesByRoomUUIDByWithAPIKey(t, roomUUID, 0, validAPIKey)
// 		assert.Len(t, messagesResponse.Messages, 20)
// 		totalMessages = append(totalMessages, messagesResponse.Messages...)
// 		assert.Len(t, totalMessages, 20)

// 		messagesResponse = common.GetMessagesByRoomUUIDByWithAPIKey(t, roomUUID, len(totalMessages), validAPIKey)
// 		assert.Len(t, messagesResponse.Messages, 20)
// 		totalMessages = append(totalMessages, messagesResponse.Messages...)
// 		assert.Len(t, totalMessages, 40)

// 		messagesResponse = common.GetMessagesByRoomUUIDByWithAPIKey(t, roomUUID, len(totalMessages), validAPIKey)
// 		assert.Len(t, messagesResponse.Messages, 20)
// 		totalMessages = append(totalMessages, messagesResponse.Messages...)
// 		assert.Len(t, totalMessages, 60)

// 		messagesResponse = common.GetMessagesByRoomUUIDByWithAPIKey(t, roomUUID, len(totalMessages), validAPIKey)
// 		assert.Len(t, messagesResponse.Messages, 20)
// 		totalMessages = append(totalMessages, messagesResponse.Messages...)
// 		assert.Len(t, totalMessages, 80)

// 		messagesResponse = common.GetMessagesByRoomUUIDByWithAPIKey(t, roomUUID, len(totalMessages), validAPIKey)
// 		assert.Len(t, messagesResponse.Messages, 20)
// 		totalMessages = append(totalMessages, messagesResponse.Messages...)
// 		assert.Len(t, totalMessages, 100)

// 		messagesResponse = common.GetMessagesByRoomUUIDByWithAPIKey(t, roomUUID, len(totalMessages), validAPIKey)
// 		assert.Len(t, messagesResponse.Messages, 20)
// 		totalMessages = append(totalMessages, messagesResponse.Messages...)
// 		assert.Len(t, totalMessages, 120)

// 		messagesResponse = common.GetMessagesByRoomUUIDByWithAPIKey(t, roomUUID, len(totalMessages), validAPIKey)
// 		assert.Len(t, messagesResponse.Messages, 20)
// 		totalMessages = append(totalMessages, messagesResponse.Messages...)
// 		assert.Len(t, totalMessages, 140)

// 		messagesResponse = common.GetMessagesByRoomUUIDByWithAPIKey(t, roomUUID, len(totalMessages), validAPIKey)
// 		assert.Len(t, messagesResponse.Messages, 10)
// 		totalMessages = append(totalMessages, messagesResponse.Messages...)
// 		assert.Len(t, totalMessages, 150)
// 	})
// }
