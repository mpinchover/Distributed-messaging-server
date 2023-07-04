package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"messaging-service/types/requests"
	"net/http"
)

// func TestSignupUserAndCreateAuthprofile(t *testing.T) {
// 	t.Run("test signup user and create auth profile", func(t *testing.T) {
// 		log.Printf("Running %s", t.Name())

// 		password := "password"
// 		confirmPassword := "password"
// 		email := "email@gmail.com"

// 		authProfile := &requests.SignupRequest{
// 			Email:           email,
// 			Password:        password,
// 			ConfirmPassword: confirmPassword,
// 		}

// 		signupResp, err := makeSignupRequest(authProfile)
// 		assert.NoError(t, err)
// 		// test private route
// 	})
// }

func makeSignupRequest(authProfile *requests.SignupRequest) (*requests.SignupResponse, error) {
	postBody, err := json.Marshal(authProfile)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:9090/signup", "application/json", reqBody)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return nil, fmt.Errorf("status code is %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &requests.SignupResponse{}
	err = json.Unmarshal(b, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
