package integrestion_testing

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestMockEndpoint(t *testing.T) {
	t.Run("mock test", func(t *testing.T) {
		assert.Equal(t, 123, 123, "they should be equal")
	})
	t.Run("mock test", func(t *testing.T) {
		assert.NotEqual(t, 124, 123, "they should be equal")
	})
}

func TestConnectWebsocket(t *testing.T) {
	t.Run("open websocket", func(t *testing.T) {

		url := "ws://localhost:9090/ws"
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		assert.NoError(t, err)

		msg := map[string]interface{}{
			"age": 25,
		}
		err = ws.WriteJSON(msg)
		assert.NoError(t, err)
	})
}
