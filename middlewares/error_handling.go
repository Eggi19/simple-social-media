package middlewares

import (
	"errors"
	"net/http"

	"github.com/Eggi19/simple-social-media/constants"
	"github.com/Eggi19/simple-social-media/custom_errors"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandling(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {

		firstError := c.Errors[0].Err
		errResponse, errCode := checkError(firstError)
		c.AbortWithStatusJSON(errCode, errResponse)
	}
}

func checkError(err error) (dtos.ErrResponse, int) {
	var ve validator.ValidationErrors
	var appErr *custom_errors.AppError

	if errors.As(err, &ve) {
		details := GenerateValidationErrs(ve)
		return dtos.ErrResponse{Message: constants.ValidationErrorMsg, Details: details}, http.StatusBadRequest
	}

	if errors.As(err, &appErr) {
		return dtos.ErrResponse{Status: appErr.Status, Message: appErr.Message}, appErr.Code
	}

	return dtos.ErrResponse{Message: constants.ResponseMsgErrorInternalServer}, http.StatusInternalServerError
}

func GenerateValidationErrs(ve validator.ValidationErrors) []dtos.ValidationErrorMsg {
	details := make([]dtos.ValidationErrorMsg, len(ve))
	for i, fe := range ve {
		details[i] = dtos.ValidationErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
	}
	return details
}

func getErrorMsg(fe validator.FieldError) string {
	containsFe := fe.Param()
	if fe.Param() == `!#$%&'()*+,-./:"\;<=>?@[]^_{|}~` {
		containsFe = "special character"
	}
	if fe.Param() == "abcdefghijklmnopqrstuvwxyz" {
		containsFe = "lowercase"
	}
	if fe.Param() == "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		containsFe = "uppercase"
	}
	if fe.Param() == "1234567890" {
		containsFe = "number"
	}
	if fe.Param() == ` []{|}"\%~#=&<>?/.` {
		containsFe = fe.Param()[1:]
	}
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "use valid email address"
	case "max":
		return "maximum value for this field is " + fe.Param() + " characters"
	case "min":
		return "minimum value for this field is " + fe.Param() + " characters"
	case "lte":
		return "maximum value for this field is " + fe.Param()
	case "gte":
		return "minimum value for this field is " + fe.Param()
	case "len=0|url":
		return "this field just receive empty string or valid url"
	case "len=0|e164":
		return "this field just receive empty string or e.164 phone number format (ex: +1123456789)"
	case "e164":
		return "use e.164 phone number format (ex: +1123456789)"
	case "datetime":
		return "use valid date format: yyyy-mm-dd (ex: 2026-01-02)"
	case "excludes":
		return "this field cannot contain space"
	case "containsany":
		return "this field must contain at least 1 " + containsFe
	case "excludesall":
		return "this field cannot contain space and any of this characters " + containsFe
	case "lowercase":
		return "this field cannot contain uppercase character"
	case "alphanum":
		return "this field just receive alphanumeric"
	default:
		return "unknown error"
	}
}