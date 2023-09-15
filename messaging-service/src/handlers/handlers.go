package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"messaging-service/src/controllers/authcontroller"
	"messaging-service/src/controllers/controltower"
	mappers "messaging-service/src/mappers/requests"
	redisClient "messaging-service/src/redis"
	"messaging-service/src/types/requests"
	"messaging-service/src/utils"
	"messaging-service/src/validation"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	"go.uber.org/fx"
)

type Handler struct {
	ControlTowerCtrlr *controltower.ControlTowerCtrlr
	AuthController    *authcontroller.AuthController
	RedisClient       *redisClient.RedisClient
}

type Params struct {
	fx.In

	RedisClient    *redisClient.RedisClient
	ControlTower   *controltower.ControlTowerCtrlr
	AuthController *authcontroller.AuthController
}

func New(p Params) *Handler {
	return &Handler{
		ControlTowerCtrlr: p.ControlTower,
		RedisClient:       p.RedisClient,
		AuthController:    p.AuthController,
	}
}

var decoder = schema.NewDecoder()

func (h *Handler) GetRoomsByUserUUID(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	req := &requests.GetRoomsByUserUUIDRequest{}
	err := decoder.Decode(req, r.URL.Query())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	rooms, err := h.ControlTowerCtrlr.GetRoomsByUserUUID(r.Context(), req.UserUUID, req.Offset)
	if err != nil {
		return nil, err
	}

	response := &requests.GetRoomsByUserUUIDResponse{
		Rooms: mappers.ToRequestRooms(rooms),
	}
	return response, nil
}

func (h *Handler) GetMessagesByRoomUUID(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	req := &requests.GetMessagesByRoomUUIDRequest{}
	err := decoder.Decode(req, r.URL.Query())
	if err != nil {
		return nil, err
	}

	err = validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	msgs, err := h.ControlTowerCtrlr.GetMessagesByRoomUUID(r.Context(), req.RoomUUID, req.Offset)
	if err != nil {
		return nil, err
	}

	resp := &requests.GetMessagesByRoomUUIDResponse{
		Messages: mappers.ToRequestMessages(msgs),
	}
	return resp, nil
}

func (h *Handler) DeleteRoom(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	req := &requests.DeleteRoomRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	err = h.ControlTowerCtrlr.DeleteRoom(r.Context(), req.RoomUUID)
	if err != nil {
		return nil, err
	}
	return &requests.DeleteRoomResponse{}, nil
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// TODO - validate the create room request
	req := &requests.CreateRoomRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	ctx := r.Context()

	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	room, err := h.ControlTowerCtrlr.CreateRoom(ctx, req.Members)
	if err != nil {
		return nil, err
	}

	return &requests.CreateRoomResponse{
		Room: room,
	}, nil
}

func (h *Handler) GetNewAPIKey(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	return h.getNewAPIKey(ctx)
}

func (h *Handler) getNewAPIKey(ctx context.Context) (*requests.APIKey, error) {
	key, err := h.AuthController.GenerateAPIKey(ctx)
	if err != nil {
		return nil, err
	}
	return &requests.APIKey{
		Key: key,
	}, nil
}

func (h *Handler) InvalidateAPIKey(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()

	req := &requests.InvalidateAPIKeyRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	return h.invalidateAPIKey(ctx, req)
}

func (h *Handler) invalidateAPIKey(ctx context.Context, req *requests.InvalidateAPIKeyRequest) (*requests.GenericResponse, error) {
	err := h.AuthController.RemoveAPIKey(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return &requests.GenericResponse{
		Success: true,
	}, nil
}

// func (h *Handler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	ctx := r.Context()
// 	return h.refreshAccessToken(ctx)
// }

// func (h *Handler) refreshAccessToken(ctx context.Context) (*requests.RefreshAccessTokenResponse, error) {
// 	authProfile, err := utils.GetAuthProfileFromCtx(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	accessToken, err := utils.GenerateJWTToken(authProfile, time.Now().Add(time.Minute*10))
// 	if err != nil {
// 		return nil, err
// 	}
// 	refreshToken, err := utils.GenerateJWTToken(authProfile, time.Now().Add(time.Hour*utils.NumberOfHoursInSixMonths))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &records.RefreshAccessTokenResponse{
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 	}, nil
// }

func (h *Handler) GenerateMessagingToken(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	req := &requests.GenerateMessagingTokenRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateMessagingToken(req.UserUUID, time.Now().Add(10*time.Minute))
	if err != nil {
		return nil, err
	}

	return &requests.GenerateMessagingTokenResponse{
		Token: accessToken,
	}, nil
}
