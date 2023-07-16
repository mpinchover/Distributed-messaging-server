package integrationtests

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"messaging-service/types/requests"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSignupUserAndCreateAuthprofile(t *testing.T) {
	t.Run("test signup user and create auth profile", func(t *testing.T) {
		log.Printf("Running %s", t.Name())

		password := uuid.New().String()
		confirmPassword := password
		email := fmt.Sprintf("%s@gmail.com", uuid.New().String())

		// create user
		signupRequest := &requests.SignupRequest{
			Email:           email,
			Password:        password,
			ConfirmPassword: confirmPassword,
		}

		signupResp, err := makeSignupRequest(signupRequest)
		assert.NoError(t, err)
		assert.NotEmpty(t, signupResp.AccessToken)
		assert.NotEmpty(t, signupResp.RefreshToken)

		// test auth token
		authProfile, err := makeTestAuthRequest(signupResp.AccessToken)
		assert.NoError(t, err)
		assert.NotNil(t, authProfile)
		assert.Equal(t, signupRequest.Email, authProfile.Email)
		assert.NotEmpty(t, authProfile.UUID)

		// give some time for the time to change and so the token will also change
		time.Sleep(1 * time.Second)
		// test refresh token
		refreshTokenResp, err := makeRefreshTokenRequest(signupResp.RefreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, refreshTokenResp)
		assert.NotEqual(t, refreshTokenResp.RefreshToken, signupResp.RefreshToken)
		assert.NotEqual(t, refreshTokenResp.AccessToken, signupResp.AccessToken)

		// test new access token
		authProfile, err = makeTestAuthRequest(refreshTokenResp.AccessToken)
		assert.NoError(t, err)
		assert.NotNil(t, authProfile)
		assert.Equal(t, signupRequest.Email, authProfile.Email)
		assert.NotEmpty(t, authProfile.UUID)

		// login user
		loginRequest := &requests.LoginRequest{
			Email:    email,
			Password: password,
		}
		loginResp, err := makeLoginRequest(loginRequest)
		assert.NoError(t, err)
		assert.NotEmpty(t, loginResp.AccessToken)
		assert.NotEmpty(t, loginResp.RefreshToken)

		// test auth token
		authProfile, err = makeTestAuthRequest(loginResp.AccessToken)
		assert.NoError(t, err)
		assert.NotNil(t, authProfile)
		assert.Equal(t, loginRequest.Email, authProfile.Email)
		assert.NotEmpty(t, authProfile.UUID)

		time.Sleep(1 * time.Second)

		// test refresh token
		refreshTokenResp, err = makeRefreshTokenRequest(loginResp.RefreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, refreshTokenResp)
		assert.NotEqual(t, refreshTokenResp.RefreshToken, loginResp.RefreshToken)
		assert.NotEqual(t, refreshTokenResp.AccessToken, loginResp.AccessToken)

		// test new access token
		authProfile, err = makeTestAuthRequest(refreshTokenResp.AccessToken)
		assert.NoError(t, err)
		assert.NotNil(t, authProfile)
		assert.Equal(t, loginRequest.Email, authProfile.Email)
		assert.NotEmpty(t, authProfile.UUID)

		// should fail to login
		loginRequest.Password = "something-else"
		loginResp, err = makeLoginRequest(loginRequest)
		assert.Error(t, err)
		assert.Nil(t, loginResp)
		assert.Contains(t, err.Error(), "400")

		// create fake token with correct data
		jwtAuthProfile := &requests.AuthProfile{
			UUID:  authProfile.UUID,
			Email: authProfile.Email,
		}

		token, err := generateJWTAccessToken(*jwtAuthProfile, "SECRET!!")
		assert.NoError(t, err)

		_, err = makeTestAuthRequest(token)
		assert.Error(t, err)

		// should work
		token, err = generateJWTAccessToken(*jwtAuthProfile, "SECRET")
		assert.NoError(t, err)

		authProfile, err = makeTestAuthRequest(token)
		assert.NoError(t, err)
		assert.NotNil(t, authProfile)
		assert.Equal(t, jwtAuthProfile.Email, authProfile.Email)
		assert.NotEmpty(t, authProfile.UUID)

	})
}

func generateJWTAccessToken(authProfile requests.AuthProfile, secret string) (string, error) {
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

// test creating a JWT with the same info

func makeTestAuthRequest(token string) (*requests.AuthProfile, error) {
	req, err := http.NewRequest("GET", "http://localhost:9090/test-auth-profile", nil)
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

	response := &requests.AuthProfile{}
	err = json.Unmarshal(b, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func makeRefreshTokenRequest(refreshToken string) (*requests.RefreshAccessTokenResponse, error) {
	req, err := http.NewRequest("GET", "http://localhost:9090/refresh-token", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", refreshToken)
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

	response := &requests.RefreshAccessTokenResponse{}
	err = json.Unmarshal(b, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
