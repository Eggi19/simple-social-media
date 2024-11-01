package handlers

import "github.com/Eggi19/simple-social-media/usecases"

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
