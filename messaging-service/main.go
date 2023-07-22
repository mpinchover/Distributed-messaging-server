package main

import (
	"context"
	"log"
	"messaging-service/controllers/authcontroller"
	"messaging-service/controllers/channelscontroller"
	"messaging-service/controllers/connectionscontroller"
	"messaging-service/controllers/controltower"
	"messaging-service/handlers"
	"messaging-service/middleware"
	redisClient "messaging-service/redis"
	"messaging-service/repo"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

// https://markphelps.me/posts/handling-errors-in-your-http-handlers/
func main() {
	godotenv.Load()
	fx.New(
		// middleware
		fx.Provide(middleware.NewAuthProfileJWT),
		fx.Provide(middleware.NewMessagingJWT),
		fx.Provide(middleware.NewAPIKeyAuthMiddleware),
		// controllers
		fx.Provide(channelscontroller.New),
		fx.Provide(authcontroller.New),
		fx.Provide(connectionscontroller.New),
		fx.Provide(NewMuxRouter),
		fx.Provide(handlers.New),
		fx.Provide(redisClient.New),
		fx.Provide(controltower.New),
		fx.Provide(repo.New),
		fx.Invoke(SetupRoutes),
	).Run()
}

type SetupRoutesParams struct {
	fx.In

	APIKeyAuthMiddleware *middleware.APIKeyAuthMiddleware
	AuthProfileJWT       *middleware.AuthProfileJWT
	MessagingJWT         *middleware.MessagingJWT
	Handler              *handlers.Handler
	Router               *mux.Router
}

// func SetupRoutes(h *handlers.Handler, r *mux.Router) {
func SetupRoutes(p SetupRoutesParams) {

	authProfileJWTMW := []middleware.Middleware{p.AuthProfileJWT}
	messagingAuthMW := []middleware.Middleware{p.MessagingJWT}
	apiKeyAuthMW := []middleware.Middleware{p.APIKeyAuthMiddleware}

	// testing
	testAuthHandler := middleware.New(p.Handler.TestAuthProfileHandler, authProfileJWTMW)
	p.Router.Handle("/test-auth-profile", testAuthHandler).Methods("GET")

	testAuthAPIKeyHandler := middleware.New(p.Handler.TestNewAPIKeyHandler, apiKeyAuthMW)
	p.Router.Handle("/test-auth-api-key", testAuthAPIKeyHandler).Methods("GET")

	// websocket
	p.Router.HandleFunc("/ws", p.Handler.SetupWebsocketConnection)

	// auth
	signupHandler := middleware.New(p.Handler.Signup, nil)
	p.Router.Handle("/signup", signupHandler).Methods("POST")

	loginHandler := middleware.New(p.Handler.Login, nil)
	p.Router.Handle("/login", loginHandler).Methods("POST")

	passwordResetHandler := middleware.New(p.Handler.GeneratePasswordResetLink, nil)
	p.Router.Handle("/create-reset-password-link", passwordResetHandler).Methods("POST")

	apiKeyHandler := middleware.New(p.Handler.GetNewAPIKey, authProfileJWTMW)
	p.Router.Handle("/get-new-api-key", apiKeyHandler).Methods("GET")

	invalidateApiKeyHandler := middleware.New(p.Handler.InvalidateAPIKey, authProfileJWTMW)
	p.Router.Handle("/invalidate-api-key", invalidateApiKeyHandler).Methods("POST")

	updatePasswordHandler := middleware.New(p.Handler.UpdatePassword, authProfileJWTMW)
	p.Router.Handle("/update-password", updatePasswordHandler).Methods("POST")

	resetPasswordHandler := middleware.New(p.Handler.ResetPassword, nil)
	p.Router.Handle("/reset-password", resetPasswordHandler).Methods("POST")

	refreshTokenHandler := middleware.New(p.Handler.RefreshAccessToken, authProfileJWTMW)
	p.Router.Handle("/refresh-token", refreshTokenHandler).Methods("GET")

	// API
	generateMessagingTokenRoute := middleware.New(p.Handler.GenerateMessagingToken, apiKeyAuthMW)
	p.Router.Handle("/generate-messaging-token", generateMessagingTokenRoute).Methods("POST")

	deleteRoomHandler := middleware.New(p.Handler.DeleteRoom, apiKeyAuthMW)
	p.Router.Handle("/delete-room", deleteRoomHandler).Methods("POST")

	leaveRoomHandler := middleware.New(p.Handler.LeaveRoom, apiKeyAuthMW)
	p.Router.Handle("/leave-room", leaveRoomHandler).Methods("POST")

	createRoomHandler := middleware.New(p.Handler.CreateRoom, apiKeyAuthMW)
	p.Router.Handle("/create-room", createRoomHandler).Methods("POST")

	// messaging
	getRoomsByUserUUIDHandler := middleware.New(p.Handler.GetRoomsByUserUUID, messagingAuthMW)
	p.Router.Handle("/get-rooms-by-user-uuid", getRoomsByUserUUIDHandler).Methods("GET")

	getMessagesByRoomUUIDHandler := middleware.New(p.Handler.GetMessagesByRoomUUID, messagingAuthMW)
	p.Router.Handle("/get-messages-by-room-uuid", getMessagesByRoomUUIDHandler).Methods("GET")

	p.Handler.SetupChannels()
}

func NewMuxRouter(lc fx.Lifecycle) *mux.Router {
	r := mux.NewRouter()
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Println("Opening server on port 9090")
				err := http.ListenAndServe(":9090", r)
				if err != nil {
					log.Println(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
	return r
}
