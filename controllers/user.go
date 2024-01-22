package controllers

import (
	"ShopService/core"
	_ "ShopService/core"
	"ShopService/helpers"
	"ShopService/models"
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
	user := authUser.(models.User)

	// Fetch updated user details from the database using the controller's DB field
	var updatedUser *models.User
	var err error
	updatedUser, err = models.GetUser(dc.DB, user.ID)
	if err != nil {
		helpers.FailureResponse(context, helpers.APIFailure{Message: helpers.StrPtr("This user does not exist"),
			Status: helpers.IntPtr(400)})
		return
	}
	core.SendToWs("echo", core.Payload{User: *updatedUser, Data: "I just logged in", Event: "LOGGED IN"})

	// Return user details as JSON
	helpers.SuccessResponse(context, helpers.APISuccess{Data: helpers.AnyPtr(user)})
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser creates a user
// @Summary creates a user
// @Tags         users
// @Router /users/ [post]
// @Param request body CreateUserRequest true "user"
func (dc *DBController) RegisterUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		helpers.FailureResponse(context, helpers.APIFailure{
			Message: helpers.StrPtr(err.Error()),
		})
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		helpers.FailureResponse(context, helpers.APIFailure{
			Message: helpers.StrPtr(err.Error()),
		})
		return
	}
	record := dc.DB.Create(&user)
	if record.Error != nil {
		helpers.FailureResponse(context, helpers.APIFailure{
			Message: helpers.StrPtr(record.Error.Error()),
		})
		return
	}
	helpers.SuccessResponse(context, helpers.APISuccess{Message: helpers.StrPtr("Registration is successful"),
		Data: helpers.AnyPtr(user.BasicUserData()), Status: helpers.IntPtr(201)})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LoginAPI allows a user to login
// @Summary allows a user to login
// @Tags         auth
// @Router /auth/login/ [post]
// @Param request body LoginRequest true "Login Credential"
// @Success 200 {object} TokenResponse
func (dc *DBController) LoginAPI(context *gin.Context) {
	var request LoginRequest
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
	accessToken, err := core.GenerateJWT(user, false)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	refreshToken, err := core.GenerateJWT(user, true)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"accessToken": accessToken,
		"refreshToken": refreshToken})
}

type AccessTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// GetAccessTokenAPI get a new access token from your refresh token
// @Summary allows a user to get a new access token from your refresh token
// @Router /auth/accessToken/ [post]
// @Tags         auth
// @Param request body AccessTokenRequest true "Refresh Token"
// @Success 200 {object} TokenResponse
func (dc *DBController) GetAccessTokenAPI(context *gin.Context) {
	var request AccessTokenRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	err, claims := core.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		helpers.FailureResponse(context, helpers.APIFailure{
			Message: helpers.StrPtr(err.Error())})
		return
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
