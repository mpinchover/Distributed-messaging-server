package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"
	"messaging-service/src/utils"
	"net/http"

	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	// "github.com/stretchr/testify/s"
)

const (
// SocketURL = "ws://%s:9090/ws"
)

var (
	ServerHost = "localhost"
	SocketURL  = fmt.Sprintf("ws://%s:9090/ws", ServerHost)
)

func (s *IntegrationTestSuite) SendSingleTextMessage(fromUserUUID string, deviceUUID string, roomUUID string, conn *websocket.Conn, token string) {
	msgText := "text"
	msgEventOut := &requests.TextMessageEvent{
		FromUUID:   fromUserUUID,
		DeviceUUID: deviceUUID,
		EventType:  enums.EVENT_TEXT_MESSAGE.String(),
		Message: &requests.Message{
			RoomUUID:    roomUUID,
			MessageText: msgText,
		},
		Token: token,
	}
	s.SendTextMessage(conn, msgEventOut)

}

func (s *IntegrationTestSuite) SendMessages(fromUserUUID string, deviceUUID string, roomUUID string, conn *websocket.Conn, token string) {
	for i := 0; i < 25; i++ {
		msgText := fmt.Sprintf("Message %d", i)
		msgEventOut := &requests.TextMessageEvent{
			FromUUID:   fromUserUUID,
			DeviceUUID: deviceUUID,
			EventType:  enums.EVENT_TEXT_MESSAGE.String(),
			Message: &requests.Message{
				RoomUUID:    roomUUID,
				MessageText: msgText,
			},
			Token: token,
		}
		// time.Sleep(time.Millisecond * 500)
		s.SendTextMessage(conn, msgEventOut)
	}
}

func (s *IntegrationTestSuite) RecvMessages(conn *websocket.Conn) {
	for i := 0; i < 25; i++ {
		// conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		resp := &requests.TextMessageEvent{}
		s.RecvMessage(conn, resp)
	}
}

func (s *IntegrationTestSuite) RecvMessage(conn *websocket.Conn, resp *requests.TextMessageEvent) {
	_, p, err := conn.ReadMessage()
	s.NoError(err, string(p))
	err = json.Unmarshal(p, resp)
	s.NoError(err, string(p))
	s.NotEmpty(resp.EventType, string(p))
	s.Equal(enums.EVENT_TEXT_MESSAGE.String(), resp.EventType, string(p))
	s.NotEmpty(resp.FromUUID, string(p))
	s.NotEmpty(resp.DeviceUUID, string(p))
	s.NotEmpty(resp.Message.RoomUUID, string(p))
	s.NotEmpty(resp.Message.MessageText, string(p))
}

func ContainsRoomUUID(s []*requests.Room, str string) bool {
	for _, v := range s {
		if v.UUID == str {
			return true
		}
	}

	return false
}

func (s *IntegrationTestSuite) ReadOpenRoomResponse(conn *websocket.Conn, expectedMembers int) *requests.OpenRoomEvent {
	// TODO - ensure correct users are in the room
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	_, p, err := conn.ReadMessage()
	s.NoError(err)
	resp := &requests.OpenRoomEvent{}
	err = json.Unmarshal(p, resp)
	s.NoError(err)
	s.NotEmpty(resp.EventType)
	s.Equal(enums.EVENT_OPEN_ROOM.String(), resp.EventType)
	s.NotEmpty(resp.Room)
	s.NotEmpty(resp.Room.UUID)
	s.Equal(expectedMembers, len(resp.Room.Members))

	for _, m := range resp.Room.Members {
		s.NotEmpty(m.UUID)
		s.NotEmpty(m.UserUUID)
	}

	return resp
}

func (s *IntegrationTestSuite) SendTextMessage(ws *websocket.Conn, msgEvent *requests.TextMessageEvent) {
	err := ws.WriteJSON(msgEvent)
	s.NoError(err)
}

func (s *IntegrationTestSuite) GetMessagesByRoomUUIDByWithAPIKey(roomUUID string, offset int, apiKey string) *requests.GetMessagesByRoomUUIDResponse {
	url := fmt.Sprintf("http://%s:9090/get-messages-by-room-uuid?roomUuid=%s&offset=%d&key=%s", ServerHost, roomUUID, offset, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	s.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)
	s.NotEmpty(resp)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	s.NoError(err)
	s.NotEmpty(b)

	s.GreaterOrEqual(resp.StatusCode, 200)
	s.Less(resp.StatusCode, 300)

	result := &requests.GetMessagesByRoomUUIDResponse{}
	err = json.Unmarshal(b, result)
	s.NoError(err)
	s.NotNil(result)
	return result
}

func (s *IntegrationTestSuite) GetMessagesByRoomUUIDByMessagingJWT(roomUUID string, offset int, jwtToken string) (*requests.GetMessagesByRoomUUIDResponse, error) {
	url := fmt.Sprintf("http://%s:9090/get-messages-by-room-uuid?roomUuid=%s&offset=%d", ServerHost, roomUUID, offset)
	req, err := http.NewRequest("GET", url, nil)
	s.NoError(err)

	req.Header.Add("Authorization", jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)
	s.NotEmpty(resp)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	s.NoError(err)
	s.NotEmpty(b)

	s.GreaterOrEqual(resp.StatusCode, 200)
	s.Less(resp.StatusCode, 300)

	result := &requests.GetMessagesByRoomUUIDResponse{}
	err = json.Unmarshal(b, result)
	s.NoError(err)
	return result, err
}

func (s *IntegrationTestSuite) CreateClientConnection(msg *requests.SetClientConnectionEvent) (*requests.SetClientConnectionEvent, *websocket.Conn) {
	conn, _, err := websocket.DefaultDialer.Dial(SocketURL, nil)
	s.NoError(err)
	s.NotNil(conn)

	err = conn.WriteJSON(msg)
	s.NoError(err)

	_, p, err := conn.ReadMessage()
	s.NoError(err)

	rsp := &requests.SetClientConnectionEvent{}
	err = json.Unmarshal(p, &rsp)
	s.NoError(err)
	s.NotEmpty(rsp.DeviceUUID)
	s.NotEmpty(rsp.UserUUID)
	return rsp, conn

}

func (s *IntegrationTestSuite) OpenRoom(openRoomEvent *requests.CreateRoomRequest, apiKey string) *requests.CreateRoomResponse {
	postBody, err := json.Marshal(openRoomEvent)
	s.NoError(err)
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(fmt.Sprintf("http://%s:9090/create-room?key=%s", ServerHost, apiKey), "application/json", reqBody)
	s.NoError(err)
	s.GreaterOrEqual(resp.StatusCode, 200)
	s.Less(resp.StatusCode, 300)

	b, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)
	response := &requests.CreateRoomResponse{}
	json.Unmarshal(b, response)
	return response
}

func (s *IntegrationTestSuite) DeleteRoom(deleteRoomRequest *requests.DeleteRoomRequest, apiKey string) {
	postBody, err := json.Marshal(deleteRoomRequest)
	s.NoError(err)
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(fmt.Sprintf("http://%s:9090/delete-room?key=%s", ServerHost, apiKey), "application/json", reqBody)
	s.NoError(err)
	defer resp.Body.Close()
	s.GreaterOrEqual(resp.StatusCode, 200)
	s.Less(resp.StatusCode, 300)
}

func (s *IntegrationTestSuite) LeaveRoom(req *requests.LeaveRoomRequest, apiKey string) {
	postBody, err := json.Marshal(req)
	s.NoError(err)
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(fmt.Sprintf("http://%s:9090/leave-room?key=%s", ServerHost, apiKey), "application/json", reqBody)
	s.NoError(err)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	s.NoError(err)

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		fmt.Println(string(b))
	}

	s.GreaterOrEqual(resp.StatusCode, 200)
	s.LessOrEqual(resp.StatusCode, 299)
}

func (s *IntegrationTestSuite) ReadEvent(conn *websocket.Conn, v interface{}) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	_, p, err := conn.ReadMessage()
	s.NoError(err)

	err = json.Unmarshal(p, v)
	s.NoError(err)
}

func (s *IntegrationTestSuite) MakeGetRoomsByUserUUIDRequest(userUUID string, offset int, apiKey string) *requests.GetRoomsByUserUUIDResponse {
	url := fmt.Sprintf("http://%s:9090/get-rooms-by-user-uuid?userUuid=%s&offset=%d&key=%s", ServerHost, userUUID, offset, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	s.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	s.NoError(err)

	s.GreaterOrEqual(resp.StatusCode, 200)
	s.Less(resp.StatusCode, 300)

	result := &requests.GetRoomsByUserUUIDResponse{}
	err = json.Unmarshal(b, result)
	s.NoError(err)
	// TODO - test ordering
	return result
}

func (s *IntegrationTestSuite) MakeGetMessagesByRoomUUIDRequest(roomUUID string, apiKey string, offset int) *requests.GetMessagesByRoomUUIDResponse {
	url := fmt.Sprintf("http://%s:9090/get-messages-by-room-uuid?roomUuid=%s&offset=%d&key=%s", ServerHost, roomUUID, offset, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	s.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	s.NoError(err)

	s.GreaterOrEqual(resp.StatusCode, 200)
	s.Less(resp.StatusCode, 300)

	result := &requests.GetMessagesByRoomUUIDResponse{}
	err = json.Unmarshal(b, result)
	s.NoError(err)
	if len(result.Messages) > 1 {
		s.Greater(result.Messages[0].CreatedAtNano, result.Messages[len(result.Messages)-1].CreatedAtNano)
	}
	return result
}

func (s *IntegrationTestSuite) RecvSeenMessageEvent(conn *websocket.Conn, messageUUID string) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	_, p, err := conn.ReadMessage()
	s.NoError(err)
	seenMessageEvent := &requests.SeenMessageEvent{}
	err = json.Unmarshal(p, seenMessageEvent)
	s.NoError(err)

	s.NotEmpty(seenMessageEvent.EventType)
	s.Equal(enums.EVENT_SEEN_MESSAGE.String(), seenMessageEvent.EventType)
	s.NotEmpty(seenMessageEvent.MessageUUID)
	s.Equal(messageUUID, seenMessageEvent.MessageUUID)
	s.NotEmpty(seenMessageEvent.RoomUUID)
	s.NotEmpty(seenMessageEvent.UserUUID)
}

func (s *IntegrationTestSuite) SendMessagesByRoomUUIDEvent(conn *websocket.Conn, event *requests.MessagesByRoomUUIDEvent) {
	conn.SetWriteDeadline(time.Now().Add(time.Second * 2))
	err := conn.WriteJSON(event)
	s.NoError(err)
}

func (s *IntegrationTestSuite) GetValidToken(userUUID string) string {
	token, err := utils.GenerateMessagingToken(userUUID, time.Now().Add(time.Minute))
	s.NoError(err)

	return token
}

// TODO clean redis before each test
func (s *IntegrationTestSuite) GetValidAPIKey() string {
	key := uuid.New().String()
	err := s.redis.Set(context.Background(), key, requests.APIKey{
		Key: key,
	})
	s.NoError(err)
	return key
}

func (s *IntegrationTestSuite) MakeGetUserConnectionRequest(userUUID string) *requests.GetUserConnectionResponse {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:9090/get-user-connection/%s", ServerHost, userUUID), nil)
	s.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)
	s.GreaterOrEqual(resp.StatusCode, 200)
	s.Less(resp.StatusCode, 300)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	s.NoError(err)

	response := &requests.GetUserConnectionResponse{}
	err = json.Unmarshal(b, response)
	s.NoError(err)

	return response
}

func (s *IntegrationTestSuite) MakeGetChannelConnectionRequest(channelUUID string) *requests.GetChannelResponse {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:9090/get-channel/%s", ServerHost, channelUUID), nil)
	s.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)
	s.GreaterOrEqual(resp.StatusCode, 200)
	s.Less(resp.StatusCode, 300)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	s.NoError(err)

	response := &requests.GetChannelResponse{}
	err = json.Unmarshal(b, response)
	s.NoError(err)

	return response
}

func (s *IntegrationTestSuite) RecvDeletedRoomMsg(conn *websocket.Conn, resp *requests.DeleteRoomEvent) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	_, p, err := conn.ReadMessage()
	s.NoError(err)
	err = json.Unmarshal(p, resp)
	s.NoError(err)
	s.NotEmpty(resp.EventType)
	s.Equal(enums.EVENT_DELETE_ROOM.String(), resp.EventType)
	s.NotEmpty(resp.RoomUUID)
}
