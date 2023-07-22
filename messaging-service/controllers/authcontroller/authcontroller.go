package authcontroller

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	redisClient "messaging-service/redis"
	"messaging-service/repo"
	"messaging-service/serrors"
	"messaging-service/types/records"
	"messaging-service/types/requests"
	"messaging-service/utils"
	"net/http"
	"net/smtp"
	"os"
	"time"

	goerrors "github.com/go-errors/errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	repo        *repo.Repo
	redisClient *redisClient.RedisClient
}

func New(
	repo *repo.Repo,
	redisClient *redisClient.RedisClient,
) *AuthController {
	authController := &AuthController{
		repo:        repo,
		redisClient: redisClient,
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
		return nil, serrors.InternalError(err)
	}

	if authProfile == nil {
		return nil, serrors.AuthError(err)
	}

	if !doPasswordsMatch(authProfile.HashedPassword, req.Password) {
		return nil, serrors.AuthErrorf("password/email does not match", nil)
	}

	rAuthProfile := &requests.AuthProfile{
		Email: authProfile.Email,
		UUID:  authProfile.UUID,
	}

	accessToken, err := a.GenerateJWTAccessToken(rAuthProfile)
	if err != nil {
		return nil, serrors.InternalError(err)
	}
	refreshToken, err := a.GenerateJWTRefreshToken(rAuthProfile)
	if err != nil {
		return nil, serrors.InternalError(err)
	}

	return &requests.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthController) UpdatePassword(
	ctx context.Context,
	req *requests.UpdatePasswordRequest,
) error {
	authProfile, err := utils.GetAuthProfileFromCtx(ctx)
	if err != nil {
		return serrors.InternalError(err)
	}

	existingAuthProfile, err := a.repo.GetAuthProfileByEmail(authProfile.Email)
	if err != nil {
		return serrors.InternalError(err)
	}

	if existingAuthProfile == nil {
		return serrors.AuthError(err)
	}
	if !doPasswordsMatch(existingAuthProfile.HashedPassword, req.CurrentPassword) {
		return serrors.AuthErrorf("password/email does not match", nil)
	}

	// validate the new and confirm password match
	// update the password
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		return serrors.InternalError(err)
	}

	// update the authprofile with the new password
	err = a.repo.UpdatePassword(authProfile.Email, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

// https://www.loginradius.com/blog/engineering/sending-emails-with-golang/#:~:text=Below%20is%20the%20complete%20code,%2C%20%7D%20%2F%2F%20smtp%20server%20configuration.
func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// validate passwordmatch
func (a *AuthController) ResetPassword(ctx context.Context, req *requests.ResetPasswordRequest) error {
	var email string
	err := a.redisClient.Get(ctx, req.Token, &email)
	if err != nil && serrors.GetStatusCode(err) == http.StatusBadRequest {
		return serrors.AuthError(nil)
	}

	if email == "" {
		return serrors.AuthError(nil)
	}

	err = a.redisClient.Del(ctx, req.Token)
	if err != nil {
		return serrors.InternalError(nil)
	}

	existingAuthProfile, err := a.repo.GetAuthProfileByEmail(email)
	if err != nil {
		return serrors.InternalError(err)
	}

	if existingAuthProfile == nil {
		return serrors.AuthError(err)
	}

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		return serrors.InternalError(err)
	}
	err = a.repo.UpdatePassword(email, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

// use redis to store some amount of random strings.
// you dont want ta person to be able to use this twice.
// so bette rto use redis for this.
func (a *AuthController) GeneratePasswordResetLink(ctx context.Context, req *requests.GeneratePasswordResetLinkRequest) error {
	existingAuthProfile, err := a.repo.GetAuthProfileByEmail(req.Email)
	if err != nil {
		return serrors.InternalError(err)
	}

	if existingAuthProfile == nil {
		return serrors.AuthError(err)
	}

	// return a.GenerateJWTToken(authProfile, 10*time.Minute)
	token, err := generateRandomString(40)
	if err != nil {
		return err
	}

	// save token with email to identify who it is
	err = a.redisClient.SetWithTTL(ctx, token, req.Email, time.Minute*15)
	if err != nil {
		return err
	}

	from := os.Getenv("SMTP_EMAIL")        // from email
	password := os.Getenv("SMTP_PASSWORD") // from email password

	// Receiver email address.
	to := []string{
		"rachel.silverstein.applepie@gmail.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// TODO – change dir when making the service the root dir
	t, err := template.ParseFiles(wd + "/assets/templates/template.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t.Execute(&body, struct {
		ResetLink string
	}{
		ResetLink: fmt.Sprintf("http://localhost:9090/reset-password/%s", token),
	})

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (a *AuthController) GenerateMessagingToken(userID string, dur time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["USER_ID"] = requests.ChatProfile{
		UserID: userID,
	}
	claims["EXP"] = time.Now().UTC().Add(dur).Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte("SECRET"))
	if err != nil {
		return "", goerrors.Wrap(err, 0)
	}

	return tokenString, nil
}

func (a *AuthController) GenerateJWTToken(authProfile *requests.AuthProfile, dur time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["AUTH_PROFILE"] = authProfile
	claims["EXP"] = time.Now().UTC().Add(dur).Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte("SECRET"))
	if err != nil {
		return "", goerrors.Wrap(err, 0)
	}

	return tokenString, nil
}

func (a *AuthController) GenerateAPIKey(ctx context.Context) (string, error) {
	key := uuid.New().String()
	apiKey := requests.APIKey{
		Key: key,
	}

	err := a.redisClient.Set(ctx, key, apiKey)
	if err != nil {
		return "", nil
	}
	return key, nil
}

func (a *AuthController) VerifyAPIKeyExists(ctx context.Context, key string) (*requests.APIKey, error) {
	apiKey := &requests.APIKey{}

	err := a.redisClient.Get(ctx, key, apiKey)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
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

	rAuthProfile := &requests.AuthProfile{
		Email: recordAuthProfile.Email,
		UUID:  recordAuthProfile.UUID,
	}

	accessToken, err := a.GenerateJWTAccessToken(rAuthProfile)
	if err != nil {
		return nil, serrors.InternalError(err)
	}
	refreshToken, err := a.GenerateJWTRefreshToken(rAuthProfile)
	if err != nil {
		return nil, serrors.InternalError(err)
	}

	return &requests.SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UUID:         authUserUUID,
		Email:        recordAuthProfile.Email,
	}, nil
}

func (a *AuthController) GenerateJWTAccessToken(authProfile *requests.AuthProfile) (string, error) {
	return a.GenerateJWTToken(authProfile, 20*time.Minute)
}

func (a *AuthController) GenerateJWTRefreshToken(authProfile *requests.AuthProfile) (string, error) {
	return a.GenerateJWTToken(authProfile, utils.NumberOfHoursInSixMonths*time.Minute)
}

func (a *AuthController) RemoveAPIKey(ctx context.Context, apiKey string) error {
	fmt.Println("REMOVING THE API KEY ", apiKey)

	err := a.redisClient.Del(ctx, apiKey)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthController) VerifyJWT(tokenString string, checkExp bool) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, utils.Keyfunc)
	if err != nil {
		return nil, serrors.InternalError(err)
	}
	isExpired, err := utils.IsTokenExpired(token)
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
