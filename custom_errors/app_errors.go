package custom_errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Eggi19/simple-social-media/constants"
)

var (
	ErrInvalidAuthToken = errors.New(constants.InvalidAuthTokenErrMsg)
	ErrNotFound         = errors.New(constants.NotFoundErrorMsg)
	ErrContextNotFound  = errors.New(constants.ContextNotFoundErrMsg)
	ErrForbidden        = errors.New(constants.ForbiddenErrMsg)
)

type AppError struct {
	Status  int
	Code    int
	Message string
	err     error
}

func (e AppError) Error() string {
	return fmt.Sprint(e.Message)
}

func BadRequest(err error, message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		err:     err,
	}
}

func InternalServerError(err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: constants.ResponseMsgErrorInternalServer,
		err:     err,
	}
}

func InvalidAuthToken() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: constants.InvalidAuthTokenErrMsg,
		err:     ErrInvalidAuthToken,
	}
}

func NotFound() *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: constants.NotFoundErrorMsg,
		err:     ErrNotFound,
	}
}

func Unauthorized(err error, message string, status int) *AppError {
	return &AppError{
		Status:  status,
		Code:    http.StatusUnauthorized,
		Message: message,
		err:     err,
	}
}

func ContextNotFound() *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: constants.ContextNotFoundErrMsg,
		err:     ErrContextNotFound,
	}
}

func Forbidden() *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: constants.ForbiddenErrMsg,
		err:     ErrForbidden,
	}
}
