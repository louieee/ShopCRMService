package handlers

import (
	"ShopService/core"
	"ShopService/models"
	"ShopService/repositories"
	"encoding/json"
	"fmt"
)




type UserPayload struct {
	Id          uint   `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	IsCustomer  bool   `json:"is_customer"`
	IsAdmin     bool   `json:"is_admin"`
	IsStaff     bool   `json:"is_staff"`
	DisplayName string `json:"display_name"`
	ProfilePic  string `json:"profile_pic"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}

func (user *UserPayload) toUserModel() (userObj models.User) {
	userType := ""
	if user.IsAdmin {
		userType = "Administrator"
	} else if user.IsStaff {
		userType = "Staff"
	} else {
		userType = "Customer"
	}
	return models.User{
		UserId:   user.Id,
		Name:     fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Email:    user.Email,
		UserType: userType,
	}
}

func HandleUser(action string, data string) {
	switch action {
	case "create":
		handleNewUser(data)
	case "update":
		handleUpdateUser(data)
	case "delete":
		handleDeleteUser(data)

	}
}

func handleNewUser(data string) {
	var user UserPayload
	db := core.GetDB()
	err1 := json.Unmarshal([]byte(data), &user)
	if err1 != nil {
		panic(fmt.Sprintf("Error unwrapping message: %s", err1.Error()))
	}
	newUser := user.toUserModel()
	if newUser.UserType != "Customer" {
		repositories.CreateUser(db, newUser)
	}
}

func handleDeleteUser(data string) {
	var user UserPayload
	db := core.GetDB()
	err1 := json.Unmarshal([]byte(data), &user)
	if err1 != nil {
		panic(fmt.Sprintf("Error unwrapping message: %s", err1.Error()))
	}
	newUser := user.toUserModel()
	repositories.DeleteUser(db, newUser.UserId)
}

func handleUpdateUser(data string) {
	var user UserPayload
	db := core.GetDB()
	err1 := json.Unmarshal([]byte(data), &user)
	if err1 != nil {
		panic(fmt.Sprintf("Error unwrapping message: %s", err1.Error()))
	}
	newUser := user.toUserModel()
	repositories.UpdateUser(db, newUser)
}
