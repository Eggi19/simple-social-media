package handlers

import (
	"net/http"

	"github.com/Eggi19/simple-social-media/constants"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/usecases"
	"github.com/Eggi19/simple-social-media/utils"
	"github.com/gin-gonic/gin"
)

type TweetHandlerOpts struct {
	TweetUsecase usecases.TweetUsecase
}

type TweetHandler struct {
	TweetUsecase usecases.TweetUsecase
}

func NewTweetHandler(thOpts *TweetHandlerOpts) *TweetHandler {
	return &TweetHandler{
		TweetUsecase: thOpts.TweetUsecase,
	}
}

func (h *TweetHandler) CreateTweet(ctx *gin.Context) {
	var payload dtos.CreateTweetRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	id, err := utils.GetIdParamOrContext(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

    err = h.TweetUsecase.CreateTweet(ctx, int64(id), payload)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Status:  0,
		Message: constants.ResponseMsgPostTweet,
		Data:    nil,
	})
}

func (h *TweetHandler) LikeTweet(ctx *gin.Context) {
	var payload dtos.TweetLikesRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

    err := h.TweetUsecase.LikeTweet(ctx, payload.TweetId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Status:  0,
		Message: constants.ResponseMsgLikeTweet,
		Data:    nil,
	})
}
