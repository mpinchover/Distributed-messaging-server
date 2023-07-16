package requests

type GetRoomsByUserUUIDRequest struct {
	UserUUID string `schema:"userUuid" validate:"required"`
	Offset   int    `schema:"offset"`
}

type GetRoomsByUserUUIDResponse struct {
	Rooms []*Room `json:"rooms"`
}

type GetMessagesByRoomUUIDRequest struct {
	RoomUUID string `schema:"roomUuid" validate:"required"`
	Offset   int    `schema:"offset"`
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
}

type InvalidateAPIKeyRequest struct {
	Key string
}

type RefreshAccessTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UpdatePasswordRequest struct {
	CurrentPassword    string `json:"currentPassword"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}

type GeneratePasswordResetLinkRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token              string `json:"token"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}

type GenericResponse struct {
	Success bool
}

type ErrorResponse struct {
	Message string
}
