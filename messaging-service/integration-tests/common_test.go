package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"messaging-service/types/enums"
	"messaging-service/types/requests"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func sendMessages(t *testing.T, fromUserUUID string, connectionUUID string, roomUUID string, conn *websocket.Conn) {
	for i := 0; i < 25; i++ {
		msgText := fmt.Sprintf("Message %d", i)
		msgEventOut := &requests.TextMessageEvent{
			FromUUID:       fromUserUUID,
			ConnectionUUID: connectionUUID,
			EventType:      enums.EVENT_TEXT_MESSAGE.String(),
			Message: &requests.Message{
				RoomUUID:    roomUUID,
				MessageText: msgText,
			},
		}
		sendTextMessage(t, conn, msgEventOut)
	}
}

func recvMessages(t *testing.T, conn *websocket.Conn) {
	for i := 0; i < 25; i++ {
		// conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		resp := &requests.TextMessageEvent{}
		recvMessage(t, conn, resp)
	}
}

func recvMessage(t *testing.T, conn *websocket.Conn, resp *requests.TextMessageEvent) {
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)
	err = json.Unmarshal(p, resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.EventType)
	assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType)
	assert.NotEmpty(t, resp.FromUUID)
	assert.NotEmpty(t, resp.ConnectionUUID)
	assert.NotEmpty(t, resp.Message.RoomUUID)
	assert.NotEmpty(t, resp.Message.MessageText)
}

func containsRoomUUID(s []*requests.Room, str string) bool {
	for _, v := range s {
		if v.UUID == str {
			return true
		}
	}

	return false
}

func readOpenRoomResponse(t *testing.T, conn *websocket.Conn, expectedMembers int) *requests.OpenRoomEvent {
	// TODO - ensure correct users are in the room
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)
	resp := &requests.OpenRoomEvent{}
	err = json.Unmarshal(p, resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.EventType)
	assert.Equal(t, enums.EVENT_OPEN_ROOM.String(), resp.EventType)
	assert.NotEmpty(t, resp.Room)
	assert.NotEmpty(t, resp.Room.UUID)
	assert.Equal(t, expectedMembers, len(resp.Room.Members))

	for _, m := range resp.Room.Members {
		assert.NotEmpty(t, m.UUID)
		assert.NotEmpty(t, m.UserUUID)
	}

	return resp
}

func sendTextMessage(t *testing.T, ws *websocket.Conn, msgEvent *requests.TextMessageEvent) {
	err := ws.WriteJSON(msgEvent)
	assert.NoError(t, err)
}

func getMessagesByRoomUUID(t *testing.T, roomUUID string, offset int) (*requests.GetMessagesByRoomUUIDResponse, error) {
	url := fmt.Sprintf("http://localhost:9090/get-messages-by-room-uuid?roomUuid=%s&offset=%d", roomUUID, offset)
	resp, err := http.Get(url)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		fmt.Println(string(b))
		return nil, fmt.Errorf("error code is %d", resp.StatusCode)
	}

	result := &requests.GetMessagesByRoomUUIDResponse{}
	err = json.Unmarshal(b, result)
	assert.NoError(t, err)
	return result, err
}

func getRoomsByUserUUID(userUUID string, offset int) (*requests.GetRoomsByUserUUIDResponse, error) {
	url := fmt.Sprintf("http://localhost:9090/get-rooms-by-user-uuid?userUuid=%s&offset=%d", userUUID, offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		fmt.Println(string(b))
		return nil, fmt.Errorf("error code is %d", resp.StatusCode)
	}

	result := &requests.GetRoomsByUserUUIDResponse{}
	err = json.Unmarshal(b, result)
	return result, err
}

// set up a client connection
func setupClientConnection(t *testing.T, userUUID string) (*requests.SetClientConnectionEvent, *websocket.Conn) {
	conn, _, err := websocket.DefaultDialer.Dial(SocketURL, nil)
	assert.NoError(t, err)

	msgOut := requests.SetClientConnectionEvent{
		FromUUID:  userUUID,
		EventType: enums.EVENT_SET_CLIENT_SOCKET.String(),
	}

	err = conn.WriteJSON(msgOut)
	assert.NoError(t, err)

	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)

	rsp := &requests.SetClientConnectionEvent{}
	err = json.Unmarshal(p, &rsp)
	assert.NoError(t, err)
	assert.NotEmpty(t, rsp.ConnectionUUID)
	assert.NotEmpty(t, rsp.FromUUID)
	return rsp, conn
}

func openRoom(openRoomEvent *requests.CreateRoomRequest) error {
	postBody, err := json.Marshal(openRoomEvent)
	if err != nil {
		return err
	}
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:9090/create-room", "application/json", reqBody)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return fmt.Errorf("status code is %d", resp.StatusCode)
	}
	return nil

}

func deleteRoom(t *testing.T, deleteRoomRequest *requests.DeleteRoomRequest) {
	postBody, err := json.Marshal(deleteRoomRequest)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:9090/delete-room", "application/json", reqBody)
	assert.NoError(t, err)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		fmt.Println(string(b))
	}

	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.LessOrEqual(t, resp.StatusCode, 299)
}

func leaveRoom(t *testing.T, req *requests.LeaveRoomRequest) {
	postBody, err := json.Marshal(req)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:9090/leave-room", "application/json", reqBody)
	assert.NoError(t, err)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		fmt.Println(string(b))
	}

	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.LessOrEqual(t, resp.StatusCode, 299)
}

func readEvent(conn *websocket.Conn, v interface{}) error {
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	_, p, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	err = json.Unmarshal(p, v)
	return err
}

func queryMessages(t *testing.T, userUUID string, roomUUID string, expectedRooms int) {
	totalMessages := []*requests.Message{}
	res, err := getRoomsByUserUUID(userUUID, 0)

	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	assert.Equal(t, expectedRooms, len(res.Rooms))

	// ensure it contains the room uuid
	assert.True(t, containsRoomUUID(res.Rooms, roomUUID))

	resp, err := getMessagesByRoomUUID(t, roomUUID, 0)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(resp.Messages))

	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, 20, len(totalMessages))

	resp, err = getMessagesByRoomUUID(t, roomUUID, len(totalMessages))
	assert.NoError(t, err)
	assert.Equal(t, 20, len(resp.Messages))
	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, 40, len(totalMessages))

	resp, err = getMessagesByRoomUUID(t, roomUUID, len(totalMessages))
	assert.NoError(t, err)
	assert.Equal(t, 10, len(resp.Messages))
	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, 50, len(totalMessages))

	// jump by 15 because the msgs are being sent too fast.
	for i := 15; i < len(totalMessages); i++ {
		// prev := totalMessages[i-1]
		// cur := totalMessages[i]
		// assert.True(t, prev.CreatedAt >= cur.CreatedAt)
	}

}

func recvSeenMessageEvent(t *testing.T, conn *websocket.Conn, messageUUID string) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)
	seenMessageEvent := &requests.SeenMessageEvent{}
	err = json.Unmarshal(p, seenMessageEvent)
	assert.NoError(t, err)

	assert.NotEmpty(t, seenMessageEvent.EventType)
	assert.Equal(t, enums.EVENT_SEEN_MESSAGE.String(), seenMessageEvent.EventType)
	assert.NotEmpty(t, seenMessageEvent.MessageUUID)
	assert.Equal(t, messageUUID, seenMessageEvent.MessageUUID)
	assert.NotEmpty(t, seenMessageEvent.RoomUUID)
	assert.NotEmpty(t, seenMessageEvent.UserUUID)
}
