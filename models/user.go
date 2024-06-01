package models

// models/user.go

import (
	"ShopService/helpers"
	"ShopService/schemas"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	UserType string `json:"user_type" enum:"Administrator,Staff,Customer"`
	Password string `json:"password"`
}

func (user *User) BasicUserData() map[string]any {
	return gin.H{"userId": user.ID, "email": user.Email}
}

func (user *User) HashPassword(password string) *helpers.CustomError {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return helpers.ValidationError(err.Error())
	}
	user.Password = string(bytes)
	return nil
}
func (user *User) CheckPassword(providedPassword string) *helpers.CustomError {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return helpers.AuthenticationError("The password is incorrect")
	}
	return nil
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
