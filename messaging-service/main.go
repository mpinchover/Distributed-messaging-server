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
		fx.Provide(middleware.NewAccessJWTAuthMiddleware),
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

	APIKeyAuthMiddleware    *middleware.APIKeyAuthMiddleware
	AccessJWTAuthMiddleware *middleware.AccessJWTAuthMiddleware
	Handler                 *handlers.Handler
	Router                  *mux.Router
}

// func SetupRoutes(h *handlers.Handler, r *mux.Router) {
func SetupRoutes(p SetupRoutesParams) {

	// testing
	testAuthHandler := middleware.New(p.Handler.TestAuthProfileHandler, []middleware.Middleware{p.AccessJWTAuthMiddleware})
	testAuthAPIKeyHandler := middleware.New(p.Handler.TestNewAPIKeyHandler, []middleware.Middleware{p.APIKeyAuthMiddleware})
	p.Router.Handle("/test-auth-profile", testAuthHandler).Methods("GET")
	p.Router.Handle("/test-auth-api-key", testAuthAPIKeyHandler).Methods("GET")

	// websocket
	p.Router.HandleFunc("/ws", p.Handler.SetupWebsocketConnection)

	// auth
	signupHandler := middleware.New(p.Handler.Signup, nil)
	loginHandler := middleware.New(p.Handler.Login, nil)
	passwordResetHandler := middleware.New(p.Handler.GeneratePasswordResetLink, nil)
	apiKeyHandler := middleware.New(p.Handler.GetNewAPIKey, []middleware.Middleware{p.AccessJWTAuthMiddleware})
	invalidateApiKeyHandler := middleware.New(p.Handler.InvalidateAPIKey, []middleware.Middleware{p.AccessJWTAuthMiddleware})
	refreshTokenHandler := middleware.New(p.Handler.RefreshAccessToken, []middleware.Middleware{p.AccessJWTAuthMiddleware})
	updatePasswordHandler := middleware.New(p.Handler.UpdatePassword, []middleware.Middleware{p.AccessJWTAuthMiddleware})
	p.Router.Handle("/create-reset-password-link", passwordResetHandler).Methods("POST")
	p.Router.Handle("/get-new-api-key", apiKeyHandler).Methods("GET")
	p.Router.Handle("/refresh-token", refreshTokenHandler).Methods("GET")
	p.Router.Handle("/invalidate-api-key", invalidateApiKeyHandler).Methods("POST")
	p.Router.Handle("/signup", signupHandler).Methods("POST")
	p.Router.Handle("/login", loginHandler).Methods("POST")
	p.Router.Handle("/update-password", updatePasswordHandler).Methods("POST")

	// API
	deleteRoomHandler := middleware.New(p.Handler.DeleteRoom, nil)
	leaveRoomHandler := middleware.New(p.Handler.LeaveRoom, nil)
	createRoomHandler := middleware.New(p.Handler.CreateRoom, nil)
	getMessagesByRoomUUIDHandler := middleware.New(p.Handler.GetMessagesByRoomUUID, nil)
	getRoomsByUserUUIDHandler := middleware.New(p.Handler.GetRoomsByUserUUID, nil)
	p.Router.Handle("/delete-room", deleteRoomHandler).Methods("POST")
	p.Router.Handle("/leave-room", leaveRoomHandler).Methods("POST")
	p.Router.Handle("/create-room", createRoomHandler).Methods("POST")
	p.Router.Handle("/get-messages-by-room-uuid", getMessagesByRoomUUIDHandler).Methods("GET")
	p.Router.Handle("/get-rooms-by-user-uuid", getRoomsByUserUUIDHandler).Methods("GET")

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
