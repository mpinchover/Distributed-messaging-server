package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"messaging-service/types/requests"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSignupUserAndGetAPIKey(t *testing.T) {
	t.Run("sign up user and get API key", func(t *testing.T) {
		log.Printf("Running %s", t.Name())

		// create user
		password := uuid.New().String()
		confirmPassword := password
		email := fmt.Sprintf("%s@gmail.com", uuid.New().String())

		signupRequest := &requests.SignupRequest{
			Email:           email,
			Password:        password,
			ConfirmPassword: confirmPassword,
		}

		signupResp, err := makeSignupRequest(signupRequest)
		assert.NoError(t, err)
		assert.NotEmpty(t, signupResp.AccessToken)

		// get the api key
		apiKeyResp, err := makeGetAPIKeyRequest(signupResp.AccessToken)
		assert.NoError(t, err)
		assert.NotNil(t, apiKeyResp)
		assert.NotEmpty(t, apiKeyResp.Key)

		// use the api key
		testApiKeyResp, err := makeTestAuthAPIKey(apiKeyResp.Key)
		assert.NoError(t, err)
		assert.NotNil(t, testApiKeyResp)
		assert.Equal(t, apiKeyResp.Key, testApiKeyResp.Key)

		// invalidate the api key
		invalidateResp, err := makeInvalidateAPIKeyRequest(signupResp.AccessToken, testApiKeyResp.Key)
		assert.NoError(t, err)
		assert.NotNil(t, invalidateResp)
		assert.True(t, invalidateResp.Success)

		// api key should fail this time
		_, err = makeTestAuthAPIKey(apiKeyResp.Key)
		assert.Error(t, err)
		// assert.NotNil(t, testApiKeyResp)
		// assert the error here too

		// log user in
		loginRequest := &requests.LoginRequest{
			Email:    email,
			Password: password,
		}
		loginResp, err := makeLoginRequest(loginRequest)
		assert.NoError(t, err)
		assert.NotEmpty(t, loginResp.AccessToken)

		// get new api key
		apiKeyResp, err = makeGetAPIKeyRequest(loginResp.AccessToken)
		assert.NoError(t, err)
		assert.NotNil(t, apiKeyResp)
		assert.NotEmpty(t, apiKeyResp.Key)

		// use api key again
		testApiKeyResp, err = makeTestAuthAPIKey(apiKeyResp.Key)
		assert.NoError(t, err)
		assert.NotNil(t, testApiKeyResp)
		assert.Equal(t, apiKeyResp.Key, testApiKeyResp.Key)

		// invalidate the api key
		invalidateResp, err = makeInvalidateAPIKeyRequest(loginResp.AccessToken, testApiKeyResp.Key)
		assert.NoError(t, err)
		assert.NotNil(t, invalidateResp)
		assert.True(t, invalidateResp.Success)

		// api key should fail this time
		_, err = makeTestAuthAPIKey(apiKeyResp.Key)
		assert.Error(t, err)
		// assert.NotNil(t, testApiKeyResp)
	})
}

func makeTestAuthAPIKey(apiKey string) (*requests.APIKey, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:9090/test-auth-api-key?key=%s", apiKey), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
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

	response := &requests.APIKey{}
	err = json.Unmarshal(b, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func makeGetAPIKeyRequest(token string) (*requests.APIKey, error) {
	req, err := http.NewRequest("GET", "http://localhost:9090/get-new-api-key", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
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

	response := &requests.APIKey{}
	err = json.Unmarshal(b, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func makeInvalidateAPIKeyRequest(token string, apiKey string) (*requests.GenericResponse, error) {
	body := requests.InvalidateAPIKeyRequest{
		Key: apiKey,
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", "http://localhost:9090/invalidate-api-key", reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
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

	response := &requests.GenericResponse{}
	err = json.Unmarshal(b, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
