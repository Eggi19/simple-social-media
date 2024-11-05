package handlers

import (
	"net/http"

	"github.com/Eggi19/simple-social-media/constants"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/usecases"
	"github.com/Eggi19/simple-social-media/utils"
	"github.com/gin-gonic/gin"
)

type CommentHandlerOpts struct {
	CommentUsecase usecases.CommentUsecase
}

type CommentHandler struct {
	CommentUsecase usecases.CommentUsecase
}

func NewCommentHandler(thOpts *CommentHandlerOpts) *CommentHandler {
	return &CommentHandler{
		CommentUsecase: thOpts.CommentUsecase,
	}
}

func (h *CommentHandler) CreateComment(ctx *gin.Context) {
	var payload dtos.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	id, err := utils.GetIdParamOrContext(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

    err = h.CommentUsecase.CreateComment(ctx, int64(id), payload)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Status:  0,
		Message: constants.ResponseMsgPostComment,
		Data:    nil,
	})
}
