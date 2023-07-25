package serrors

import (
	"errors"
	"net/http"

	goerrors "github.com/go-errors/errors"
)

type Error struct {
	ErrCode int
	Msg     string
	Stack   error
}

func (e Error) Error() string {
	if e.Stack != nil {
		return e.Stack.Error()
	}
	return e.Msg
}

func GetStatusCode(e error) int {
	er, _ := e.(Error)
	return er.ErrCode
}

func GetStackTrace(e error) error {
	er, _ := e.(Error)
	return er.Stack
}

func New(
	errCode int,
	msg string,
	err error,
) Error {

	return Error{
		ErrCode: errCode,
		Msg:     msg,
		Stack:   err,
	}
}

func AuthErrorf(msg string, err error) Error {
	if err == nil {
		err = goerrors.Wrap(errors.New(msg), 0)
	}
	_, ok := err.(*goerrors.Error)
	if !ok {
		err = goerrors.Wrap(err, 0)
	}
	return Error{
		ErrCode: http.StatusUnauthorized,
		Msg:     msg,
		Stack:   err,
	}
}

func AuthError(err error) Error {
	if err == nil {
		err = goerrors.Wrap(errors.New("authorization error"), 0)
	}
	_, ok := err.(*goerrors.Error)
	if !ok {
		err = goerrors.Wrap(err, 0)
	}
	return Error{
		ErrCode: http.StatusUnauthorized,
		Msg:     "Authorization error",
		Stack:   err,
	}
}

func InvalidArgumentErrorf(msg string, err error) Error {
	if err == nil {
		err = goerrors.Wrap(errors.New(msg), 0)
	}
	_, ok := err.(*goerrors.Error)
	if !ok {
		err = goerrors.Wrap(err, 0)
	}
	return Error{
		ErrCode: http.StatusBadRequest,
		Msg:     msg,
		Stack:   err,
	}
}

func InvalidArgumentError(err error) Error {
	if err == nil {
		err = goerrors.Wrap(errors.New("invalid argument"), 0)
	}
	_, ok := err.(*goerrors.Error)
	if !ok {
		err = goerrors.Wrap(err, 0)
	}
	return Error{
		ErrCode: http.StatusBadRequest,
		Msg:     "Invalid argument",
		Stack:   err,
	}
}

func InternalErrorf(msg string, err error) Error {
	if err == nil {
		err = goerrors.Wrap(errors.New(msg), 0)
	}
	_, ok := err.(*goerrors.Error)
	if !ok {
		err = goerrors.Wrap(err, 0)
	}
	return Error{
		ErrCode: http.StatusInternalServerError,
		Msg:     msg,
		Stack:   err,
	}
}
func InternalError(err error) Error {
	if err == nil {
		err = goerrors.Wrap(errors.New("internal error"), 0)
	}
	_, ok := err.(*goerrors.Error)
	if !ok {
		err = goerrors.Wrap(err, 0)
	}
	return Error{
		ErrCode: http.StatusInternalServerError,
		Msg:     "Internal error",
		Stack:   err,
	}
}
