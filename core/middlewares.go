package core

import (
	"ShopService/helpers"
	"github.com/gin-gonic/gin"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			helpers.FailureResponse(context, *helpers.AuthenticationError(
				"request does not contain an access token"))
		}
		if "bearer" != strings.ToLower(tokenString[:6]) {
			helpers.FailureResponse(context, *helpers.AuthenticationError(
				"Token must start with Bearer"))
		}
		tokenString = strings.Split(tokenString, " ")[1]
		claim, err := ValidateAccessToken(tokenString)
		if err != nil {
			helpers.FailureResponse(context, *helpers.AuthenticationError(err.Error()))
			return
		}
		userData := claim.User.ConvertPayloadToUserResponse()
		context.Set("user", *userData)
		context.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(CORSConfig["ALLOWED_ORIGIN"].([]string), ","))
		context.Writer.Header().Set("Access-Control-Allow-Credentials", CORSConfig["ALLOW_CREDENTIALS"].(string))
		context.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(CORSConfig["ALLOWED_HEADERS"].([]string), ","))
		context.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(CORSConfig["ALLOWED_METHODS"].([]string), ","))

		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(204)
			return
		}
		context.Next()
	}
}
