package handlers

import (
	"context"
	"messaging-service/src/types/requests"
	"messaging-service/src/utils"
	"net/http"
)

// func (h *Handler) TestAuthProfileHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	// validation

// 	res, err := h.testAuthProfileHandler(r.Context())
// 	return res, err
// }

func (h *Handler) TestNewAPIKeyHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// validation

	res, err := h.testNewAPIKeyHandler(r.Context())
	return res, err
}

// func (h *Handler) testAuthProfileHandler(ctx context.Context) (*requests.AuthProfile, error) {
// 	// validation

// 	authProfile, err := utils.GetAuthProfileFromCtx(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return authProfile, nil
// }

func (h *Handler) testNewAPIKeyHandler(ctx context.Context) (*requests.APIKey, error) {
	// validation

	apiKey, err := utils.GetAPIKeyFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return apiKey, nil
}
