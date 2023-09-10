package integrationtests

import (
	"encoding/json"
	"io/ioutil"
	"messaging-service/integration-tests/common"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIPing(t *testing.T) {
	t.Parallel()

	requestURL := "http://localhost:9090/ping"
	res, err := http.Get(requestURL)
	assert.NoError(t, err)

	bytes, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	resp := struct {
		Message string
	}{}

	err = json.Unmarshal(bytes, &resp)
	assert.NoError(t, err)
	assert.Equal(t, "pong", resp.Message)
}

func TestOpenSocket(t *testing.T) {
	// t.Skip()
	t.Run("test set sign up user and setup client", func(t *testing.T) {
		t.Parallel()

		// get token
		newUser := uuid.New().String()
		token := common.GetValidToken(t, newUser)
		assert.NotEmpty(t, token)

		// set up the client
		setupClientConnEvent := &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			UserUUID:  newUser,
			Token:     token,
		}
		setupClientConnResp, conn := common.CreateClientConnection(t, setupClientConnEvent)

		assert.NotNil(t, setupClientConnResp, t.Name())
		assert.NotEmpty(t, setupClientConnResp.DeviceUUID, t.Name())
		assert.NotEmpty(t, setupClientConnResp.UserUUID, t.Name())
		assert.Equal(t, setupClientConnResp.UserUUID, newUser, t.Name())

		pingHandler := conn.PingHandler()
		err := pingHandler("PING")
		assert.NoError(t, err, t.Name())

		_, p, err := conn.ReadMessage()
		assert.NoError(t, err, t.Name())
		assert.Equal(t, "PONG", string(p), t.Name())

	})
}

func TestSocketConnection(t *testing.T) {
	t.Run("test set sign up user and setup client", func(t *testing.T) {
		t.Parallel()

		tomUUID := uuid.New().String()
		tomToken := common.GetValidToken(t, tomUUID)
		assert.NotEmpty(t, tomToken)

		jerryUUID := uuid.New().String()
		jerryToken := common.GetValidToken(t, jerryUUID)
		assert.NotEmpty(t, jerryToken)

		apiKey := common.GetValidAPIKey(t)

		clientTom, tom := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			UserUUID:  tomUUID,
			Token:     tomToken,
		})

		clientJerry, jerry := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			UserUUID:  jerryUUID,
			Token:     jerryToken,
		})

		openRoomEvent := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: clientTom.UserUUID,
				},
				{
					UserUUID: clientJerry.UserUUID,
				},
			},
		}

		room := common.OpenRoom(t, openRoomEvent, apiKey)
		assert.NotNil(t, room)
		assert.NotNil(t, room.Room)
		assert.NotEmpty(t, room.Room.UUID)

		roomUUID := room.Room.UUID
		time.Sleep(2 * time.Second)

		// test mappings
		userConnectionTom := common.MakeGetUserConnectionRequest(t, tomUUID)
		assert.NotNil(t, userConnectionTom.UserConnection)
		assert.Len(t, userConnectionTom.UserConnection.Devices, 1)
		assert.Equal(t, tomUUID, userConnectionTom.UserConnection.UUID)

		userConnectionJerry := common.MakeGetUserConnectionRequest(t, jerryUUID)
		assert.NotNil(t, userConnectionJerry.UserConnection)
		assert.Len(t, userConnectionJerry.UserConnection.Devices, 1)
		assert.Equal(t, jerryUUID, userConnectionJerry.UserConnection.UUID)

		channel := common.MakeGetChannelConnectionRequest(t, roomUUID)
		assert.NotNil(t, channel)
		assert.Len(t, channel.Users, 2)

		// tom leaves
		tom.Close()

		time.Sleep(2 * time.Second)

		// test mappings
		userConnectionTom = common.MakeGetUserConnectionRequest(t, tomUUID)
		assert.Nil(t, userConnectionTom.UserConnection)

		userConnectionJerry = common.MakeGetUserConnectionRequest(t, jerryUUID)
		assert.NotNil(t, userConnectionJerry.UserConnection)
		assert.Len(t, userConnectionJerry.UserConnection.Devices, 1)
		assert.Equal(t, jerryUUID, userConnectionJerry.UserConnection.UUID)

		channel = common.MakeGetChannelConnectionRequest(t, roomUUID)
		assert.NotNil(t, channel)
		assert.Len(t, channel.Users, 1)
		assert.True(t, channel.Users[jerryUUID])
		assert.False(t, channel.Users[tomUUID])

		// tom comes back
		clientTom, tom = common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			UserUUID:  tomUUID,
			Token:     tomToken,
		})

		// time.Sleep(2 * time.Second)

		// test mappings
		userConnectionTom = common.MakeGetUserConnectionRequest(t, tomUUID)
		assert.NotNil(t, userConnectionTom.UserConnection)
		assert.Len(t, userConnectionTom.UserConnection.Devices, 1)
		assert.Equal(t, tomUUID, userConnectionTom.UserConnection.UUID)

		userConnectionJerry = common.MakeGetUserConnectionRequest(t, jerryUUID)
		assert.NotNil(t, userConnectionJerry.UserConnection)
		assert.Len(t, userConnectionJerry.UserConnection.Devices, 1)
		assert.Equal(t, jerryUUID, userConnectionJerry.UserConnection.UUID)

		channel = common.MakeGetChannelConnectionRequest(t, roomUUID)
		assert.NotNil(t, channel)
		assert.Len(t, channel.Users, 2)
		assert.True(t, channel.Users[jerryUUID])
		assert.True(t, channel.Users[tomUUID])

		// tom goes away again
		tom.Close()
		jerry.Close()

		// test mappings
		userConnectionTom = common.MakeGetUserConnectionRequest(t, tomUUID)
		assert.Nil(t, userConnectionTom.UserConnection)
		userConnectionJerry = common.MakeGetUserConnectionRequest(t, jerryUUID)
		assert.Nil(t, userConnectionJerry.UserConnection)

		channel = common.MakeGetChannelConnectionRequest(t, roomUUID)
		assert.NotNil(t, channel)
		assert.Len(t, channel.Users, 0)
		assert.False(t, channel.Users[jerryUUID])
		assert.False(t, channel.Users[tomUUID])

		time.Sleep(2 * time.Second)

		tomUUID = uuid.New().String()
		tomToken = common.GetValidToken(t, tomUUID)
		assert.NotEmpty(t, tomToken)

		jerryUUID = uuid.New().String()
		jerryToken = common.GetValidToken(t, jerryUUID)
		assert.NotEmpty(t, jerryToken)

		common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			UserUUID:  clientTom.UserUUID,
			Token:     tomToken,
		})

		common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			UserUUID:  clientJerry.UserUUID,
			Token:     jerryToken,
		})

		openRoomEvent = &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: clientTom.UserUUID,
				},
				{
					UserUUID: clientJerry.UserUUID,
				},
			},
		}

		common.OpenRoom(t, openRoomEvent, apiKey)
		time.Sleep(2 * time.Second)

		jerry.Close()
		tom.Close()
	})
}
