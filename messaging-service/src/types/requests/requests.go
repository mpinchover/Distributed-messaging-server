package requests

type GetRoomsByUserUUIDRequest struct {
	UserUUID string `schema:"userUuid" validate:"required"`
	Offset   int    `schema:"offset"`
	Key      string `schema:"key,-"`
}

type GetRoomsByUserUUIDResponse struct {
	Rooms []*Room `json:"rooms"`
}

type GetMessagesByRoomUUIDRequest struct {
	RoomUUID string `schema:"roomUuid" validate:"required"`
	Offset   int    `schema:"offset"`
	Key      string `schema:"key,-"`
}

type GetMessagesByRoomUUIDResponse struct {
	Messages []*Message `json:"messages"`
}

type CreateRoomRequest struct {
	Members []*Member `json:"participants"`
}

type CreateRoomResponse struct {
	Room *Room `json:"room"`
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
	UserID string `json:"userId" validate:"required"`
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
