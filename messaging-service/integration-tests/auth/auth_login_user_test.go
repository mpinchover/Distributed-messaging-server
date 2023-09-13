package integrationtests

// func TestSignupUserAndCreateAuthprofile(t *testing.T) {
// 	// t.Skip()
// 	t.Run("test signup user and create auth profile", func(t *testing.T) {
// 		t.Logf("Runningg test %s at %d", t.Name(), time.Now().UnixNano())

// 		password := uuid.New().String()
// 		confirmPassword := password
// 		email := fmt.Sprintf("%s@gmail.com", uuid.New().String())

// 		signupResponse := common.MakeSignupRequest(t, &records.SignupRequest{
// 			Email:           email,
// 			Password:        password,
// 			ConfirmPassword: confirmPassword,
// 		})
// 		common.MakeTestAuthRequest(t, signupResponse.AccessToken)
// 		// give some time for the time to change and so the token will also change
// 		time.Sleep(1 * time.Second)

// 		// test refresh token
// 		refreshTokenResp := common.MakeRefreshTokenRequest(t, signupResponse.RefreshToken)

// 		// test new access token
// 		common.MakeTestAuthRequest(t, refreshTokenResp.AccessToken)

// 		// create fake token with correct data
// 		jwtAuthProfile := &records.AuthProfile{
// 			UUID:  signupResponse.UUID,
// 			Email: signupResponse.Email,
// 		}

// 		token, err := common.GenerateJWTAccessToken(*jwtAuthProfile, "SECRET!!")
// 		assert.NoError(t, err)

// 		common.MakeTestAuthRequestFailAuth(t, token)

// 		// login user
// 		loginRequest := &records.LoginRequest{
// 			Email:    email,
// 			Password: password,
// 		}

// 		loginResp := common.MakeLoginRequest(t, loginRequest)

// 		// test auth token
// 		common.MakeTestAuthRequest(t, loginResp.AccessToken)
// 		time.Sleep(1 * time.Second)

// 		// test refresh token
// 		common.MakeRefreshTokenRequest(t, loginResp.RefreshToken)

// 		// test new access token
// 		common.MakeTestAuthRequest(t, refreshTokenResp.AccessToken)

// 		// should fail to login
// 		loginRequest.Password = "something-else"
// 		common.MakeLoginRequestFailAuth(t, loginRequest)

// 		token, err = common.GenerateJWTAccessToken(*jwtAuthProfile, "SECRET!!")
// 		assert.NoError(t, err)

// 		common.MakeTestAuthRequestFailAuth(t, token)

// 		// should work
// 		_, err = common.GenerateJWTAccessToken(*jwtAuthProfile, "SECRET")
// 		assert.NoError(t, err)
// 	})
// }
