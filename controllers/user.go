package controllers

import (
	"ShopService/core"
	_ "ShopService/core"
	"ShopService/helpers"
	"ShopService/repositories"
	"ShopService/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetUser retrieves a user
// @Summary retrieves a user
// @Router /users/me [get]
// @Tags         users
// @Produce      json
// @Success      200  {object} schemas.UserResponse
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) GetUser(context *gin.Context) {
	// Retrieve user ID from the URL parameter
	authUser, _ := context.Get("user")
	user := authUser.(schemas.UserResponse)
	// Fetch updated user details from the database using the controller's DB field

	newUser, err := repositories.GetUser(dc.DB, uint(user.Id))
	if err != nil {
		helpers.FailureResponse(context, *err)
		return
	}
	user = *newUser
	core.SendToWs("echo", core.Payload{User: user, Data: "I just logged in", Event: "LOGGED IN"})

	// Return user details as JSON
	helpers.SuccessResponse(context, helpers.APISuccess{Data: helpers.AnyPtr(user)})
}

// GetAccessTokenAPI get a new access token from your refresh token
// @Summary allows a user to get a new access token from your refresh token
// @Router /auth/accessToken/ [post]
// @Tags         auth
// @Param request body schemas.AccessTokenRequest true "Refresh Token"
// @Success 200 {object} schemas.TokenResponse
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) GetAccessTokenAPI(context *gin.Context) {
	var request schemas.AccessTokenRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	err, claims := core.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		helpers.FailureResponse(context, *helpers.AuthenticationError(err.Error()))
	}
	accessToken, err := core.GenerateJWT(claims.User, false)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"accessToken": accessToken,
		"refreshToken": request.RefreshToken})
}

// GetUserList UserList retrieves all users
// @Summary retrieves all users
// @Router /users/ [get]
// @Param filter query schemas.FilterUser true "filters"
// @Tags         users
// @Produce      json
// @Success      200  {array} schemas.UserResponse[]
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) GetUserList(c *gin.Context) {

	params := schemas.FilterUser{
		Search:   c.Query("search"),
		UserType: c.Query("user_type"),
	}

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	offset := uint((page - 1) * pageSize)
	limit := uint(pageSize)
	users, err := repositories.UserList(dc.DB, limit, offset, params)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Data: helpers.AnyPtr(users), Status: helpers.IntPtr(200)})
}
