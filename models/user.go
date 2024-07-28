package models

// models/user.go

import (
	"ShopService/schemas"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
)

type User struct {
	gorm.Model
	UserId   uint   `json:"user_id" gorm:"unique"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	UserType string `json:"user_type" enum:"Administrator,Staff,Customer"`
}

func (user *User) BasicUserData() map[string]any {
	return gin.H{"userId": user.ID, "email": user.Email}
}

func (user *User) ToTokenPayload() *schemas.TokenUserPayload {
	names := strings.Split(user.Name, " ")
	return &schemas.TokenUserPayload{
		Id:        int(user.ID),
		FirstName: names[0],
		LastName:  strings.Join(names[1:], " "),
		UserType:  user.UserType,
		Email:     user.Email,
	}
}

func (user *User) ToUserResponse() *schemas.UserResponse {
	return &schemas.UserResponse{
		Id:       int(user.ID),
		Name:     user.Name,
		UserType: user.UserType,
		Email:    user.Email,
	}
}

// models/shop.go
