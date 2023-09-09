package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"messaging-service/src/controllers/authcontroller"
	"messaging-service/src/controllers/controltower"
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

	RedisClient *redisClient.RedisClient
	// AuthMiddleware *middleware.AuthMiddleware
	ControlTower   *controltower.ControlTowerCtrlr
	AuthController *authcontroller.AuthController
}

func New(p Params) *Handler {
	return &Handler{
		ControlTowerCtrlr: p.ControlTower,
		RedisClient:       p.RedisClient,
		AuthController:    p.AuthController,
		// AuthMiddleware:    p.AuthMiddleware,
	}
}

var decoder = schema.NewDecoder()

func (h *Handler) GetRoomsByUserUUID(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// validation
	req := &requests.GetRoomsByUserUUIDRequest{}
	err := decoder.Decode(req, r.URL.Query())
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	return h.getRoomsByUserUUID(ctx, req)
}

func (h *Handler) APIGetRoomsByUserUUID(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// validation
	req := &requests.GetRoomsByUserUUIDRequest{}
	err := decoder.Decode(req, r.URL.Query())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	ctx := r.Context()
	return h.getRoomsByUserUUID(ctx, req)
}

func (h *Handler) getRoomsByUserUUID(ctx context.Context, req *requests.GetRoomsByUserUUIDRequest) (*requests.GetRoomsByUserUUIDResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		fmt.Println("ERR 1")
		fmt.Println(err)
		return nil, err
	}

	rooms, err := h.ControlTowerCtrlr.GetRoomsByUserUUID(ctx, req.UserUUID, req.Offset)
	if err != nil {
		fmt.Println("ERR 2")
		fmt.Println(err)
		return nil, err
	}

	response := &requests.GetRoomsByUserUUIDResponse{
		Rooms: rooms,
	}
	return response, nil
}

func (h *Handler) GetMessagesByRoomUUID(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// validation
	req := &requests.GetMessagesByRoomUUIDRequest{}
	err := decoder.Decode(req, r.URL.Query())
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	return h.getMessagesByRoomUUID(ctx, req)
}

func (h *Handler) APIGetMessagesByRoomUUID(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// validation
	req := &requests.GetMessagesByRoomUUIDRequest{}
	err := decoder.Decode(req, r.URL.Query())
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	return h.getMessagesByRoomUUID(ctx, req)
}

func (h *Handler) getMessagesByRoomUUID(ctx context.Context, req *requests.GetMessagesByRoomUUIDRequest) (*requests.GetMessagesByRoomUUIDResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	msgs, err := h.ControlTowerCtrlr.GetMessagesByRoomUUID(ctx, req.RoomUUID, req.Offset)
	if err != nil {
		return nil, err
	}

	requestMsgs := make([]*requests.Message, len(msgs))
	for i, msg := range msgs {

		seenBy := make([]*requests.SeenBy, len(msg.SeenBy))
		for j, sb := range msg.SeenBy {
			seenBy[j] = &requests.SeenBy{
				MessageUUID: sb.MessageUUID,
				UserUUID:    sb.UserUUID,
			}
		}

		requestMsgs[i] = &requests.Message{
			UUID:          msg.UUID,
			FromUUID:      msg.FromUUID,
			RoomUUID:      msg.RoomUUID,
			MessageText:   msg.MessageText,
			CreatedAtNano: msg.CreatedAtNano,
			SeenBy:        seenBy,
		}
	}

	resp := &requests.GetMessagesByRoomUUIDResponse{
		Messages: requestMsgs,
	}
	return resp, nil
}

func (h *Handler) DeleteRoom(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// TODO - validation
	req := &requests.DeleteRoomRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	ctx := r.Context()
	return h.deleteRoom(ctx, req)
}

func (h *Handler) deleteRoom(ctx context.Context, req *requests.DeleteRoomRequest) (*requests.DeleteRoomResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	roomUUID := req.RoomUUID
	// verify user has permissions
	err = h.ControlTowerCtrlr.DeleteRoom(ctx, roomUUID)
	if err != nil {
		return nil, err
	}
	return &requests.DeleteRoomResponse{}, nil
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// validation
	// run all the middleware here

	req := &requests.CreateRoomRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	ctx := r.Context()
	return h.createRoom(ctx, req)
}

func (h *Handler) createRoom(ctx context.Context, req *requests.CreateRoomRequest) (*requests.CreateRoomResponse, error) {
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

func (h *Handler) LeaveRoom(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// validation
	req := &requests.LeaveRoomRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	ctx := r.Context()
	return h.leaveRoom(ctx, req)
}

func (h *Handler) leaveRoom(ctx context.Context, req *requests.LeaveRoomRequest) (*requests.LeaveRoomResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	err = h.ControlTowerCtrlr.LeaveRoom(ctx, req.UserUUID, req.RoomUUID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// func (h *Handler) Login(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	// validation
// 	req := &requests.LoginRequest{}
// 	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
// 		return nil, err
// 	}
// 	ctx := r.Context()
// 	return h.login(ctx, req)
// }

// func (h *Handler) login(ctx context.Context, req *requests.LoginRequest) (*requests.LoginResponse, error) {
// 	err := validation.ValidateRequest(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// login user
// 	// send back token
// 	return h.AuthController.Login(req)
// }

// func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	// validation
// 	req := &requests.SignupRequest{}
// 	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
// 		return nil, err
// 	}
// 	ctx := r.Context()
// 	return h.signup(ctx, req)
// }

// func (h *Handler) signup(ctx context.Context, req *requests.SignupRequest) (*requests.SignupResponse, error) {
// 	err := validation.ValidateRequest(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return h.AuthController.Signup(req)
// }

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
// 	return &requests.RefreshAccessTokenResponse{
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 	}, nil
// }

// func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	ctx := r.Context()
// 	req := &requests.UpdatePasswordRequest{}
// 	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
// 		return nil, err
// 	}
// 	return h.updatePassword(ctx, req)
// }

// func (h *Handler) updatePassword(ctx context.Context, req *requests.UpdatePasswordRequest) (*requests.GenericResponse, error) {
// 	err := h.AuthController.UpdatePassword(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &requests.GenericResponse{
// 		Success: true,
// 	}, nil
// }

// func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	ctx := r.Context()
// 	req := &requests.ResetPasswordRequest{}
// 	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
// 		return nil, err
// 	}
// 	return h.resetPassword(ctx, req)
// }

// func (h *Handler) resetPassword(ctx context.Context, req *requests.ResetPasswordRequest) (*requests.GenericResponse, error) {
// 	err := h.AuthController.ResetPassword(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &requests.GenericResponse{
// 		Success: true,
// 	}, nil
// }

// func (h *Handler) GeneratePasswordResetLink(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	ctx := r.Context()
// 	req := &requests.GeneratePasswordResetLinkRequest{}
// 	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
// 		return nil, err
// 	}
// 	return h.generatePasswordResetLink(ctx, req)
// }

// func (h *Handler) generatePasswordResetLink(ctx context.Context, req *requests.GeneratePasswordResetLinkRequest) (*requests.GenericResponse, error) {
// 	err := h.AuthController.GeneratePasswordResetLink(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &requests.GenericResponse{
// 		Success: true,
// 	}, nil
// }

func (h *Handler) GenerateMessagingToken(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	req := &requests.GenerateMessagingTokenRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	return h.generateMessagingToken(ctx, req)
}

// create a user for this person
// have they exceeded quota for monthly active users?
func (h *Handler) generateMessagingToken(ctx context.Context, req *requests.GenerateMessagingTokenRequest) (*requests.GenerateMessagingTokenResponse, error) {
	userUUID := req.UserUUID
	accessToken, err := utils.GenerateMessagingToken(userUUID, time.Now().Add(10*time.Minute))
	if err != nil {
		return nil, err
	}

	return &requests.GenerateMessagingTokenResponse{
		Token: accessToken,
	}, nil
}
