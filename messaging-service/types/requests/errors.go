package requests

import (
	"encoding/json"
	"net/http"
)

func MakeUnauthorized(w http.ResponseWriter, msg string) (interface{}, error) {
	errResponse := ErrorResponse{
		Message: msg,
	}
	bytes, err := json.Marshal(errResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return w.Write([]byte("could not send error response"))
	}

	w.WriteHeader(http.StatusUnauthorized)
	return w.Write(bytes)
}

func MakeInternalError(w http.ResponseWriter, msg string) (interface{}, error) {
	errResponse := ErrorResponse{
		Message: msg,
	}
	bytes, err := json.Marshal(errResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return w.Write([]byte("could not send error response"))
	}

	w.WriteHeader(http.StatusInternalServerError)
	return w.Write(bytes)
}

