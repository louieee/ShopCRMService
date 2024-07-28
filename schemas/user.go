package schemas

import (
	"fmt"
	"strconv"
	"strings"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AccessTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	UserType string `json:"user_type" enum:"Administrator,Staff,Customer"`
	UserId   uint   `json:"user_id" gorm:"unique"`
}

func (user *UserResponse) ToTokenPayload() *TokenUserPayload {
	names := strings.Split(user.Name, " ")
	return &TokenUserPayload{
		Id:        user.Id,
		FirstName: names[0],
		LastName:  strings.Join(names[1:], " "),
		UserType:  user.UserType,
		Email:     user.Email,
	}
}

type TokenUserPayload struct {
	Id        int    `json:"id"`
	UserID    string `json:"user_id"`
	UserType  string `json:"user_type" enum:"Administrator,Staff,Customer"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (userPayload TokenUserPayload) ConvertPayloadToUserResponse() *UserResponse {
	Id, _ := strconv.Atoi(userPayload.UserID)
	return &UserResponse{
		Id:       Id,
		Name:     fmt.Sprintf("%s %s", userPayload.FirstName, userPayload.LastName),
		Email:    userPayload.Email,
		UserType: userPayload.UserType,
		UserId:   uint(Id),
	}
}

type FilterUser struct {
	*PageFilter
	Search   string `form:"search"`
	UserType string `form:"user_type" enum:"Administrator,Staff"`
}

type UserListResponse struct {
	Count   int            `json:"count"`
	Results []UserResponse `json:"results"`
}
