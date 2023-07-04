package authcontroller

import (
	"fmt"
	"messaging-service/repo"
	"messaging-service/serrors"
	"messaging-service/types/records"
	"messaging-service/types/requests"
	"messaging-service/utils"
	"time"

	goerrors "github.com/go-errors/errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	repo *repo.Repo
}

func New(
	repo *repo.Repo,
) *AuthController {
	authController := &AuthController{
		repo: repo,
	}
	return authController
}

// Hash password using the bcrypt hashing algorithm
func hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// Check if two passwords match using Bcrypt's CompareHashAndPassword
// which return nil on success and an error on failure.
func doPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func (a *AuthController) Login(req *requests.LoginRequest) (*requests.LoginResponse, error) {
	authProfile, err := a.repo.GetAuthProfileByEmail(req.Email)
	if err != nil {
		fmt.Println("POINT 1")
		return nil, serrors.InternalError(err)
	}

	if authProfile == nil {
		fmt.Println("POINT 2")
		return nil, serrors.AuthError(err)
	}

	if !doPasswordsMatch(authProfile.HashedPassword, req.Password) {
		fmt.Println("POINT 3")
		return nil, serrors.AuthErrorf("password/email does not match", nil)
	}

	rAuthProfile := requests.AuthProfile{
		Email: authProfile.Email,
		UUID:  authProfile.UUID,
	}

	token, err := utils.GenerateAPIToken(rAuthProfile)
	if err != nil {
		return nil, serrors.InternalError(err)
	}

	return &requests.LoginResponse{
		Token: token,
	}, nil
}

func GenerateJWTToken(authProfile requests.AuthProfile) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["AUTH_PROFILE"] = authProfile
	claims["EXP"] = time.Now().UTC().Add(20 * time.Minute).Unix()
	token.Claims = claims

	// tkn, _ := utils.Keyfunc(token)
	tokenString, err := token.SignedString([]byte("SECRET"))
	if err != nil {
		return "", goerrors.Wrap(err, 0)
	}

	return tokenString, nil
}

func (a *AuthController) Signup(req *requests.SignupRequest) (*requests.SignupResponse, error) {
	authProfile, err := a.repo.GetAuthProfileByEmail(req.Email)
	if err != nil {
		return nil, serrors.InternalError(err)
	}

	if authProfile != nil {
		return nil, serrors.InvalidArgumentErrorf("User already exists", nil)
	}

	authUserUUID := uuid.New().String()
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, serrors.InternalError(err)
	}

	recordAuthProfile := &records.AuthProfile{
		UUID:           authUserUUID,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	err = a.repo.SaveAuthProfile(recordAuthProfile)
	if err != nil {
		return nil, serrors.InternalError(err)
	}

	rAuthProfile := requests.AuthProfile{
		Email: recordAuthProfile.Email,
		UUID:  recordAuthProfile.UUID,
	}

	token, err := GenerateJWTToken(rAuthProfile)
	if err != nil {
		return nil, serrors.InternalError(err)
	}

	return &requests.SignupResponse{
		Token: token,
	}, nil
}

func (a *AuthController) VerifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, utils.Keyfunc)
	if err != nil {
		return nil, serrors.InternalError(err)
	}
	isExpired, err := utils.IsTokenExpired(token)
	if err != nil {
		return nil, err
	}

	if isExpired {
		return nil, serrors.AuthErrorf("token is expired", nil)
	}

	if !token.Valid {
		return nil, serrors.InternalErrorf("token is not valid", nil)
	}
	return token, nil
}
