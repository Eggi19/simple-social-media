package usecases

import "github.com/Eggi19/simple-social-media/repositories"

type UserUsecaseOpts struct {
	UserRepository repositories.UserRepository
}

type UserUsecase interface {
}

type UserUsecaseImpl struct {
	UserRepository repositories.UserRepository
}

func NewUserUsecaseImpl(uuOpts *UserUsecaseOpts) UserUsecase {
	return &UserUsecaseImpl{
		UserRepository: uuOpts.UserRepository,
	}
}
