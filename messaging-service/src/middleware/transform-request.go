package middleware

import (
	"net/http"
)

type TransformRequest struct {
}

func NewTransformRequest() *TransformRequest {
	return &TransformRequest{}
}

func (a *TransformRequest) execute(h HTTPHandler) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {

		h(w, r)
		return nil, nil
	}
}
