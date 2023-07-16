package handlers

import (
	"context"
	"messaging-service/controllers/authcontroller"
	"messaging-service/controllers/controltower"
	redisClient "messaging-service/redis"
	"messaging-service/types/requests"
	"messaging-service/utils"
	"messaging-service/validation"

	"go.uber.org/fx"
)

type Handler struct {
	ControlTowerCtrlr *controltower.ControlTowerCtrlr
	AuthController    *authcontroller.AuthController
	RedisClient       *redisClient.RedisClient

	// middleware
	// AuthMiddleware *middleware.AuthMiddleware
}

type Params struct {
	fx.In

	RedisClient *redisClient.RedisClient
	// AuthMiddleware *middleware.AuthMiddleware
	ControlTower   *controltower.ControlTowerCtrlr
	AuthController *authcontroller.AuthController
}

// func New(redisClient *redisClient.RedisClient, controlTower *controltower.ControlTowerCtrlr) *Handler {
// 	return &Handler{
// 		ControlTowerCtrlr: controlTower,
// 		RedisClient:       redisClient,
// 	}
// }

func New(p Params) *Handler {
	return &Handler{
		ControlTowerCtrlr: p.ControlTower,
		RedisClient:       p.RedisClient,
		AuthController:    p.AuthController,
		// AuthMiddleware:    p.AuthMiddleware,
	}
}

func (h *Handler) getRoomsByUserUUID(ctx context.Context, req *requests.GetRoomsByUserUUIDRequest) (*requests.GetRoomsByUserUUIDResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	rooms, err := h.ControlTowerCtrlr.GetRoomsByUserUUID(ctx, req.UserUUID, req.Offset)
	if err != nil {
		return nil, err
	}

	// TODO - put this all in the controller
	response := &requests.GetRoomsByUserUUIDResponse{
		Rooms: rooms,
	}
	return response, nil
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
			UUID:        msg.UUID,
			FromUUID:    msg.FromUUID,
			RoomUUID:    msg.RoomUUID,
			MessageText: msg.MessageText,
			CreatedAt:   msg.Model.CreatedAt.UnixMilli(),
			SeenBy:      seenBy,
		}
	}

	resp := &requests.GetMessagesByRoomUUIDResponse{
		Messages: requestMsgs,
	}
	return resp, nil
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

func (h *Handler) login(ctx context.Context, req *requests.LoginRequest) (*requests.LoginResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// login user
	// send back token
	return h.AuthController.Login(req)
}

func (h *Handler) signup(ctx context.Context, req *requests.SignupRequest) (*requests.SignupResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	return h.AuthController.Signup(req)
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

func (h *Handler) invalidateAPIKey(ctx context.Context, req *requests.InvalidateAPIKeyRequest) (*requests.GenericResponse, error) {
	err := h.AuthController.RemoveAPIKey(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return &requests.GenericResponse{
		Success: true,
	}, nil
}

func (h *Handler) refreshAccessToken(ctx context.Context) (*requests.RefreshAccessTokenResponse, error) {
	authProfile, err := utils.GetAuthProfileFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	accessToken, err := h.AuthController.GenerateJWTAccessToken(authProfile)
	if err != nil {
		return nil, err
	}
	refreshToken, err := h.AuthController.GenerateJWTRefreshToken(authProfile)
	if err != nil {
		return nil, err
	}
	return &requests.RefreshAccessTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *Handler) updatePassword(ctx context.Context, req *requests.UpdatePasswordRequest) (*requests.GenericResponse, error) {
	err := h.AuthController.UpdatePassword(ctx, req)
	if err != nil {
		return nil, err
	}
	return &requests.GenericResponse{
		Success: true,
	}, nil
}

func (h *Handler) generatePasswordResetLink(ctx context.Context, req *requests.GeneratePasswordResetLinkRequest) (*requests.GenericResponse, error) {
	err := h.AuthController.GeneratePasswordResetLink(ctx, req)
	if err != nil {
		return nil, err
	}
	return &requests.GenericResponse{
		Success: true,
	}, nil
}
