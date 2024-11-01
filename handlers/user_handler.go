package handlers

import (
	"net/http"

	"github.com/Eggi19/simple-social-media/constants"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/usecases"
	"github.com/gin-gonic/gin"
)

type UserHandlerOpts struct {
	UserUsecase usecases.UserUsecase
}

type UserHandler struct {
	UserUsecase usecases.UserUsecase
}

func NewUserHandler(uhOpts *UserHandlerOpts) *UserHandler {
	return &UserHandler{
		UserUsecase: uhOpts.UserUsecase,
	}
}

func (h *UserHandler) RegisterUser(ctx *gin.Context) {
	var payload dtos.UserRegisterData
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	err := h.UserUsecase.RegisterUser(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Status:  0,
		Message: constants.ResponseMsgRegistered,
		Data:    nil,
	})
}
