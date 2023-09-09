package apitests

// func TestGetRoomByUUIDWithApiKey(t *testing.T) {

// 	t.Run("test get rooms with api key", func(t *testing.T) {
// 		t.Parallel()
// 		t.Logf("Runningg test %s at %d", t.Name(), time.Now().UnixNano())

// 		// need to get valid API key as well
// 		_, validAPIKey := common.GetValidToken(t)

// 		tom := uuid.New().String() + "_1"
// 		jerry := uuid.New().String() + "_2"

// 		for i := 0; i < 25; i++ {
// 			createRoomRequest := &requests.CreateRoomRequest{
// 				Members: []*requests.Member{
// 					{
// 						UserUUID: tom,
// 					},
// 					{
// 						UserUUID: uuid.New().String(),
// 					},
// 				},
// 			}

// 			common.OpenRoom(t, createRoomRequest, validAPIKey)
// 		}

// 		for i := 0; i < 5; i++ {
// 			createRoomRequest := &requests.CreateRoomRequest{
// 				Members: []*requests.Member{
// 					{
// 						UserUUID: jerry,
// 					},
// 					{
// 						UserUUID: uuid.New().String(),
// 					},
// 				},
// 			}

// 			common.OpenRoom(t, createRoomRequest, validAPIKey)
// 		}

// 		// test tom
// 		totalRooms := []*requests.Room{}
// 		roomsResponse := common.GetRoomsByUserUUIDWithApiKey(t, tom, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		totalRooms = append(totalRooms, roomsResponse.Rooms...)

// 		assert.Len(t, roomsResponse.Rooms, 10)
// 		assert.Len(t, totalRooms, 10)

// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, tom, len(totalRooms), validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		totalRooms = append(totalRooms, roomsResponse.Rooms...)

// 		assert.Len(t, roomsResponse.Rooms, 10)
// 		assert.Len(t, totalRooms, 20)

// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, tom, len(totalRooms), validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		totalRooms = append(totalRooms, roomsResponse.Rooms...)

// 		assert.Len(t, roomsResponse.Rooms, 5)
// 		assert.Len(t, totalRooms, 25)

// 		// test jerry
// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, jerry, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)

// 		assert.Len(t, roomsResponse.Rooms, 5)

// 		// send messages between A and B
// 		// send 150 messages
// 		// common.SendMessages(t, aClient.UserUUID, aClient.ConnectionUUID, roomUUID, aConn, validMessagingToken)
// 		// common.SendMessages(t, bClient.UserUUID, bClient.ConnectionUUID, roomUUID, bConn, validMessagingToken)
// 		// common.SendMessages(t, aClient.UserUUID, aClient.ConnectionUUID, roomUUID, aConn, validMessagingToken)
// 		// common.SendMessages(t, bClient.UserUUID, bClient.ConnectionUUID, roomUUID, bConn, validMessagingToken)
// 		// common.SendMessages(t, aClient.UserUUID, aClient.ConnectionUUID, roomUUID, aConn, validMessagingToken)
// 		// common.SendMessages(t, bClient.UserUUID, bClient.ConnectionUUID, roomUUID, bConn, validMessagingToken)

// 		// time.Sleep(1 * time.Second)
// 		// common.GetRoomsByUserUUIDWithApiKey(t, aClient.UserUUID, roomUUID, 2, validAPIKey)
// 	})
// 	t.Run("test get rooms in correct order with api key", func(t *testing.T) {
// 		t.Parallel()
// 		t.Logf("Runningg test %s at %d", t.Name(), time.Now().UnixNano())

// 		// need to get valid API key as well
// 		validMessagingToken, validAPIKey := common.GetValidToken(t)

// 		tom := uuid.New().String() + "_3"
// 		jerry := uuid.New().String() + "_4"
// 		alice := uuid.New().String() + "_5"
// 		dave := uuid.New().String() + "_6"

// 		clientTom, tomConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     validMessagingToken,
// 			UserUUID:  tom,
// 		})

// 		_, jerryConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     validMessagingToken,
// 			UserUUID:  jerry,
// 		})

// 		clientAlice, aliceConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     validMessagingToken,
// 			UserUUID:  alice,
// 		})

// 		clientDave, daveConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			Token:     validMessagingToken,
// 			UserUUID:  dave,
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
// 		common.ReadOpenRoomResponse(t, jerryConn, 2)
// 		openRoomEvent := common.ReadOpenRoomResponse(t, tomConn, 2)
// 		roomUUID1 := openRoomEvent.Room.UUID
// 		_ = roomUUID1

// 		createRoomRequest = &requests.CreateRoomRequest{
// 			Members: []*requests.Member{
// 				{
// 					UserUUID: tom,
// 				},
// 				{
// 					UserUUID: alice,
// 				},
// 			},
// 		}

// 		common.OpenRoom(t, createRoomRequest, validAPIKey)
// 		common.ReadOpenRoomResponse(t, aliceConn, 2)
// 		openRoomEvent = common.ReadOpenRoomResponse(t, tomConn, 2)
// 		roomUUID2 := openRoomEvent.Room.UUID

// 		createRoomRequest = &requests.CreateRoomRequest{
// 			Members: []*requests.Member{
// 				{
// 					UserUUID: tom,
// 				},
// 				{
// 					UserUUID: dave,
// 				},
// 			},
// 		}

// 		common.OpenRoom(t, createRoomRequest, validAPIKey)
// 		common.ReadOpenRoomResponse(t, daveConn, 2)
// 		openRoomEvent = common.ReadOpenRoomResponse(t, tomConn, 2)
// 		roomUUID3 := openRoomEvent.Room.UUID
// 		// fmt.Println("Room 1", roomUUID1)
// 		// fmt.Println("Room 2", roomUUID2)
// 		// fmt.Println("Room 3", roomUUID3)

// 		roomsResponse := common.GetRoomsByUserUUIDWithApiKey(t, tom, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 3)

// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, jerry, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 1)

// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, alice, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 1)

// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, dave, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 1)

// 		common.SendSingleTextMessage(t, tom, clientTom.ConnectionUUID, roomUUID3, tomConn, validMessagingToken)
// 		time.Sleep(time.Second)
// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, tom, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 3)
// 		assert.Equal(t, roomUUID3, roomsResponse.Rooms[0].UUID)

// 		common.SendSingleTextMessage(t, tom, clientTom.ConnectionUUID, roomUUID2, tomConn, validMessagingToken)
// 		time.Sleep(time.Second)
// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, tom, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 3)
// 		assert.Equal(t, roomUUID2, roomsResponse.Rooms[0].UUID)
// 		assert.Equal(t, roomUUID3, roomsResponse.Rooms[1].UUID)

// 		common.SendSingleTextMessage(t, tom, clientTom.ConnectionUUID, roomUUID1, tomConn, validMessagingToken)
// 		time.Sleep(time.Second)
// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, tom, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 3)
// 		assert.Equal(t, roomUUID1, roomsResponse.Rooms[0].UUID)
// 		assert.Equal(t, roomUUID2, roomsResponse.Rooms[1].UUID)
// 		assert.Equal(t, roomUUID3, roomsResponse.Rooms[2].UUID)

// 		common.SendSingleTextMessage(t, tom, clientTom.ConnectionUUID, roomUUID3, tomConn, validMessagingToken)
// 		time.Sleep(time.Second)
// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, tom, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 3)
// 		assert.Equal(t, roomUUID3, roomsResponse.Rooms[0].UUID)
// 		assert.Equal(t, roomUUID1, roomsResponse.Rooms[1].UUID)
// 		assert.Equal(t, roomUUID2, roomsResponse.Rooms[2].UUID)

// 		common.SendSingleTextMessage(t, dave, clientDave.ConnectionUUID, roomUUID3, daveConn, validMessagingToken)
// 		time.Sleep(time.Second)
// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, tom, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 3)
// 		assert.Equal(t, roomUUID3, roomsResponse.Rooms[0].UUID)
// 		assert.Equal(t, roomUUID1, roomsResponse.Rooms[1].UUID)
// 		assert.Equal(t, roomUUID2, roomsResponse.Rooms[2].UUID)

// 		common.SendSingleTextMessage(t, alice, clientAlice.ConnectionUUID, roomUUID2, aliceConn, validMessagingToken)
// 		time.Sleep(time.Second)
// 		roomsResponse = common.GetRoomsByUserUUIDWithApiKey(t, tom, 0, validAPIKey)
// 		assert.NotNil(t, roomsResponse)
// 		assert.Len(t, roomsResponse.Rooms, 3)
// 		assert.Equal(t, roomUUID2, roomsResponse.Rooms[0].UUID)
// 		assert.Equal(t, roomUUID3, roomsResponse.Rooms[1].UUID)
// 		assert.Equal(t, roomUUID1, roomsResponse.Rooms[2].UUID)

// 	})
// }
