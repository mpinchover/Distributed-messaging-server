package main

import (
	"context"
	"encoding/json"
	"log"
	"messaging-service/src/controllers/authcontroller"
	"messaging-service/src/route"

	"messaging-service/src/controllers/controltower"
	"messaging-service/src/handlers"
	"messaging-service/src/middleware"
	redisClient "messaging-service/src/redis"
	"messaging-service/src/repo"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	godotenv.Load()
	fx.New(
		// middleware
		fx.Provide(middleware.NewAPIKeyAuthMiddleware),

		// controllers
		fx.Provide(authcontroller.New),
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
	Handler              *handlers.Handler
	Router               *mux.Router
}

func SetupRoutes(p SetupRoutesParams) {

	apiKeyAuthMW := []middleware.Middleware{p.APIKeyAuthMiddleware}

	testAuthAPIKeyHandler := route.New(p.Handler.TestNewAPIKeyHandler, apiKeyAuthMW)
	p.Router.Handle("/test-auth-api-key", testAuthAPIKeyHandler).Methods("GET")

	// websocket
	p.Router.HandleFunc("/ws", p.Handler.SetupWebsocketConnection)

	// API
	generateMessagingTokenRoute := route.New(p.Handler.GenerateMessagingToken, apiKeyAuthMW)
	p.Router.Handle("/generate-messaging-token", generateMessagingTokenRoute).Methods("POST")

	p.Router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Message string
		}{
			Message: "pong",
		}
		b, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		w.Write(b)
	})

	deleteRoomHandler := route.New(p.Handler.DeleteRoom, apiKeyAuthMW)
	p.Router.Handle("/delete-room", deleteRoomHandler).Methods("POST")

	createRoomHandler := route.New(p.Handler.CreateRoom, apiKeyAuthMW)
	p.Router.Handle("/create-room", createRoomHandler).Methods("POST")

	getRoomsByUserUUIDHandler := route.New(p.Handler.GetRoomsByUserUUID, apiKeyAuthMW)
	p.Router.Handle("/get-rooms-by-user-uuid", getRoomsByUserUUIDHandler).Methods("GET")

	getMessagesByRoomUUIDHandler := route.New(p.Handler.GetMessagesByRoomUUID, apiKeyAuthMW)
	p.Router.Handle("/get-messages-by-room-uuid", getMessagesByRoomUUIDHandler).Methods("GET")

	getUserConnection := route.New(p.Handler.GetUserConnection, nil)
	p.Router.Handle("/get-user-connection/{userUuid}", getUserConnection).Methods("GET")

	getChannel := route.New(p.Handler.GetChannel, nil)
	p.Router.Handle("/get-channel/{channelUuid}", getChannel).Methods("GET")

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
