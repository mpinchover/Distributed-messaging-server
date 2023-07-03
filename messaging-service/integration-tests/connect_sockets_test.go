package integrationtests

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestMockEndpoint(t *testing.T) {

	t.Run("mock test", func(t *testing.T) {
		log.Printf("Running %s", t.Name())
		assert.Equal(t, 123, 123, "they should be equal")
	})
	t.Run("mock test", func(t *testing.T) {
		log.Printf("Running %s", t.Name())
		assert.NotEqual(t, 124, 123, "they should be equal")
	})
}

func TestConnectWebsocket(t *testing.T) {
	t.Run("test opening websocket", func(t *testing.T) {
		log.Printf("Running %s", t.Name())

		ws, _, err := websocket.DefaultDialer.Dial(SocketURL, nil)
		assert.NoError(t, err)
		pingHandler := ws.PingHandler()
		err = pingHandler("PING")
		assert.NoError(t, err)

		_, p, err := ws.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, "PONG", string(p))
	})
}

func TestOpenSocket(t *testing.T) {
	t.Run("test set open socket info", func(t *testing.T) {
		log.Printf("Running %s", t.Name())

		clientUUID := uuid.New().String()
		setupClientConnection(t, clientUUID)

	})
}
