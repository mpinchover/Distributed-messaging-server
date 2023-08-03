package authcontroller

import (
	"context"
	"errors"
	"fmt"
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"
	"messaging-service/src/utils"
	"testing"
	"time"

	mockRedis "messaging-service/mocks/src/redis"
	mockRepo "messaging-service/mocks/src/repo"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthControllerSuite struct {
	suite.Suite

	authCtrlr       *AuthController
	mockRepoClient  *mockRepo.RepoInterface
	mockRedisClient *mockRedis.RedisInterface
}

// this function executes before the test suite begins execution
func (s *AuthControllerSuite) SetupSuite() {
	fmt.Println(">>> From SetupSuite")

	s.mockRepoClient = mockRepo.NewRepoInterface(s.T())
	s.mockRedisClient = mockRedis.NewRedisInterface(s.T())
	s.authCtrlr = &AuthController{
		repo:        s.mockRepoClient,
		redisClient: s.mockRedisClient,
	}
	//
}

func TestAuthControllerSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerSuite))
}

func (s *AuthControllerSuite) TestLogin() {

	hashedPassword, err := hashPassword("password")
	s.NoError(err)
	tests := []struct {
		test        string
		request     *requests.LoginRequest
		expectedErr string
		mocks       func()
	}{
		{
			test: "database throws an error",
			request: &requests.LoginRequest{
				Email:    "email@gmail.com",
				Password: "some-password",
			},
			expectedErr: "database error",
			mocks: func() {
				s.mockRepoClient.On("GetAuthProfileByEmail", mock.Anything).Return(nil, errors.New("database error")).Once()
			},
		},
		{
			test: "database throws an error",
			request: &requests.LoginRequest{
				Email:    "email@gmail.com",
				Password: "some-password",
			},
			expectedErr: "authorization error",
			mocks: func() {
				s.mockRepoClient.On("GetAuthProfileByEmail", mock.Anything).Return(nil, nil).Once()
			},
		},
		{
			test: "email/password doesn't match",
			request: &requests.LoginRequest{
				Email:    "email@gmail.com",
				Password: "some-password",
			},
			expectedErr: "old/new passwords do not match",
			mocks: func() {
				mockAuthProfile := &records.AuthProfile{
					Email:          "email@gmail.com",
					HashedPassword: "unhashedpassword",
					UUID:           "some-uuid",
				}
				s.mockRepoClient.On("GetAuthProfileByEmail", mock.Anything).Return(mockAuthProfile, nil).Once()
			},
		},
		{
			test: "success",
			request: &requests.LoginRequest{
				Email:    "email@gmail.com",
				Password: "password",
			},
			mocks: func() {
				mockAuthProfile := &records.AuthProfile{
					Email:          "email@gmail.com",
					HashedPassword: hashedPassword,
					UUID:           "some-uuid",
				}
				s.mockRepoClient.On("GetAuthProfileByEmail", mock.Anything).Return(mockAuthProfile, nil).Once()
			},
		},
	}

	for _, t := range tests {
		t.mocks()
		res, err := s.authCtrlr.Login(t.request)
		if t.expectedErr != "" {
			s.Error(err, t.test)
			s.Nil(res, t.test)
			s.Contains(err.Error(), t.expectedErr, t.test)
		} else {
			s.NoError(err, t.test)
			s.NotNil(res, t.test)
			s.NotEmpty(res.AccessToken, t.test)
			s.NotEmpty(res.RefreshToken, t.test)
		}
	}
}

func (s *AuthControllerSuite) TestUpdatePassword() {
	// hashedPassword, err := hashPassword("password")
	// s.NoError(err)

	tests := []struct {
		test        string
		request     *requests.UpdatePasswordRequest
		expectedErr string
		mocks       func()
		ctx         func() context.Context
	}{
		{
			test: "could not extract auth profile from token",
			request: &requests.UpdatePasswordRequest{
				CurrentPassword:    "email@gmail.com",
				NewPassword:        "some-password",
				ConfirmNewPassword: "some-password",
			},
			expectedErr: "could not get auth profile",
			mocks: func() {
				// s.mockRepoClient.On("GetAuthProfileByEmail", mock.Anything).Return(nil, errors.New("database error")).Once()
			},
			ctx: func() context.Context {
				return context.Background()
			},
		},
		{
			test: "getAuthProfileByEmail throws error",
			request: &requests.UpdatePasswordRequest{
				CurrentPassword:    "password",
				NewPassword:        "some-password",
				ConfirmNewPassword: "some-password",
			},
			expectedErr: "database error",
			mocks: func() {
				s.mockRepoClient.On("GetAuthProfileByEmail", mock.Anything).Return(nil, errors.New("database error")).Once()
			},
			ctx: func() context.Context {
				authProfile := &requests.AuthProfile{
					Email: "somewhere@gmail.com",
					UUID:  "some-uuid",
				}
				token, err := utils.GenerateJWTToken(authProfile, time.Now().Add(10*time.Minute))
				s.NoError(err)
				s.NotEmpty(token)

				return context.WithValue(context.Background(), "AUTH_PROFILE", authProfile)
			},
		},
		{
			test: "getAuthProfileByEmail cannot find account with email address",
			request: &requests.UpdatePasswordRequest{
				CurrentPassword:    "password",
				NewPassword:        "some-password",
				ConfirmNewPassword: "some-password",
			},
			expectedErr: "no account matching email found",
			mocks: func() {
				s.mockRepoClient.On("GetAuthProfileByEmail", mock.Anything).Return(nil, nil).Once()
			},
			ctx: func() context.Context {
				authProfile := &requests.AuthProfile{
					Email: "somewhere@gmail.com",
					UUID:  "some-uuid",
				}
				token, err := utils.GenerateJWTToken(authProfile, time.Now().Add(10*time.Minute))
				s.NoError(err)
				s.NotEmpty(token)

				return context.WithValue(context.Background(), "AUTH_PROFILE", authProfile)
			},
		},
		{
			test: "old passwords dont match",
			request: &requests.UpdatePasswordRequest{
				CurrentPassword:    "password",
				NewPassword:        "some-password",
				ConfirmNewPassword: "some-password",
			},
			expectedErr: "old/new passwords do not match",
			mocks: func() {
				hashedPass, err := hashPassword("wrong-password")
				s.NoError(err)
				s.NotEmpty(hashedPass)

				s.mockRepoClient.On("GetAuthProfileByEmail", mock.Anything).Return(&records.AuthProfile{
					UUID:           "some-uuid",
					Email:          "somewhere@gmail.com",
					HashedPassword: hashedPass,
				}, nil).Once()
			},
			ctx: func() context.Context {
				authProfile := &requests.AuthProfile{
					Email: "somewhere@gmail.com",
					UUID:  "some-uuid",
				}
				token, err := utils.GenerateJWTToken(authProfile, time.Now().Add(10*time.Minute))
				s.NoError(err)
				s.NotEmpty(token)

				return context.WithValue(context.Background(), "AUTH_PROFILE", authProfile)
			},
		},
		{
			test: "update password fails",
			request: &requests.UpdatePasswordRequest{
				CurrentPassword:    "password",
				NewPassword:        "some-password",
				ConfirmNewPassword: "some-password",
			},
			expectedErr: "error on update password",
			mocks: func() {
				hashedPass, err := hashPassword("password")
				s.NoError(err)
				s.NotEmpty(hashedPass)

				s.mockRepoClient.On("GetAuthProfileByEmail", mock.Anything).Return(&records.AuthProfile{
					UUID:           "some-uuid",
					Email:          "somewhere@gmail.com",
					HashedPassword: hashedPass,
				}, nil).Once()

				s.mockRepoClient.On("UpdatePassword", mock.Anything, mock.Anything).Return(errors.New("error on update password")).Once()
			},
			ctx: func() context.Context {
				authProfile := &requests.AuthProfile{
					Email: "somewhere@gmail.com",
					UUID:  "some-uuid",
				}
				token, err := utils.GenerateJWTToken(authProfile, time.Now().Add(10*time.Minute))
				s.NoError(err)
				s.NotEmpty(token)

				return context.WithValue(context.Background(), "AUTH_PROFILE", authProfile)
			},
		},
		{
			test: "update password succeeds",
			request: &requests.UpdatePasswordRequest{
				CurrentPassword:    "password",
				NewPassword:        "some-password",
				ConfirmNewPassword: "some-password",
			},
			mocks: func() {
				hashedPass, err := hashPassword("password")
				s.NoError(err)
				s.NotEmpty(hashedPass)

				s.mockRepoClient.On("GetAuthProfileByEmail", mock.Anything).Return(&records.AuthProfile{
					UUID:           "some-uuid",
					Email:          "somewhere@gmail.com",
					HashedPassword: hashedPass,
				}, nil).Once()

				s.mockRepoClient.On("UpdatePassword", mock.Anything, mock.Anything).Return(nil).Once()
			},
			ctx: func() context.Context {
				authProfile := &requests.AuthProfile{
					Email: "somewhere@gmail.com",
					UUID:  "some-uuid",
				}
				token, err := utils.GenerateJWTToken(authProfile, time.Now().Add(10*time.Minute))
				s.NoError(err)
				s.NotEmpty(token)

				return context.WithValue(context.Background(), "AUTH_PROFILE", authProfile)
			},
		},
	}

	for _, t := range tests {
		t.mocks()
		err := s.authCtrlr.UpdatePassword(t.ctx(), t.request)
		fmt.Println("ERROR IS")
		fmt.Println(err)
		if t.expectedErr != "" {
			s.Error(err, t.test)
			s.Contains(err.Error(), t.expectedErr, t.test)
		} else {
			s.NoError(err, t.test)
		}
	}
}

// func (s *AuthControllerSuite) TestResetPassword() {
// 	tests := []struct {
// 		test  string
// 		req   *requests.ResetPasswordRequest
// 		expecteErr string
// 		mocks func()
// 	}{
// 		{
// 			test: "GetEmailByPasswordResetToken failure",
// 			expecteErr:"int error"
// 			req: &requests.ResetPasswordRequest{
// 				Token: "some-token",
// 			},
// 			mocks: func() {
// 				s.mockRedisClient.On("GetEmailByPasswordResetToken", mock.Anything, mock.Anything).Return("", errors.New("int error")).Once()
// 			},
// 		},
// 	}

// 	for _, t := range tests {
// 		t.mocks()
// 		err := s.authCtrlr.ResetPassword(context.Background(), t.req)
// 	}
// }

func (s *AuthControllerSuite) TestVerifyAPIKeyExists() {
	tests := []struct {
		test        string
		expectedErr string
		mocks       func()
	}{
		{
			test:        "redis error",
			expectedErr: "internal error",
			mocks: func() {
				s.mockRedisClient.On("GetAPIKey", mock.Anything, mock.Anything).Return(nil, errors.New("internal error")).Once()
			},
		},
		{
			test:        "no api key found",
			expectedErr: "could not find API key",
			mocks: func() {
				s.mockRedisClient.On("GetAPIKey", mock.Anything, mock.Anything).Return(nil, nil).Once()
			},
		},
		{
			test: "api key found, success",
			mocks: func() {
				s.mockRedisClient.On("GetAPIKey", mock.Anything, mock.Anything).Return(&requests.APIKey{Key: "some-key"}, nil).Once()
			},
		},
	}

	for _, t := range tests {
		t.mocks()
		apiKey, err := s.authCtrlr.VerifyAPIKeyExists(context.Background(), "key")

		if t.expectedErr != "" {
			s.Error(err, t.test)
			s.Nil(apiKey, t.test)
			s.Contains(err.Error(), t.expectedErr, t.test)
		} else {
			s.NoError(err, t.test)
			s.NotNil(apiKey, t.test)
			s.NotEmpty(apiKey.Key, t.test)
			// TODO - match regex of API key
		}
	}
}

func (s *AuthControllerSuite) TestRemoveAPIKey() {
	tests := []struct {
		test        string
		expectedErr string
		mocks       func()
	}{
		{
			test:        "remove key failure",
			expectedErr: "internal error",
			mocks: func() {
				s.mockRedisClient.On("Del", mock.Anything, mock.Anything).Return(errors.New("internal error")).Once()
			},
		},
	}

	for _, t := range tests {
		t.mocks()
		err := s.authCtrlr.RemoveAPIKey(context.Background(), "key")
		if t.expectedErr != "" {
			s.Error(err, t.test)
			s.Contains(err.Error(), t.expectedErr)
		} else {
			s.NoError(err, t.test)
		}
	}

}
