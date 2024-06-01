package controllers

import (
	"ShopService/core"
	_ "ShopService/core"
	"ShopService/helpers"
	"ShopService/models"
	"ShopService/repositories"
	"ShopService/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUser retrieves a user
// @Summary retrieves a user
// @Router /users/ [get]
// @Tags         users
// @Produce      json
// @Success      200  {object} models.BasicUserDataType
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) GetUser(context *gin.Context) {
	// Retrieve user ID from the URL parameter
	authUser, _ := context.Get("user")
	user := authUser.(schemas.UserResponse)
	println("user_id: ", user.Id, "  user_name: ", user.Name,
		"  user_email: ", user.Email, "  user_type: ", user.UserType)

	// Fetch updated user details from the database using the controller's DB field
	_, err := repositories.GetUser(dc.DB, uint(user.Id))
	if err != nil {
		helpers.FailureResponse(context, *err)
		return
	}
	core.SendToWs("echo", core.Payload{User: user, Data: "I just logged in", Event: "LOGGED IN"})

	// Return user details as JSON
	helpers.SuccessResponse(context, helpers.APISuccess{Data: helpers.AnyPtr(user)})
}

// RegisterUser creates a user
// @Summary creates a user
// @Tags         users
// @Router /users/ [post]
// @Param request body schemas.CreateUserRequest true "user"
// @Produce      json
// @Success      201  {object} schemas.UserResponse
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) RegisterUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		helpers.FailureResponse(context, *helpers.ValidationError(err.Error()))
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		helpers.FailureResponse(context, *err)
		return
	}
	record := dc.DB.Create(&user)
	if record.Error != nil {
		helpers.FailureResponse(context,
			*helpers.ValidationError(record.Error.Error()))
		return
	}
	helpers.SuccessResponse(context, helpers.APISuccess{Message: helpers.StrPtr("Registration is successful"),
		Data: helpers.AnyPtr(user.BasicUserData()), Status: helpers.IntPtr(201)})
}

// LoginAPI allows a user to login
// @Summary allows a user to login
// @Tags         auth
// @Router /auth/login/ [post]
// @Param request body schemas.LoginRequest true "Login Credential"
// @Success 200 {object} schemas.TokenResponse
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) LoginAPI(context *gin.Context) {
	var request schemas.LoginRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// check if email exists and password is correct
	record := dc.DB.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}
	userPayload := user.ToTokenPayload()
	accessToken, err := core.GenerateJWT(*userPayload, false)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	refreshToken, err := core.GenerateJWT(*userPayload, true)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"accessToken": accessToken,
		"refreshToken": refreshToken})
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
