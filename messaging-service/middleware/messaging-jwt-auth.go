package middleware

import (
	"context"
	"fmt"
	"messaging-service/controllers/authcontroller"
	"messaging-service/serrors"
	"messaging-service/utils"
	"net/http"
)

type MessagingJWT struct {
	authController *authcontroller.AuthController
}

func NewMessagingJWT(authController *authcontroller.AuthController) *MessagingJWT {
	return &MessagingJWT{
		authController: authController,
	}
}

/*
Web Socket Connection

When a token expires during an active web socket connection, the web socket will continue to
stay connected allowing the user to stay in the App. Once the user has left the App and the web
socket connection disconnects, the client will then need a new token to access the Stream API once again.
 The Web Socket connection could disconnect from either the user quitting or backgrounding the App.
https://support.getstream.io/hc/en-us/articles/360060576774-Token-Creation-Best-Practices-Chat#:~:text=Static%20tokens%20do%20not%20have%20an%20expiration%20time.

*/
// an external service gives this service the APP_ID and the CLIENT_SECRET
// use those to auth the user and generate a JWT
// the APP_ID will track to which ext service is using this messaging service and
// the user uuid is tracking which user is using it
// if the JWT has expired, external service should call this service to generate a new token
// todo - move this to utils

func (a *MessagingJWT) execute(h HTTPHandler) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		var tokenString, apiKey *string

		apiKey = utils.GetAPIKeyFromURL(r)
		if apiKey == nil {
			tokenString = utils.GetAuthTokenFromHeaders(r)
		}

		if apiKey == nil && tokenString == nil {
			return nil, serrors.AuthErrorf("could not get auth header or api key", nil)
		}

		if apiKey != nil {
			existingApiKey, err := a.authController.VerifyAPIKeyExists(r.Context(), *apiKey)
			if err != nil {
				return nil, err
			}
			if existingApiKey == nil {
				return nil, serrors.AuthError(nil)
			}
			ctx := context.WithValue(r.Context(), "API_KEY", existingApiKey)
			r = r.WithContext(ctx)
			fmt.Println("MADE IT OTO END OF MIDDLWARE")
		} else {
			tokenString = utils.GetAuthTokenFromHeaders(r)
			if tokenString == nil {
				return nil, serrors.AuthErrorf("could not get auth header", nil)
			}

			jwtToken, err := a.authController.VerifyJWT(*tokenString, true)
			if err != nil {
				return nil, err
			}
			ctx, err := utils.SetChatProfileToContext(jwtToken, r.Context())
			if err != nil {
				return nil, err
			}
			r = r.WithContext(ctx)
		}

		return h(w, r)
	}
}
