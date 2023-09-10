package integrationtests

// func TestSeenBy(t *testing.T) {
// 	// t.Skip()
// 	t.Run("delete room and messages", func(t *testing.T) {
// 		t.Parallel()

// 		tomUUID := uuid.New().String()
// 		jerryUUID := uuid.New().String()
// 		aliceUUID := uuid.New().String()
// 		deanUUID := uuid.New().String()

// 		tomToken := common.GetValidToken(t, tomUUID)
// 		aliceToken := common.GetValidToken(t, aliceUUID)
// 		jerryToken := common.GetValidToken(t, jerryUUID)
// 		deanToken := common.GetValidToken(t, deanUUID)

// 		apiKey := common.GetValidAPIKey(t)

// 		tomClient, tomConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     tomToken,
// 			UserUUID:  tomUUID,
// 		})

// 		aliceClient, aliceConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     aliceToken,
// 			UserUUID:  uuid.New().String(),
// 		})

// 		jerryClient, jerryConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     jerryToken,
// 			UserUUID:  jerryUUID,
// 		})

// 		deanClient, deanConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     deanToken,
// 			UserUUID:  deanUUID,
// 		})

// 		_, deanMobileConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     deanToken,
// 			UserUUID:  deanUUID,
// 		})

// 		// create a room
// 		createRoomRequest := &requests.CreateRoomRequest{
// 			Members: []*requests.Member{
// 				{
// 					UserUUID: tomUUID,
// 				},
// 				{
// 					UserUUID: jerryUUID,
// 				},
// 				{
// 					UserUUID: aliceUUID,
// 				},
// 				{
// 					UserUUID: deanUUID,
// 				},
// 			},
// 		}

// 		common.OpenRoom(t, createRoomRequest, apiKey)
// 		openRoomResponse := common.ReadOpenRoomResponse(t, tomConn, 4)
// 		common.ReadOpenRoomResponse(t, jerryConn, 4)
// 		common.ReadOpenRoomResponse(t, aliceConn, 4)
// 		common.ReadOpenRoomResponse(t, deanConn, 4)
// 		common.ReadOpenRoomResponse(t, deanMobileConn, 4)
// 		roomUUID := openRoomResponse.Room.UUID

// 		// check mappings

// 		// send out a message tom -> room

// 		msgEventOut := &requests.TextMessageEvent{
// 			FromUUID:   tomClient.UserUUID,
// 			DeviceUUID: tomClient.DeviceUUID,
// 			EventType:  enums.EVENT_TEXT_MESSAGE.String(),
// 			Message: &requests.Message{
// 				MessageText: "TEXT",
// 				RoomUUID:    roomUUID,
// 			},
// 			Token: tomToken,
// 		}
// 		common.SendTextMessage(t, tomConn, msgEventOut)
// 		time.Sleep(time.Second)
// 		// everyone should recv the message

// 		// clear out recv msg
// 		resp := &requests.TextMessageEvent{}
// 		common.RecvMessage(t, jerryConn, resp)
// 		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
// 		common.RecvMessage(t, aliceConn, resp)
// 		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
// 		common.RecvMessage(t, deanConn, resp)
// 		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
// 		common.RecvMessage(t, deanMobileConn, resp)
// 		assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)

// 		// send seen event jerry -> room
// 		seenEvent := &requests.SeenMessageEvent{
// 			EventType:   enums.EVENT_SEEN_MESSAGE.String(),
// 			MessageUUID: resp.Message.UUID,
// 			UserUUID:    jerryClient.UserUUID,
// 			RoomUUID:    roomUUID,
// 			Token:       jerryToken,
// 		}

// 		err := jerryConn.WriteJSON(seenEvent)
// 		assert.NoError(t, err)

// 		// everyone should get the seen event message
// 		common.RecvSeenMessageEvent(t, tomConn, resp.Message.UUID)
// 		common.RecvSeenMessageEvent(t, aliceConn, resp.Message.UUID)
// 		common.RecvSeenMessageEvent(t, deanConn, resp.Message.UUID)
// 		common.RecvSeenMessageEvent(t, deanMobileConn, resp.Message.UUID)

// 		// should be a websocket event
// 		res, err := common.GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, 0, validMessagingToken)

// 		assert.NoError(t, err)
// 		assert.Len(t, res.Messages, 1)
// 		assert.Len(t, res.Messages[0].SeenBy, 1)

// 		assert.Equal(t, res.Messages[0].SeenBy[0].MessageUUID, resp.Message.UUID)
// 		assert.Equal(t, res.Messages[0].SeenBy[0].UserUUID, jerryClient.UserUUID)

// 	})
// }
