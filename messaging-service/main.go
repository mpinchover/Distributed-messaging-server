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
		fx.Provide(middleware.NewJWTAuthMiddleware),
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

/*
Solution
pass the middleware thing in here as a param
and then you can call GetAuthMiddleware, and run it in the array like youa re down below
*/

type SetupRoutesParams struct {
	fx.In

	AuthMiddleware    *middleware.APIKeyAuthMiddleware
	JWTAuthMiddleware *middleware.JWTAuthMiddleware
	Handler           *handlers.Handler
	Router            *mux.Router
}

// func SetupRoutes(h *handlers.Handler, r *mux.Router) {
func SetupRoutes(p SetupRoutesParams) {

	commonMiddleware := []middleware.Middleware{p.AuthMiddleware}

	// pass in all the middleware to the handler
	// handler should have methods like SetupCreateRoom
	// which will then invoke all the midleware and stuff
	deleteRoomHandler := middleware.New(p.Handler.DeleteRoom, commonMiddleware)
	leaveRoomHandler := middleware.New(p.Handler.LeaveRoom, commonMiddleware)
	createRoomHandler := middleware.New(p.Handler.CreateRoom, commonMiddleware)
	getMessagesByRoomUUIDHandler := middleware.New(p.Handler.GetMessagesByRoomUUID, commonMiddleware)
	getRoomsByUserUUIDHandler := middleware.New(p.Handler.GetRoomsByUserUUID, commonMiddleware)

	signupHandler := middleware.New(p.Handler.Signup, nil)
	loginHandler := middleware.New(p.Handler.Login, nil)

	// // probably need a subscribe fn here

	testAuthHandler := middleware.New(p.Handler.TestAuthProfileHandler, []middleware.Middleware{p.JWTAuthMiddleware})

	// websocket
	p.Router.HandleFunc("/ws", p.Handler.SetupWebsocketConnection)

	// API
	// p.Router.Handle("/test", testHandler).Methods("GET")
	p.Router.Handle("/test-auth-profile", testAuthHandler).Methods("GET")
	p.Router.Handle("/signup", signupHandler).Methods("POST")
	p.Router.Handle("/login", loginHandler).Methods("POST")
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
