package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"messaging-service/src/serrors"
	"messaging-service/src/types/requests"
	"net/http"
	"os"
	"time"

	goerrors "github.com/go-errors/errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var Keyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
	// return []byte(os.Getenv("JWT_SECRET")), nil
	return []byte(os.Getenv("JWT_SECRET")), nil
}

func GetAPIKeyFromCtx(ctx context.Context) (*requests.APIKey, error) {
	_apiKey := ctx.Value("API_KEY")

	apiKey := &requests.APIKey{}
	b, err := json.Marshal(_apiKey)
	if err != nil {
		return nil, serrors.AuthError(err)
	}
	err = json.Unmarshal(b, apiKey)
	if err != nil {
		return nil, serrors.AuthError(err)
	}
	return apiKey, nil
}

// func GetAuthProfileFromCtx(ctx context.Context) (*requests.AuthProfile, error) {
// 	_authProfile := ctx.Value("AUTH_PROFILE")
// 	if _authProfile == nil {
// 		return nil, serrors.AuthErrorf("could not get auth profile", nil)
// 	}
// 	authProfile := &requests.AuthProfile{}
// 	b, err := json.Marshal(_authProfile)
// 	if err != nil {
// 		return nil, serrors.AuthError(err)
// 	}
// 	err = json.Unmarshal(b, authProfile)
// 	if err != nil {
// 		return nil, serrors.AuthError(err)
// 	}
// 	return authProfile, nil
// }

func GetClaimsFromJWT(jwtToken *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, serrors.InternalErrorf("Could not get claims from token", nil)
	}
	return claims, nil
}

func IsTokenExpired(jwtToken *jwt.Token) (bool, error) {
	claims, err := GetClaimsFromJWT(jwtToken)
	if err != nil {
		return false, err
	}

	expirationTime := claims["EXP"].(float64)
	now := time.Now().Unix()
	return expirationTime < float64(now), nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func GetAuthTokenFromHeaders(r *http.Request) *string {
	if r.Header["Authorization"] == nil {
		return nil
	}

	// fmt.Println("1")
	if len(r.Header["Authorization"]) == 0 {
		return nil
	}
	tokenString := r.Header["Authorization"][0]
	if tokenString == "" {
		return nil
	}
	return &tokenString
}

// TODO – regex to check if it's a valid API key
func GetAPIKeyFromURL(r *http.Request) *string {
	apiKey := r.URL.Query().Get("key")
	if apiKey == "" {
		return nil
	}
	return &apiKey
}

// func SetAuthProfileToContext(jwtToken *jwt.Token, oldContext context.Context) (context.Context, error) {
// 	claims, ok := jwtToken.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return nil, serrors.InternalErrorf("could not get token claims", nil)
// 	}
// 	authProfile, err := GetAuthProfileFromTokenClaims(claims)
// 	if err != nil {
// 		return nil, serrors.InternalError(err)
// 	}

// 	// fmt.Println("6")
// 	ctx := context.WithValue(oldContext, "AUTH_PROFILE", authProfile)
// 	return ctx, nil
// }

// func GetAuthProfileFromTokenClaims(claims jwt.MapClaims) (*requests.AuthProfile, error) {
// 	_authProfile, ok := claims["AUTH_PROFILE"]
// 	if !ok {
// 		return nil, serrors.InternalErrorf("could not get auth profile", nil)
// 	}
// 	bytes, err := json.Marshal(_authProfile)
// 	if err != nil {
// 		return nil, serrors.InternalErrorf("could not marshall auth profile", nil)
// 	}
// 	authProfile := &requests.AuthProfile{}
// 	err = json.Unmarshal(bytes, authProfile)
// 	return authProfile, err
// }

func SetChatProfileToContext(jwtToken *jwt.Token, oldContext context.Context) (context.Context, error) {
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, serrors.InternalErrorf("could not get token claims", nil)
	}

	chatProfile, err := GetChatProfileFromTokenClaims(claims)
	if err != nil {
		return nil, serrors.InternalError(err)
	}

	// fmt.Println("6")
	ctx := context.WithValue(oldContext, "USER_ID", chatProfile)
	return ctx, nil
}

func GetChatProfileFromTokenClaims(claims jwt.MapClaims) (*requests.ChatProfile, error) {
	_chatProfile, ok := claims["USER_ID"]
	if !ok {
		fmt.Println("COULD NOT GET CHAT PROFILE")
		return nil, nil
	}
	bytes, err := json.Marshal(_chatProfile)
	if err != nil {
		return nil, serrors.AuthErrorf("could not marshall chat profile", nil)
	}
	chatProfile := &requests.ChatProfile{}
	err = json.Unmarshal(bytes, chatProfile)
	return chatProfile, err
}

func GenerateMessagingToken(userUUID string, exp time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["USER_ID"] = requests.ChatProfile{
		UserUUID: userUUID,
	}
	claims["EXP"] = exp.Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", goerrors.Wrap(err, 0)
	}

	return tokenString, nil
}

// func GenerateJWTToken(authProfile *requests.AuthProfile, exp time.Time) (string, error) {
// 	token := jwt.New(jwt.SigningMethodHS256)
// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["AUTH_PROFILE"] = authProfile
// 	claims["EXP"] = exp.Unix()
// 	token.Claims = claims

// 	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
// 	if err != nil {
// 		return "", goerrors.Wrap(err, 0)
// 	}

// 	return tokenString, nil
// }

func VerifyJWT(tokenString string, checkExp bool) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, Keyfunc)
	if err != nil {
		return nil, serrors.InternalError(err)
	}
	isExpired, err := IsTokenExpired(token)
	if err != nil {
		return nil, err
	}

	if checkExp && isExpired {
		return nil, serrors.AuthErrorf("token is expired", nil)
	}

	if !token.Valid {
		return nil, serrors.InternalErrorf("token is not valid", nil)
	}
	return token, nil
}
