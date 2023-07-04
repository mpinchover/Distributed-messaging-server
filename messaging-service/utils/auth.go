package utils

import (
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

func GenerateAPIToken(authProfile requests.AuthProfile) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["AUTH_PROFILE"] = authProfile
	claims["EXP"] = time.Now().UTC().Add(20 * time.Minute).Unix()
	tkn, _ := Keyfunc(token)
	tokenString, err := token.SignedString(tkn)
	if err != nil {
		return "", serrors.InternalError(err)
	}

	return tokenString, nil
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
