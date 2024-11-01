package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/Eggi19/simple-social-media/config"
	"github.com/Eggi19/simple-social-media/constants"
	"github.com/Eggi19/simple-social-media/custom_errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthTokenProvider interface {
	CreateAndSign(data map[string]interface{}) (*JwtToken, error)
	ParseAndVerify(signed string) (jwt.MapClaims, error)
	IsAuthorized(ctx *gin.Context) (bool, *ClaimsData, error)
	GetToken(ctx *gin.Context) (string, error)
}

type JwtProvider struct {
	config config.Config
}

type JwtToken struct {
	AccessToken string `json:"token"`
}

type ClaimsData struct {
	Id   int64
	Role string
}

func NewJwtProvider(config config.Config) AuthTokenProvider {
	return &JwtProvider{
		config: config,
	}
}

func (j *JwtProvider) CreateAndSign(data map[string]interface{}) (*JwtToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  j.config.Issuer,
		"exp":  time.Now().Add(time.Duration(j.config.ExpDurationHour) * time.Hour).Unix(),
		"iat":  time.Now(),
		"data": data,
	})

	signed, err := token.SignedString([]byte(j.config.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	return &JwtToken{
		AccessToken: signed,
	}, nil
}

func (j *JwtProvider) ParseAndVerify(signed string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(signed, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.JwtSecretKey), nil
	}, jwt.WithIssuer(j.config.Issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if err.Error() == "token has invalid claims: token is expired" {
			return nil, custom_errors.BadRequest(err, constants.InvalidAuthTokenErrMsg)
		}
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, custom_errors.InvalidAuthToken()
}

func (j *JwtProvider) IsAuthorized(ctx *gin.Context) (bool, *ClaimsData, error) {
	token, err := j.GetToken(ctx)
	if err != nil {
		return false, nil, err
	}

	claims, err := j.ParseAndVerify(token)
	if err != nil {
		return false, nil, err
	}
	dataMap := claims["data"]
	data, _ := dataMap.(map[string]interface{})

	id := int64(data["id"].(float64))

	if id != 0 {
		return true, &ClaimsData{Id: id}, nil
	}

	return false, nil, custom_errors.InvalidAuthToken()
}

func (j *JwtProvider) GetToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.Request.Header.Get("Authorization")
	t := strings.Fields(authHeader)
	if len(t) == 2 && t[0] == "Bearer" {
		authToken := t[1]
		return authToken, nil
	}

	return "", errors.New("token not found")
}
