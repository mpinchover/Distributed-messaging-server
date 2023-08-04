package integrationtests

import (
	"fmt"
	"messaging-service/integration-tests/common"
	"messaging-service/src/types/requests"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestResetPassword(t *testing.T) {
	// t.Skip()
	t.Run("test signup user and create auth profile", func(t *testing.T) {
		t.Logf("Runningg test %s at %d", t.Name(), time.Now().UnixNano())

		password := uuid.New().String()
		confirmPassword := password
		email := fmt.Sprintf("%s@gmail.com", uuid.New().String())

		// create user
		signupRequest := &requests.SignupRequest{
			Email:           email,
			Password:        password,
			ConfirmPassword: confirmPassword,
		}
		//
		signupResp := common.MakeSignupRequest(t, signupRequest)
		assert.NotEmpty(t, signupResp.AccessToken)
		assert.NotEmpty(t, signupResp.RefreshToken)

		// test auth token
		authProfile := common.MakeTestAuthRequest(t, signupResp.AccessToken)
		assert.NotNil(t, authProfile)
		assert.Equal(t, signupRequest.Email, authProfile.Email)
		assert.NotEmpty(t, authProfile.UUID)

		// login user
		loginRequest := &requests.LoginRequest{
			Email:    email,
			Password: password,
		}
		loginResp := common.MakeLoginRequest(t, loginRequest)
		assert.NotEmpty(t, loginResp.AccessToken)
		assert.NotEmpty(t, loginResp.RefreshToken)

		updatePasswordRequest := &requests.UpdatePasswordRequest{
			CurrentPassword:    password,
			NewPassword:        "something-else",
			ConfirmNewPassword: "something-else",
		}
		// reset password
		common.MakeUpdatePasswordRequest(t, updatePasswordRequest, loginResp.AccessToken)

		// // should success
		loginRequest.Password = "something-else"
		common.MakeLoginRequest(t, loginRequest)

		// should fail
		loginRequest.Password = "something-else11"
		common.MakeLoginRequestFailAuth(t, loginRequest)
	})
}
