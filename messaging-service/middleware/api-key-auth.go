package middleware

import (
	"context"
	"log"
	redisClient "messaging-service/redis"
	"messaging-service/serrors"
	"messaging-service/types/requests"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var keyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
	return []byte("SECRET"), nil
}

type APIKeyAuthMiddleware struct {
	redisClient *redisClient.RedisClient
}

func NewAPIKeyAuthMiddleware(
	redisClient *redisClient.RedisClient,
) *APIKeyAuthMiddleware {
	return &APIKeyAuthMiddleware{
		redisClient: redisClient,
	}
}

func (a *APIKeyAuthMiddleware) apiKeyExists(ctx context.Context, apiKey string) (bool, error) {
	var existingApiKey *bool
	err := a.redisClient.Get(ctx, apiKey, existingApiKey)
	if err != nil {
		return false, err
	}

	return existingApiKey != nil, nil
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
func (a *APIKeyAuthMiddleware) generateJWT(userID string, appID string) (string, error) {

	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["USER_ID"] = userID
	claims["APP_ID"] = appID // identify which external service this API key belongs to.
	claims["EXP"] = time.Now().UTC().Add(20 * time.Minute).Unix()
	tkn, _ := keyfunc(token)
	tokenString, err := token.SignedString(tkn)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// TODO – test this
func (a *APIKeyAuthMiddleware) verifyJWT(tokenString string, checkExp bool) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, keyfunc)
	if err != nil {
		return nil, serrors.InternalError(err)
	}
	isExpired, err := isTokenExpired(token)
	if err != nil {
		return nil, err
	}

	if checkExp && isExpired {
		return nil, serrors.AuthError(nil)
	}

	if !token.Valid {
		return nil, serrors.InternalError(nil)
	}
	return token, nil
}

func getClaimsFromJWT(jwtToken *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, serrors.InternalErrorf("Could not get claims from token", nil)
	}
	return claims, nil
}

func isTokenExpired(jwtToken *jwt.Token) (bool, error) {
	claims, err := getClaimsFromJWT(jwtToken)
	if err != nil {
		return false, err
	}

	expiration := int64(claims["EXP"].(float64))
	return expiration < time.Now().Unix(), nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (a *APIKeyAuthMiddleware) execute(h HTTPHandler) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		// if there's a token then handle the token auth
		// ctx := context.WithValue(r.Context(), "Username1", "Bob Moses")
		// fmt.Println("EXECUTING AUTH MIDDLEWARE HERE")
		// r = r.WithContext(ctx)
		params := mux.Vars(r)
		apiKey, ok := params["key"]
		if !ok || !IsValidUUID(apiKey) {
			log.Println("Unauthorized")
			return requests.MakeUnauthorized(w, "Unauthorized")
		}

		doesExist, err := a.apiKeyExists(r.Context(), apiKey)
		if err != nil {
			log.Println("Internal server")
			return requests.MakeInternalError(w, "Internal server error")
		}
		if !doesExist {
			log.Println("Unauthorized")
			return requests.MakeUnauthorized(w, "Unauthorized")
		}
		h(w, r)
		return nil, nil
	}
}

// TODO – handle expired jwts
// TODO - refresh the JWT
// TODO – create this as a struct with FX and middleware
// func (a *AuthMiddleware) AuthenticateClientConnection(h HTTPHandler) HTTPHandler {
// 	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// 		if r.Header["Authorization"] == nil {
// 			return requests.MakeUnauthorized(w, "Not Authorized")
// 		}

// 		tokenString := r.Header["Authorization"][0]
// 		jwtToken, err := a.verifyJWT(tokenString, false)
// 		if err != nil {
// 			if serrors.GetStatusCode(err) == http.StatusInternalServerError {
// 				return requests.MakeInternalError(w, "Internal server error")
// 			}
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return w.Write([]byte("Internal server error"))
// 		}

// 		claims, ok := jwtToken.Claims.(jwt.MapClaims)
// 		if !ok {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return w.Write([]byte("Internal server error"))
// 		}

// 		userID, ok := claims["USER_ID"]
// 		if !ok {
// 			return requests.MakeUnauthorized(w, "Not Authorized")
// 		}

// 		ctx := context.WithValue(r.Context(), "UserID", userID)
// 		r = r.WithContext(ctx)
// 		h(w, r)
// 		return nil, nil
// 	}
// }

func (a *APIKeyAuthMiddleware) AuthenticateClientConnection(tokenString string, checkExp bool) (*requests.User, error) {
	if tokenString == "" {
		return nil, serrors.AuthError(nil)
	}

	jwtToken, err := a.verifyJWT(tokenString, false)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, serrors.InternalError(nil)
	}

	userID, ok := claims["USER_ID"]
	if !ok {
		return nil, serrors.AuthError(nil)
	}

	return &requests.User{
		UserID: userID.(string),
	}, nil
}
