package integrationtests

import (
	"log"
	"messaging-service/integration-tests/common"
	"messaging-service/types/enums"
	"messaging-service/types/requests"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRoomAndMessagesPagination(t *testing.T) {
	// t.Skip()
	t.Run("test rooms and messages pagination", func(t *testing.T) {
		log.Printf("Running %s", t.Name())

		// need to get valid API key as well
		validMessagingToken, validAPIKey := common.GetValidToken(t)

		aClient, aConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		bClient, bConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		cClient, cConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		dClient, dConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		createRoomRequest1 := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aClient.UserUUID,
				},
				{
					UserUUID: bClient.UserUUID,
				},
			},
		}
		common.OpenRoom(t, createRoomRequest1, validAPIKey)

		openRoomRes1 := common.ReadOpenRoomResponse(t, aConn, 2)
		openRoomRes1 = common.ReadOpenRoomResponse(t, bConn, 2)
		roomUUID1 := openRoomRes1.Room.UUID

		createRoomRequest2 := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aClient.UserUUID,
				},
				{
					UserUUID: cClient.UserUUID,
				},
			},
		}
		common.OpenRoom(t, createRoomRequest2, validAPIKey)
		openRoomRes2 := common.ReadOpenRoomResponse(t, cConn, 2)
		openRoomRes2 = common.ReadOpenRoomResponse(t, aConn, 2)
		roomUUID2 := openRoomRes2.Room.UUID

		// send messages between A and B
		common.SendMessages(t, aClient.UserUUID, aClient.ConnectionUUID, roomUUID1, aConn, validMessagingToken)
		common.SendMessages(t, bClient.UserUUID, bClient.ConnectionUUID, roomUUID1, bConn, validMessagingToken)

		common.RecvMessages(t, bConn)
		common.RecvMessages(t, aConn)

		time.Sleep(1 * time.Second)
		common.QueryMessagesByMessagingJWT(t, aClient.UserUUID, roomUUID1, 2, validMessagingToken)
		common.QueryMessagesByMessagingJWT(t, bClient.UserUUID, roomUUID1, 1, validMessagingToken)

		// send messages between A and C
		common.SendMessages(t, aClient.UserUUID, aClient.ConnectionUUID, roomUUID2, aConn, validMessagingToken)
		common.SendMessages(t, cClient.UserUUID, cClient.ConnectionUUID, roomUUID2, cConn, validMessagingToken)

		common.RecvMessages(t, cConn)
		common.RecvMessages(t, aConn)

		time.Sleep(1 * time.Second)
		common.QueryMessagesByMessagingJWT(t, aClient.UserUUID, roomUUID2, 2, validMessagingToken)
		common.QueryMessagesByMessagingJWT(t, cClient.UserUUID, roomUUID2, 1, validMessagingToken)

		// create room between A and D
		createRoomReq3 := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aClient.UserUUID,
				},
				{
					UserUUID: dClient.UserUUID,
				},
			},
		}
		common.OpenRoom(t, createRoomReq3, validAPIKey)

		openRoomRes3 := common.ReadOpenRoomResponse(t, dConn, 2)
		openRoomRes3 = common.ReadOpenRoomResponse(t, aConn, 2)
		roomUUID3 := openRoomRes3.Room.UUID

		// send messages between A and D
		common.SendMessages(t, aClient.UserUUID, aClient.ConnectionUUID, roomUUID3, aConn, validMessagingToken)
		common.SendMessages(t, dClient.UserUUID, dClient.ConnectionUUID, roomUUID3, dConn, validMessagingToken)

		common.RecvMessages(t, aConn)
		common.RecvMessages(t, dConn)

		time.Sleep(1 * time.Second)
		common.QueryMessagesByMessagingJWT(t, aClient.UserUUID, roomUUID3, 3, validMessagingToken)
		common.QueryMessagesByMessagingJWT(t, dClient.UserUUID, roomUUID3, 1, validMessagingToken)

		// create room between B and C
		openRoomReq4 := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: bClient.UserUUID,
				},
				{
					UserUUID: cClient.UserUUID,
				},
			},
		}

		common.OpenRoom(t, openRoomReq4, validAPIKey)
		openRoomRes4 := common.ReadOpenRoomResponse(t, bConn, 2)
		openRoomRes4 = common.ReadOpenRoomResponse(t, cConn, 2)
		roomUUID4 := openRoomRes4.Room.UUID

		// send messages between B and C
		common.SendMessages(t, bClient.UserUUID, bClient.ConnectionUUID, roomUUID4, bConn, validMessagingToken)
		common.SendMessages(t, cClient.UserUUID, cClient.ConnectionUUID, roomUUID4, cConn, validMessagingToken)

		common.RecvMessages(t, bConn)
		common.RecvMessages(t, cConn)

		time.Sleep(1 * time.Second)
		common.QueryMessagesByMessagingJWT(t, bClient.UserUUID, roomUUID4, 2, validMessagingToken)
		common.QueryMessagesByMessagingJWT(t, cClient.UserUUID, roomUUID4, 2, validMessagingToken)

		// create room between B and D
		openRoomRequest5 := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: bClient.UserUUID,
				},
				{
					UserUUID: dClient.UserUUID,
				},
			},
		}
		common.OpenRoom(t, openRoomRequest5, validAPIKey)

		openRoomRes5 := common.ReadOpenRoomResponse(t, dConn, 2)
		openRoomRes5 = common.ReadOpenRoomResponse(t, bConn, 2)

		// the mobiel device will get the open room msg as well
		roomUUID5 := openRoomRes5.Room.UUID

		// send messages between B and D
		common.SendMessages(t, bClient.UserUUID, bClient.ConnectionUUID, roomUUID5, bConn, validMessagingToken)
		common.SendMessages(t, dClient.UserUUID, dClient.ConnectionUUID, roomUUID5, dConn, validMessagingToken)

		common.RecvMessages(t, bConn)
		common.RecvMessages(t, dConn)

		time.Sleep(100 * time.Millisecond)
		common.QueryMessagesByMessagingJWT(t, bClient.UserUUID, roomUUID5, 3, validMessagingToken)
		common.QueryMessagesByMessagingJWT(t, dClient.UserUUID, roomUUID5, 2, validMessagingToken)
	})
}

// Need to get the room id first and pass it to the text message id
func TestAllConnectionsRcvMessages(t *testing.T) {
	// t.Skip()
	t.Run("test all connections get msgs", func(t *testing.T) {
		log.Printf("Running test %s", t.Name())

		// need to get valid API key as well
		validMessagingToken, validAPIKey := common.GetValidToken(t)

		aClient, aConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		bClient, bConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		_, bMobileConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  bClient.UserUUID,
		})

		openRoomEvent := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aClient.UserUUID,
				},
				{
					UserUUID: bClient.UserUUID,
				},
			},
		}
		common.OpenRoom(t, openRoomEvent, validAPIKey)

		openRoomRes := common.ReadOpenRoomResponse(t, aConn, 2)
		common.ReadOpenRoomResponse(t, bConn, 2)
		common.ReadOpenRoomResponse(t, bMobileConn, 2)
		roomUUID := openRoomRes.Room.UUID

		common.SendMessages(t, aClient.UserUUID, aClient.ConnectionUUID, roomUUID, aConn, validMessagingToken)
		common.SendMessages(t, bClient.UserUUID, bClient.ConnectionUUID, roomUUID, bConn, validMessagingToken)

		common.RecvMessages(t, bConn)
		common.RecvMessages(t, aConn)

		// need to recv double the msgs
		common.RecvMessages(t, bMobileConn)
		common.RecvMessages(t, bMobileConn)
		common.QueryMessagesByMessagingJWT(t, bClient.UserUUID, roomUUID, 1, validMessagingToken)
		common.QueryMessagesByMessagingJWT(t, aClient.UserUUID, roomUUID, 1, validMessagingToken)

		// add new connection

		_, aMobileConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  aClient.UserUUID,
		})

		common.SendMessages(t, aClient.UserUUID, aClient.ConnectionUUID, roomUUID, aConn, validMessagingToken)
		common.SendMessages(t, bClient.UserUUID, bClient.ConnectionUUID, roomUUID, bConn, validMessagingToken)

		common.RecvMessages(t, bConn)
		common.RecvMessages(t, aConn)

		// need to recv double the msgs
		common.RecvMessages(t, bMobileConn)
		common.RecvMessages(t, bMobileConn)

		// need to recv double the msgs
		common.RecvMessages(t, aMobileConn)
		common.RecvMessages(t, aMobileConn)

	})
}
