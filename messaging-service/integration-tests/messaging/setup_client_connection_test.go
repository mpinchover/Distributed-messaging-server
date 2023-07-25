package integrationtests

import (
	"log"
	"messaging-service/integration-tests/common"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenSocket(t *testing.T) {
	// t.Skip()
	t.Run("test set sign up user and setup client", func(t *testing.T) {
		log.Printf("Running %s", t.Name())

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
