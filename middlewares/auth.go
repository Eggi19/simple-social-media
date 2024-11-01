package middlewares

import (
	"net/http"

	"github.com/Eggi19/simple-social-media/config"
	"github.com/Eggi19/simple-social-media/constants"
	"github.com/Eggi19/simple-social-media/custom_errors"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/utils"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(config config.Config) func(*gin.Context) {
	return func(ctx *gin.Context) {
		authorized, data, err := utils.NewJwtProvider(config).IsAuthorized(ctx)
		if !authorized && err != nil && data == nil {
			if err.Error() == custom_errors.InvalidAuthToken().Error() {
				ctx.AbortWithStatusJSON(http.StatusForbidden, dtos.ErrResponse{
					Message: constants.InvalidAuthTokenErrMsg,
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.InvalidAuthTokenErrMsg,
			})
			return
		}
		ctx.Set("data", data)
		ctx.Next()
	}
}