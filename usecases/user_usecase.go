package usecases

import (
	"context"

	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/repositories"
	"github.com/Eggi19/simple-social-media/utils"
)

type UserUsecaseOpts struct {
	HashAlgorithm  utils.Hasher
	UserRepository repositories.UserRepository
}

type UserUsecase interface {
	RegisterUser(ctx context.Context, req dtos.UserRegisterData) error
}

type UserUsecaseImpl struct {
	HashAlgorithm  utils.Hasher
	UserRepository repositories.UserRepository
}

func NewUserUsecaseImpl(uuOpts *UserUsecaseOpts) UserUsecase {
	return &UserUsecaseImpl{
		HashAlgorithm:  uuOpts.HashAlgorithm,
		UserRepository: uuOpts.UserRepository,
	}
}

func (u *UserUsecaseImpl) RegisterUser(ctx context.Context, req dtos.UserRegisterData) error {
	pwd := req.Password
	pwdHash, err := u.HashAlgorithm.HashPassword(pwd)
	if err != nil {
		return err
	}

	req.Password = string(pwdHash)

	user := dtos.ConvertUserRegisterData(&req)

	err = u.UserRepository.RegisterUser(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}
