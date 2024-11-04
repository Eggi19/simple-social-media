package handlers

import (
	"net/http"

	"github.com/Eggi19/simple-social-media/constants"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/usecases"
	"github.com/gin-gonic/gin"
)

type FirebaseHandlerOpts struct {
	FirebaseUsecase usecases.FirebaseUsecase
}

type FirebaseHandler struct {
	FirebaseUsecase usecases.FirebaseUsecase
}

func NewFirebaseHandler(fhOpts *FirebaseHandlerOpts) *FirebaseHandler {
	return &FirebaseHandler{
		FirebaseUsecase: fhOpts.FirebaseUsecase,
	}
}

func (h *FirebaseHandler) SubscribeTopic(ctx *gin.Context) {
	var payload dtos.SubscribeTopicRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}
	
    err := h.FirebaseUsecase.SubscribeTopic(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	
	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Status:  0,
		Message: constants.ResponseMsgSubscribeTopic,
		Data:    nil,
	})
}

func (h *FirebaseHandler) UnsubscribeTopic(ctx *gin.Context) {
	var payload dtos.UnsubscribeTopicRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	err := h.FirebaseUsecase.UnsubscribeTopic(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Status:  0,
		Message: constants.ResponseMsgUnsubscribeTopic,
		Data:    nil,
	})
}
