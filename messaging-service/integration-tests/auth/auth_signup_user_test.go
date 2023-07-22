package integrationtests

import (
	"log"
	"messaging-service/integration-tests/common"
	"messaging-service/types/requests"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignupUser(t *testing.T) {
	// t.Skip()
	t.Run("test signup user and create auth profile", func(t *testing.T) {
		log.Printf("Running %s", t.Name())

		signupResponse := common.CreateRandomUser(t)
		common.MakeTestAuthRequest(t, signupResponse.AccessToken)
		// give some time for the time to change and so the token will also change
		time.Sleep(1 * time.Second)

		// test refresh token
		refreshTokenResp := common.MakeRefreshTokenRequest(t, signupResponse.RefreshToken)

		// test new access token
		common.MakeTestAuthRequest(t, refreshTokenResp.AccessToken)

		// create fake token with correct data
		jwtAuthProfile := &requests.AuthProfile{
			UUID:  signupResponse.UUID,
			Email: signupResponse.Email,
		}

		token, err := common.GenerateJWTAccessToken(*jwtAuthProfile, "SECRET!!")
		assert.NoError(t, err)

		common.MakeTestAuthRequestFailAuth(t, token)
	})
}

// test creating a JWT with the same info
