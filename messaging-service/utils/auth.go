package utils

import (
	"context"
	"encoding/json"
	"messaging-service/serrors"
	"messaging-service/types/requests"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var Keyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
	// return []byte(os.Getenv("JWT_SECRET")), nil
	return []byte("SECRET"), nil
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

func GetAuthProfileFromCtx(ctx context.Context) (*requests.AuthProfile, error) {
	_authProfile := ctx.Value("AUTH_PROFILE")

	authProfile := &requests.AuthProfile{}
	b, err := json.Marshal(_authProfile)
	if err != nil {
		return nil, serrors.AuthError(err)
	}
	err = json.Unmarshal(b, authProfile)
	if err != nil {
		return nil, serrors.AuthError(err)
	}
	return authProfile, nil
}

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

	expiration := int64(claims["EXP"].(float64))
	return expiration < time.Now().Unix(), nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
