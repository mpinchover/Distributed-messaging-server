package common

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

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

const (
	SocketURL = "ws://localhost:9090/ws"
)

func SendMessages(t *testing.T, fromUserUUID string, connectionUUID string, roomUUID string, conn *websocket.Conn, token string) {
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
			Token: token,
		}
		SendTextMessage(t, conn, msgEventOut)
	}
}

func RecvMessages(t *testing.T, conn *websocket.Conn) {
	for i := 0; i < 25; i++ {
		// conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		resp := &requests.TextMessageEvent{}
		RecvMessage(t, conn, resp)
	}
}

func RecvMessage(t *testing.T, conn *websocket.Conn, resp *requests.TextMessageEvent) {
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err, string(p))
	err = json.Unmarshal(p, resp)
	assert.NoError(t, err, string(p))
	assert.NotEmpty(t, resp.EventType, string(p))
	assert.Equal(t, enums.EVENT_TEXT_MESSAGE.String(), resp.EventType, string(p))
	assert.NotEmpty(t, resp.FromUUID, string(p))
	assert.NotEmpty(t, resp.ConnectionUUID, string(p))
	assert.NotEmpty(t, resp.Message.RoomUUID, string(p))
	assert.NotEmpty(t, resp.Message.MessageText, string(p))
}

func ContainsRoomUUID(s []*requests.Room, str string) bool {
	for _, v := range s {
		if v.UUID == str {
			return true
		}
	}

	return false
}

func ReadOpenRoomResponse(t *testing.T, conn *websocket.Conn, expectedMembers int) *requests.OpenRoomEvent {
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

func SendTextMessage(t *testing.T, ws *websocket.Conn, msgEvent *requests.TextMessageEvent) {
	err := ws.WriteJSON(msgEvent)
	assert.NoError(t, err)
}

func GetMessagesByRoomUUIDByWithAPIKey(t *testing.T, roomUUID string, offset int, apiKey string) *requests.GetMessagesByRoomUUIDResponse {
	url := fmt.Sprintf("http://localhost:9090/get-messages-by-room-uuid?roomUuid=%s&offset=%d&key=%s", roomUUID, offset, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)

	result := &requests.GetMessagesByRoomUUIDResponse{}
	err = json.Unmarshal(b, result)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	return result
}

func GetMessagesByRoomUUIDByMessagingJWT(t *testing.T, roomUUID string, offset int, jwtToken string) (*requests.GetMessagesByRoomUUIDResponse, error) {
	url := fmt.Sprintf("http://localhost:9090/get-messages-by-room-uuid?roomUuid=%s&offset=%d", roomUUID, offset)
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	req.Header.Add("Authorization", jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)

	result := &requests.GetMessagesByRoomUUIDResponse{}
	err = json.Unmarshal(b, result)
	assert.NoError(t, err)
	return result, err
}

func GetRoomsByUserUUIDByMessagingJWT(t *testing.T, userUUID string, offset int, jwtToken string) *requests.GetRoomsByUserUUIDResponse {
	url := fmt.Sprintf("http://localhost:9090/get-rooms-by-user-uuid?userUuid=%s&offset=%d", userUUID, offset)
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	req.Header.Add("Authorization", jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)

	result := &requests.GetRoomsByUserUUIDResponse{}
	err = json.Unmarshal(b, result)
	assert.NoError(t, err)
	return result
}

func GetRoomsByUserUUIDWithApiKey(t *testing.T, userUUID string, offset int, apiKey string) *requests.GetRoomsByUserUUIDResponse {
	url := fmt.Sprintf("http://localhost:9090/get-rooms-by-user-uuid?userUuid=%s&offset=%d&key=%s", userUUID, offset, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)

	result := &requests.GetRoomsByUserUUIDResponse{}
	err = json.Unmarshal(b, result)
	assert.NoError(t, err)
	return result
}

func CreateClientConnection(t *testing.T, msg *requests.SetClientConnectionEvent) (*requests.SetClientConnectionEvent, *websocket.Conn) {
	conn, _, err := websocket.DefaultDialer.Dial(SocketURL, nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	err = conn.WriteJSON(msg)
	assert.NoError(t, err)

	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)

	rsp := &requests.SetClientConnectionEvent{}
	err = json.Unmarshal(p, &rsp)
	assert.NoError(t, err)
	assert.NotEmpty(t, rsp.ConnectionUUID)
	assert.NotEmpty(t, rsp.UserUUID)
	return rsp, conn

}

func OpenRoom(t *testing.T, openRoomEvent *requests.CreateRoomRequest, apiKey string) {
	postBody, err := json.Marshal(openRoomEvent)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(fmt.Sprintf("http://localhost:9090/create-room?key=%s", apiKey), "application/json", reqBody)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)
}

func DeleteRoom(t *testing.T, deleteRoomRequest *requests.DeleteRoomRequest, apiKey string) {
	postBody, err := json.Marshal(deleteRoomRequest)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(fmt.Sprintf("http://localhost:9090/delete-room?key=%s", apiKey), "application/json", reqBody)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)
}

func LeaveRoom(t *testing.T, req *requests.LeaveRoomRequest, apiKey string) {
	postBody, err := json.Marshal(req)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(fmt.Sprintf("http://localhost:9090/leave-room?key=%s", apiKey), "application/json", reqBody)
	// resp, err := http.Post("http://localhost:9090/leave-room", "application/json", reqBody)
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

func ReadEvent(t *testing.T, conn *websocket.Conn, v interface{}) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	_, p, err := conn.ReadMessage()
	assert.NoError(t, err)

	err = json.Unmarshal(p, v)
	assert.NoError(t, err)
}

// func QueryMessagesWithApiKey(t *testing.T, userUUID string, roomUUID string, expectedRooms int, apiKey string) {
// 	resp, err := GetMessagesByRoomUUIDByWithAPIKey(t, roomUUID, 0, jwtToken)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 20, len(resp.Messages))

// 	totalMessages = append(totalMessages, resp.Messages...)
// 	assert.Equal(t, 20, len(totalMessages))

// 	resp, err = GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, len(totalMessages), jwtToken)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 20, len(resp.Messages))
// 	totalMessages = append(totalMessages, resp.Messages...)
// 	assert.Equal(t, 40, len(totalMessages))

// 	resp, err = GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, len(totalMessages), jwtToken)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 10, len(resp.Messages))
// 	totalMessages = append(totalMessages, resp.Messages...)
// 	assert.Equal(t, 50, len(totalMessages))

// 	// jump by 15 because the msgs are being sent too fast.
// 	for i := 15; i < len(totalMessages); i++ {
// 		// prev := totalMessages[i-1]
// 		// cur := totalMessages[i]
// 		// assert.True(t, prev.CreatedAt >= cur.CreatedAt)
// 	}
// }

func QueryMessagesByMessagingJWT(t *testing.T, userUUID string, roomUUID string, expectedRooms int, jwtToken string) {
	totalMessages := []*requests.Message{}
	res := GetRoomsByUserUUIDByMessagingJWT(t, userUUID, 0, jwtToken)

	assert.NotEmpty(t, res)
	assert.Equal(t, expectedRooms, len(res.Rooms))

	// ensure it contains the room uuid
	assert.True(t, ContainsRoomUUID(res.Rooms, roomUUID))

	resp, err := GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, 0, jwtToken)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(resp.Messages))

	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, 20, len(totalMessages))

	resp, err = GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, len(totalMessages), jwtToken)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(resp.Messages))
	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, 40, len(totalMessages))

	resp, err = GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, len(totalMessages), jwtToken)
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

func QueryMessagesByAPIKey(t *testing.T, userUUID string, roomUUID string, expectedRooms int, apiKey string) {
	totalMessages := []*requests.Message{}
	res := GetRoomsByUserUUIDByMessagingJWT(t, userUUID, 0, apiKey)

	assert.NotEmpty(t, res)
	assert.Equal(t, expectedRooms, len(res.Rooms))

	// ensure it contains the room uuid
	assert.True(t, ContainsRoomUUID(res.Rooms, roomUUID))

	resp, err := GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, 0, apiKey)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(resp.Messages))

	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, 20, len(totalMessages))

	resp, err = GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, len(totalMessages), apiKey)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(resp.Messages))
	totalMessages = append(totalMessages, resp.Messages...)
	assert.Equal(t, 40, len(totalMessages))

	resp, err = GetMessagesByRoomUUIDByMessagingJWT(t, roomUUID, len(totalMessages), apiKey)
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

func RecvSeenMessageEvent(t *testing.T, conn *websocket.Conn, messageUUID string) {
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

func MakeSignupRequest(t *testing.T, authProfile *requests.SignupRequest) *requests.SignupResponse {
	postBody, err := json.Marshal(authProfile)
	assert.NoError(t, err)
	assert.NotNil(t, postBody)

	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:9090/signup", "application/json", reqBody)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	response := &requests.SignupResponse{}
	err = json.Unmarshal(b, response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.NotEmpty(t, response.UUID)
	return response
}

func MakeLoginRequest(t *testing.T, loginReq *requests.LoginRequest) *requests.LoginResponse {
	postBody, err := json.Marshal(loginReq)
	assert.NoError(t, err)

	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:9090/login", "application/json", reqBody)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	response := &requests.LoginResponse{}
	err = json.Unmarshal(b, response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)

	return response
}

func MakeLoginRequestFailAuth(t *testing.T, loginReq *requests.LoginRequest) {
	postBody, err := json.Marshal(loginReq)
	assert.NoError(t, err)

	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:9090/login", "application/json", reqBody)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func MakeUpdatePasswordRequest(t *testing.T, request *requests.UpdatePasswordRequest, token string) (*requests.GenericResponse, error) {

	postBody, err := json.Marshal(request)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", "http://localhost:9090/update-password", reqBody)
	assert.NoError(t, err)
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	response := &requests.GenericResponse{}
	err = json.Unmarshal(b, response)
	assert.NoError(t, err)

	return response, nil
}
func MakeGenerateMessagingTokenRequest(t *testing.T, request *requests.GenerateMessagingTokenRequest, apiKey string) *requests.GenerateMessagingTokenResponse {
	postBody, err := json.Marshal(request)
	assert.NoError(t, err)
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(fmt.Sprintf("http://localhost:9090/generate-messaging-token?key=%s", apiKey), "application/json", reqBody)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)

	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	response := &requests.GenerateMessagingTokenResponse{}
	err = json.Unmarshal(b, response)
	assert.NoError(t, err)
	return response
}

func GetValidToken(t *testing.T) (string, string) {
	newUserResponse := CreateRandomUser(t)
	// get a messaging token for this user
	apiKey := MakeGetAPIKeyRequest(t, newUserResponse.AccessToken)

	generateMessagingToken := &requests.GenerateMessagingTokenRequest{
		UserID: newUserResponse.UUID,
	}
	generateMessagingTokenResp := MakeGenerateMessagingTokenRequest(t, generateMessagingToken, apiKey.Key)
	assert.NotEmpty(t, generateMessagingTokenResp.Token)

	return generateMessagingTokenResp.Token, apiKey.Key
}

func CreateRandomUser(t *testing.T) *requests.SignupResponse {
	password := uuid.New().String()
	confirmPassword := password
	email := fmt.Sprintf("%s@gmail.com", uuid.New().String())

	signupResponse := MakeSignupRequest(t, &requests.SignupRequest{
		Email:           email,
		Password:        password,
		ConfirmPassword: confirmPassword,
	})
	assert.NotNil(t, signupResponse)
	assert.NotEmpty(t, signupResponse.AccessToken)
	assert.NotEmpty(t, signupResponse.RefreshToken)
	assert.NotEmpty(t, signupResponse.UUID)
	return signupResponse
}

func MakeTestAuthRequest(t *testing.T, token string) *requests.AuthProfile {
	req, err := http.NewRequest("GET", "http://localhost:9090/test-auth-profile", nil)
	assert.NoError(t, err)

	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	response := &requests.AuthProfile{}
	err = json.Unmarshal(b, response)
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.Email, response.Email)
	assert.NotEmpty(t, response.UUID)
	return response
}

func MakeTestAuthRequestFailAuth(t *testing.T, token string) {
	req, err := http.NewRequest("GET", "http://localhost:9090/test-auth-profile", nil)
	assert.NoError(t, err)

	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func MakeRefreshTokenRequest(t *testing.T, refreshToken string) *requests.RefreshAccessTokenResponse {
	req, err := http.NewRequest("GET", "http://localhost:9090/refresh-token", nil)
	assert.NoError(t, err)

	req.Header.Add("Authorization", refreshToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	response := &requests.RefreshAccessTokenResponse{}
	err = json.Unmarshal(b, response)
	assert.NoError(t, err)

	return response
}

func GenerateJWTAccessToken(authProfile requests.AuthProfile, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["AUTH_PROFILE"] = authProfile
	claims["EXP"] = time.Now().UTC().Add(20 * time.Minute).Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func MakeGetAPIKeyRequest(t *testing.T, token string) *requests.APIKey {
	req, err := http.NewRequest("GET", "http://localhost:9090/get-new-api-key", nil)
	assert.NoError(t, err)

	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.StatusCode, 200)
	assert.Less(t, resp.StatusCode, 300)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	response := &requests.APIKey{}
	err = json.Unmarshal(b, response)
	assert.NoError(t, err)

	return response
}