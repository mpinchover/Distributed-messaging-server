package middleware

import (
	"context"
	"encoding/json"
	"messaging-service/controllers/authcontroller"
	"messaging-service/serrors"
	"messaging-service/types/requests"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type JWTAuthMiddleware struct {
	authController *authcontroller.AuthController
}

func NewJWTAuthMiddleware(authController *authcontroller.AuthController) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
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

func (a *JWTAuthMiddleware) execute(h HTTPHandler) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		if r.Header["Authorization"] == nil {
			return nil, serrors.AuthErrorf("missing auth header", nil)
		}

		tokenString := r.Header["Authorization"][0]
		jwtToken, err := a.authController.VerifyJWT(tokenString)
		if err != nil {
			return nil, err
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return nil, serrors.InternalErrorf("could not get token claims", nil)
		}
		_authProfile, ok := claims["AUTH_PROFILE"]
		if !ok {
			return nil, serrors.AuthErrorf("could not get auth profile", nil)
		}

		bytes, err := json.Marshal(_authProfile)
		if err != nil {
			return nil, serrors.AuthErrorf("could not marshall auth profile", nil)
		}

		authProfile := &requests.AuthProfile{}
		err = json.Unmarshal(bytes, authProfile)
		if err != nil {
			return nil, serrors.AuthErrorf("Not Authorized", nil)
		}

		ctx := context.WithValue(r.Context(), "AUTH_PROFILE", authProfile)
		r = r.WithContext(ctx)
		return h(w, r)
	}
}
