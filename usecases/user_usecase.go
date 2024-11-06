package usecases

import (
	"context"
	"fmt"
	"strconv"

	"firebase.google.com/go/v4/messaging"
	"github.com/Eggi19/simple-social-media/constants"
	"github.com/Eggi19/simple-social-media/custom_errors"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/repositories"
	"github.com/Eggi19/simple-social-media/utils"
)

type UserUsecaseOpts struct {
	HashAlgorithm           utils.Hasher
	AuthTokenProvider       utils.AuthTokenProvider
	FirebaseMessagingClient *messaging.Client
	Transactor              repositories.Transactor
	UserRepository          repositories.UserRepository
}

type UserUsecase interface {
	RegisterUser(ctx context.Context, req dtos.UserRegisterData) error
	Login(ctx context.Context, req dtos.UserLoginRequest) (*utils.JwtToken, error)
	AddFollower(ctx context.Context, userId int64, req dtos.AddFollowerRequest) error
}

type UserUsecaseImpl struct {
	HashAlgorithm           utils.Hasher
	AuthTokenProvider       utils.AuthTokenProvider
	FirebaseMessagingClient *messaging.Client
	Transactor              repositories.Transactor
	UserRepository          repositories.UserRepository
}

func NewUserUsecaseImpl(uuOpts *UserUsecaseOpts) UserUsecase {
	return &UserUsecaseImpl{
		HashAlgorithm:           uuOpts.HashAlgorithm,
		AuthTokenProvider:       uuOpts.AuthTokenProvider,
		FirebaseMessagingClient: uuOpts.FirebaseMessagingClient,
		Transactor:              uuOpts.Transactor,
		UserRepository:          uuOpts.UserRepository,
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

	_, err = u.Transactor.WithinTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		newUser, err := u.UserRepository.RegisterUser(ctx, *user)
		if err != nil {
			return nil, err
		}

		err = u.UserRepository.AddUserToFirestore(ctx, *newUser)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
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

func (u *UserUsecaseImpl) AddFollower(ctx context.Context, userId int64, req dtos.AddFollowerRequest) error {
	followingId := strconv.Itoa(int(req.FollowingId))
	followerId := strconv.Itoa(int(userId))

	err := u.UserRepository.AddFollowerToFirestore(ctx, followingId, followerId)
	if err != nil {
		return err
	}

	following, err := u.UserRepository.GetUserById(ctx, req.FollowingId)
	if err != nil {
		return err
	}

	follower, err := u.UserRepository.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	if following.FcmToken.Valid {
		message := &messaging.Message{
			Token: following.FcmToken.String,
			Notification: &messaging.Notification{
				Title: fmt.Sprintf(`%s followed you`, follower.Name),
			},
		}
	
		_, err = u.FirebaseMessagingClient.Send(ctx, message)
		if err != nil {
			return err
		}
	}

	return nil
}
