package core

import (
	"ShopService/helpers"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			helpers.FailureResponse(context, helpers.APIFailure{
				Message: helpers.StrPtr("request does not contain an access token")})
			return
		}
		claim, err := ValidateAccessToken(tokenString)
		if err != nil {
			helpers.FailureResponse(context, helpers.APIFailure{
				Message: helpers.StrPtr(err.Error())})
			return
		}
		context.Set("user", claim.User)
		context.Next()
	}
}
