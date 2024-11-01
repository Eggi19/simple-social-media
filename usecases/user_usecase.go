package usecases

import (
	"context"

	"github.com/Eggi19/simple-social-media/constants"
	"github.com/Eggi19/simple-social-media/custom_errors"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/repositories"
	"github.com/Eggi19/simple-social-media/utils"
)

type UserUsecaseOpts struct {
	HashAlgorithm     utils.Hasher
	AuthTokenProvider utils.AuthTokenProvider
	UserRepository    repositories.UserRepository
}

type UserUsecase interface {
	RegisterUser(ctx context.Context, req dtos.UserRegisterData) error
	Login(ctx context.Context, req dtos.UserLoginRequest) (*utils.JwtToken, error)
}

type UserUsecaseImpl struct {
	HashAlgorithm     utils.Hasher
	AuthTokenProvider utils.AuthTokenProvider
	UserRepository    repositories.UserRepository
}

func NewUserUsecaseImpl(uuOpts *UserUsecaseOpts) UserUsecase {
	return &UserUsecaseImpl{
		HashAlgorithm:     uuOpts.HashAlgorithm,
		AuthTokenProvider: uuOpts.AuthTokenProvider,
		UserRepository:    uuOpts.UserRepository,
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

func (u *UserUsecaseImpl) Login(ctx context.Context, req dtos.UserLoginRequest) (*utils.JwtToken, error) {
	user, err := u.UserRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, custom_errors.Unauthorized(err, constants.InvalidCredentialsErrMsg, 103)
	}

	isCorrectPassword, err := u.HashAlgorithm.CheckPassword(req.Password, []byte(user.Password))
	if !isCorrectPassword {
		return nil, custom_errors.Unauthorized(err, constants.InvalidCredentialsErrMsg, 103)
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap["id"] = user.Id

	accessToken, err := u.AuthTokenProvider.CreateAndSign(dataTokenMap)
	if err != nil {
		return nil, err
	}

	tokens := &utils.JwtToken{
		AccessToken: accessToken.AccessToken,
	}

	return tokens, nil
}
