package validation

import (
	"errors"
	"fmt"
	"messaging-service/types/requests"

	"github.com/go-playground/validator"
)

func ValidateGetRoomsUserUUID(req *requests.GetRoomsByUserUUIDRequest) error {
	validate := validator.New()
	err := validate.Struct(req)
	errMsg := ""

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsg += fmt.Sprintf("%s type (%s) is %s", err.Field(), err.Kind(), err.Tag())
			// fmt.Println(err.Namespace())
			// fmt.Println(err.Field())
			// fmt.Println(err.StructNamespace())
			// fmt.Println(err.StructField())
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())
			// fmt.Println()
		}

	}
	if errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}

func ValidateGetMessagesByRoomUUID(req *requests.GetMessagesByRoomUUIDRequest) error {
	validate := validator.New()
	err := validate.Struct(req)
	errMsg := ""

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsg += fmt.Sprintf("%s type (%s) is %s", err.Field(), err.Kind(), err.Tag())
		}

	}
	if errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}

func ValidateRequest(req interface{}) error {
	validate := validator.New()
	err := validate.Struct(req)
	errMsg := ""

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsg += fmt.Sprintf("%s type (%s) is %s", err.Field(), err.Kind(), err.Tag())
		}

	}
	if errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}
