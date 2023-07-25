package apitests

import (
	"encoding/json"
	"log"
	"messaging-service/integration-tests/common"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRoom(t *testing.T) {
	// t.Skip()
	t.Run("test delete a room", func(t *testing.T) {
		log.Printf("Running test %s", t.Name())
		validMessagingToken, validAPIKey := common.GetValidToken(t)

		// TODO - problem is here.
		// You're setting the same token for sockets are you are for chat profile
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
		common.ReadOpenRoomResponse(t, aConn, 2)
		openRoomRes := common.ReadOpenRoomResponse(t, bConn, 2)
		roomUUID1 := openRoomRes.Room.UUID

		openRoomEvent = &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aClient.UserUUID,
				},
				{
					UserUUID: cClient.UserUUID,
				},
			},
		}
		common.OpenRoom(t, openRoomEvent, validAPIKey)
		common.ReadOpenRoomResponse(t, aConn, 2)
		openRoomRes = common.ReadOpenRoomResponse(t, cConn, 2)
		roomUUID2 := openRoomRes.Room.UUID

		// TODO - change this to the API endpoint
		res := common.GetRoomsByUserUUIDByMessagingJWT(t, aClient.UserUUID, 0, validMessagingToken)
		assert.Equal(t, 2, len(res.Rooms))
		assert.Equal(t, 2, len(res.Rooms[0].Members))
		assert.Equal(t, 2, len(res.Rooms[1].Members))

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, bClient.UserUUID, 0, validMessagingToken)
		assert.Equal(t, 1, len(res.Rooms))
		assert.Equal(t, 2, len(res.Rooms[0].Members))

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, cClient.UserUUID, 0, validMessagingToken)
		assert.Equal(t, 1, len(res.Rooms))
		assert.Equal(t, 2, len(res.Rooms[0].Members))

		common.DeleteRoom(t, &requests.DeleteRoomRequest{
			RoomUUID: roomUUID1,
		}, validAPIKey)

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, aClient.UserUUID, 0, validMessagingToken)
		assert.Equal(t, 1, len(res.Rooms))
		assert.Equal(t, 2, len(res.Rooms[0].Members))

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, bClient.UserUUID, 0, validMessagingToken)
		assert.NotEmpty(t, res)
		assert.Equal(t, 0, len(res.Rooms))

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, cClient.UserUUID, 0, validMessagingToken)
		assert.NotEmpty(t, res)
		assert.Equal(t, 1, len(res.Rooms))
		assert.Equal(t, 2, len(res.Rooms[0].Members))

		// ensure delete event is recd
		resp := &requests.DeleteRoomEvent{}
		common.ReadEvent(t, aConn, resp)
		assert.Equal(t, enums.EVENT_DELETE_ROOM.String(), resp.EventType)
		assert.Equal(t, roomUUID1, resp.RoomUUID)

		resp = &requests.DeleteRoomEvent{}
		common.ReadEvent(t, bConn, resp)
		assert.Equal(t, enums.EVENT_DELETE_ROOM.String(), resp.EventType)
		assert.Equal(t, roomUUID1, resp.RoomUUID)

		common.DeleteRoom(t, &requests.DeleteRoomRequest{
			RoomUUID: roomUUID2,
		}, validAPIKey)

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, aClient.UserUUID, 0, validMessagingToken)
		assert.NotEmpty(t, res)
		assert.Equal(t, 0, len(res.Rooms))

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, bClient.UserUUID, 0, validMessagingToken)
		assert.NotEmpty(t, res)
		assert.Equal(t, 0, len(res.Rooms))

		res = common.GetRoomsByUserUUIDByMessagingJWT(t, cClient.UserUUID, 0, validMessagingToken)
		assert.NotEmpty(t, res)
		assert.Equal(t, 0, len(res.Rooms))

		// ensure delete event is recd
		resp = &requests.DeleteRoomEvent{}
		common.ReadEvent(t, aConn, resp)
		assert.Equal(t, enums.EVENT_DELETE_ROOM.String(), resp.EventType)
		assert.Equal(t, roomUUID2, resp.RoomUUID)

		resp = &requests.DeleteRoomEvent{}
		common.ReadEvent(t, cConn, resp)
		assert.Equal(t, enums.EVENT_DELETE_ROOM.String(), resp.EventType)
		assert.Equal(t, roomUUID2, resp.RoomUUID)
	})
}

func TestDeleteRoomAndMessages(t *testing.T) {
	// t.Skip()
	t.Run("delete room and messages", func(t *testing.T) {
		log.Printf("Running test %s", t.Name())

		validMessagingToken, validAPIKey := common.GetValidToken(t)
		tomClient, tomConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		jerryClient, jerryConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		aliceClient, aliceConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		benClient, benConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  uuid.New().String(),
		})

		_, benDeviceOneConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  benClient.UserUUID,
		})

		_, benDeviceTwoConn := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			Token:     validMessagingToken,
			UserUUID:  benClient.UserUUID,
		})

		// create a room
		createRoomRequest := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: tomClient.UserUUID,
				},
				{
					UserUUID: jerryClient.UserUUID,
				},
				{
					UserUUID: aliceClient.UserUUID,
				},
				{
					UserUUID: benClient.UserUUID,
				},
			},
		}

		common.OpenRoom(t, createRoomRequest, validAPIKey)

		// // get open room response over socket
		tomOpenRoomEventResponse := common.ReadOpenRoomResponse(t, tomConn, 4)

		common.ReadOpenRoomResponse(t, jerryConn, 4)
		common.ReadOpenRoomResponse(t, aliceConn, 4)
		common.ReadOpenRoomResponse(t, benConn, 4)
		common.ReadOpenRoomResponse(t, benDeviceOneConn, 4)
		common.ReadOpenRoomResponse(t, benDeviceTwoConn, 4)

		roomUUID := tomOpenRoomEventResponse.Room.UUID
		common.SendMessages(t, tomClient.UserUUID, tomClient.ConnectionUUID, roomUUID, tomConn, validMessagingToken)
		common.RecvMessages(t, jerryConn)
		common.RecvMessages(t, aliceConn)
		common.RecvMessages(t, benConn)
		common.RecvMessages(t, benDeviceOneConn)
		common.RecvMessages(t, benDeviceTwoConn)

		res, _ := common.GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, 0, validMessagingToken)
		assert.Len(t, res.Messages, 20)

		deleteRoomRequest := &requests.DeleteRoomRequest{
			RoomUUID: roomUUID,
		}
		common.DeleteRoom(t, deleteRoomRequest, validAPIKey)
		res, _ = common.GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, 0, validMessagingToken)
		assert.Len(t, res.Messages, 0)

		// ensure everyone got the deletedRoom event
		e := &requests.DeleteRoomEvent{}
		recvDeletedRoomMsg(t, tomConn, e)
		recvDeletedRoomMsg(t, jerryConn, e)
		recvDeletedRoomMsg(t, aliceConn, e)
		recvDeletedRoomMsg(t, benConn, e)
		recvDeletedRoomMsg(t, benDeviceOneConn, e)
		recvDeletedRoomMsg(t, benDeviceTwoConn, e)
	})
}

func recvDeletedRoomMsg(t *testing.T, conn *websocket.Conn, resp *requests.DeleteRoomEvent) {
	// conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)
	err = json.Unmarshal(p, resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.EventType)
	assert.Equal(t, enums.EVENT_DELETE_ROOM.String(), resp.EventType)
	assert.NotEmpty(t, resp.RoomUUID)
}
