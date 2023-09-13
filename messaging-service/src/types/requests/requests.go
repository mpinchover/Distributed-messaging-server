package requests

import (
	"messaging-service/src/types/connections"
	"messaging-service/src/types/records"
)

type GetRoomsByUserUUIDRequest struct {
	UserUUID string `schema:"userUuid" validate:"required"`
	Offset   int    `schema:"offset"`
	Key      string `schema:"key,-"`
}

type GetRoomsByUserUUIDResponse struct {
	Rooms []*records.Room `json:"rooms"`
}

type GetMessagesByRoomUUIDRequest struct {
	RoomUUID string `schema:"roomUuid" validate:"required"`
	Offset   int    `schema:"offset"`
	Key      string `schema:"key,-"`
}

type GetUserConnectionRequest struct {
	UserUUID string
}

type GetUserConnectionResponse struct {
	UserConnections map[string]*connections.UserConnection
}

type GetChannelRequest struct {
	ChannelUUID string
}

type GetChannelResponse struct {
	Channel map[string]*connections.Channel
}

type GetMessagesByRoomUUIDResponse struct {
	Messages []*records.Message `json:"messages"`
}

type CreateRoomRequest struct {
	Members []*records.Member `json:"participants"`
}

type CreateRoomResponse struct {
	Room *records.Room `json:"room"`
}

type DeleteRoomRequest struct {
	RoomUUID string `json:"roomUuid" validate:"required"`
	// UserUUID string `json:"userUuid" validate:"required"`
}

type DeleteRoomResponse struct {
}

type LeaveRoomRequest struct {
	UserUUID string `json:"userUuid" validate:"required"`
	RoomUUID string `json:"roomUuid" validate:"required"`
}

type LeaveRoomResponse struct {
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignupRequest struct {
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

type SignupResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UUID         string `json:"uuid"`
	Email        string `json:"email" validate:"required"`
}

type InvalidateAPIKeyRequest struct {
	Key string `validate:"required"`
}

type RefreshAccessTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UpdatePasswordRequest struct {
	CurrentPassword    string `json:"currentPassword" validate:"required"`
	NewPassword        string `json:"newPassword" validate:"required"`
	ConfirmNewPassword string `json:"confirmNewPassword" validate:"required"`
}

type GeneratePasswordResetLinkRequest struct {
	Email string `json:"email" validate:"required"`
}

type ResetPasswordRequest struct {
	Token              string `json:"token" validate:"required"`
	NewPassword        string `json:"newPassword" validate:"required"`
	ConfirmNewPassword string `json:"confirmNewPassword" validate:"required"`
}

type GenerateMessagingTokenRequest struct {
	UserUUID string `json:"userUuid" validate:"required"`
}

type GenerateMessagingTokenResponse struct {
	Token string `json:"token"`
}

type GenericResponse struct {
	Success bool
}

type ErrorResponse struct {
	Message string
}
