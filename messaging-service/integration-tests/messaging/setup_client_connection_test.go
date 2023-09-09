package integrationtests

import (
	"encoding/json"
	"io/ioutil"
	"messaging-service/integration-tests/common"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIPing(t *testing.T) {
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
		// t.cLogf("Runningg test %s at %d", t.Name(), time.Now().UnixNano())

		// get token
		newUser := uuid.New().String()
		token := common.GetValidToken(t, newUser)

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

// func TestCloseSocketConnection(t *testing.T) {

// 	t.Run("test set sign up user and setup client", func(t *testing.T) {
// 		// log.Printf("Running %s", t.Name())
// 		t.Parallel()
// 		t.Logf("Runningg test %s at %d", t.Name(), time.Now().UnixNano())

// 		// create a user
// 		newUserResponse := common.CreateRandomUser(t)

// 		apiKey := common.MakeGetAPIKeyRequest(t, newUserResponse.AccessToken)

// 		generateMessagingTokenRequest := &requests.GenerateMessagingTokenRequest{
// 			UserID: newUserResponse.UUID,
// 		}
// 		generateMessagingTokenResp := common.MakeGenerateMessagingTokenRequest(t, generateMessagingTokenRequest, apiKey.Key)
// 		assert.NotEmpty(t, generateMessagingTokenResp.Token)

// 		clientTom, tom := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			UserUUID:  uuid.New().String() + "_51",
// 			Token:     generateMessagingTokenResp.Token,
// 		})
// 		clientJerry, jerry := common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			UserUUID:  uuid.New().String() + "_52",
// 			Token:     generateMessagingTokenResp.Token,
// 		})

// 		openRoomEvent := &requests.CreateRoomRequest{
// 			Members: []*requests.Member{
// 				{
// 					UserUUID: clientTom.UserUUID,
// 				},
// 				{
// 					UserUUID: clientJerry.UserUUID,
// 				},
// 			},
// 		}
// 		// fmt.Println("CREATING ROOM 20")
// 		common.OpenRoom(t, openRoomEvent, apiKey.Key)
// 		time.Sleep(2 * time.Second)
// 		tom.Close()
// 		jerry.Close()

// 		common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			UserUUID:  clientTom.UserUUID,
// 			Token:     generateMessagingTokenResp.Token,
// 		})

// 		common.CreateClientConnection(t, &requests.SetClientConnectionEvent{
// 			EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
// 			UserUUID:  clientJerry.UserUUID,
// 			Token:     generateMessagingTokenResp.Token,
// 		})

// 		openRoomEvent = &requests.CreateRoomRequest{
// 			Members: []*requests.Member{
// 				{
// 					UserUUID: clientTom.UserUUID,
// 				},
// 				{
// 					UserUUID: clientJerry.UserUUID,
// 				},
// 			},
// 		}
// 		// fmt.Println("CREATING ROOM 21")
// 		common.OpenRoom(t, openRoomEvent, apiKey.Key)
// 		time.Sleep(2 * time.Second)
// 	})
// }
