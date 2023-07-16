package integrationtests

import (
	"fmt"
	"log"
	"messaging-service/types/requests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestResetPassword(t *testing.T) {
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

		// login user
		loginRequest := &requests.LoginRequest{
			Email:    email,
			Password: password,
		}
		loginResp, err := makeLoginRequest(loginRequest)
		assert.NoError(t, err)
		assert.NotEmpty(t, loginResp.AccessToken)
		assert.NotEmpty(t, loginResp.RefreshToken)

		updatePasswordRequest := &requests.UpdatePasswordRequest{
			CurrentPassword:    password,
			NewPassword:        "something-else",
			ConfirmNewPassword: "something-else",
		}
		// reset password
		_, err = makeUpdatePasswordRequest(updatePasswordRequest, loginResp.AccessToken)
		assert.NoError(t, err)

		// // should success
		// loginRequest.Password = "something-else"
		// _, err = makeLoginRequest(loginRequest)
		// assert.NoError(t, err)

		// // should fail
		// loginRequest.Password = "something-else11"
		// loginResp, err = makeLoginRequest(loginRequest)
		// assert.Error(t, err)
		// assert.Nil(t, loginResp)
		// assert.Contains(t, err.Error(), "400")
	})
}
