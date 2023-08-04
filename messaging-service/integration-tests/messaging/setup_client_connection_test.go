package integrationtests

import (
	"messaging-service/integration-tests/common"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOpenSocket(t *testing.T) {
	// t.Skip()
	t.Run("test set sign up user and setup client", func(t *testing.T) {
		// log.Printf("Running %s", t.Name())
		t.Parallel()
		t.Logf("Runningg test %s at %d", t.Name(), time.Now().UnixNano())

		// create a user
		newUserResponse := common.CreateRandomUser(t)

		apiKey := common.MakeGetAPIKeyRequest(t, newUserResponse.AccessToken)

		generateMessagingTokenRequest := &requests.GenerateMessagingTokenRequest{
			UserID: newUserResponse.UUID,
		}
		generateMessagingTokenResp := common.MakeGenerateMessagingTokenRequest(t, generateMessagingTokenRequest, apiKey.Key)
		assert.NotEmpty(t, generateMessagingTokenResp.Token)

		// set up the client
		setupClientConnEvent := &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			UserUUID:  newUserResponse.UUID,
			Token:     generateMessagingTokenResp.Token,
		}
		setupClientConnResp, conn := common.CreateClientConnection(t, setupClientConnEvent)

		assert.NotNil(t, setupClientConnResp)
		assert.NotEmpty(t, setupClientConnResp.ConnectionUUID)
		assert.NotEmpty(t, setupClientConnResp.UserUUID)
		assert.Equal(t, setupClientConnResp.UserUUID, newUserResponse.UUID)

		pingHandler := conn.PingHandler()
		err := pingHandler("PING")
		assert.NoError(t, err)

		_, p, err := conn.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, "PONG", string(p))
	})
}

func TestCloseSocketConnection(t *testing.T) {

	t.Run("test set sign up user and setup client", func(t *testing.T) {
		// log.Printf("Running %s", t.Name())
		t.Parallel()
		t.Logf("Runningg test %s at %d", t.Name(), time.Now().UnixNano())

		// create a user
		newUserResponse := common.CreateRandomUser(t)

		apiKey := common.MakeGetAPIKeyRequest(t, newUserResponse.AccessToken)

		generateMessagingTokenRequest := &requests.GenerateMessagingTokenRequest{
			UserID: newUserResponse.UUID,
		}
		generateMessagingTokenResp := common.MakeGenerateMessagingTokenRequest(t, generateMessagingTokenRequest, apiKey.Key)
		assert.NotEmpty(t, generateMessagingTokenResp.Token)

		clientTom, tom := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			UserUUID:  uuid.New().String() + "_51",
			Token:     generateMessagingTokenResp.Token,
		})
		clientJerry, jerry := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			UserUUID:  uuid.New().String() + "_52",
			Token:     generateMessagingTokenResp.Token,
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
		// fmt.Println("CREATING ROOM 20")
		common.OpenRoom(t, openRoomEvent, apiKey.Key)
		time.Sleep(2 * time.Second)
		tom.Close()
		jerry.Close()

		common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			UserUUID:  clientTom.UserUUID,
			Token:     generateMessagingTokenResp.Token,
		})

		common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
			UserUUID:  clientJerry.UserUUID,
			Token:     generateMessagingTokenResp.Token,
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
		// fmt.Println("CREATING ROOM 21")
		common.OpenRoom(t, openRoomEvent, apiKey.Key)
		time.Sleep(2 * time.Second)
	})
}
