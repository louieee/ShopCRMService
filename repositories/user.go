package repositories

import (
	"ShopService/helpers"
	"ShopService/models"
	"ShopService/schemas"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

func GetUser(db *gorm.DB, userID uint) (*schemas.UserResponse, *helpers.CustomError) {
	var users []schemas.UserResponse
	db.Raw("select * from users where id = ? ;", userID).Find(&users)
	if len(users) == 0 {
		return nil, helpers.NotFoundError("This user does not exist")
	}
	return &users[0], nil
}

func CreateUser(db *gorm.DB, user models.User) (*models.User, *helpers.CustomError) {
	var count struct{ Count int }
	db.Raw("select count(*) from users where email ilike ? ", "%"+user.Email+"%").Find(&count)
	if count.Count > 0 {
		return nil, helpers.ValidationError("A user with this email address already exist")
	}
	res := db.Create(&user)
	if res.Error != nil {
		return nil, helpers.ValidationError(res.Error.Error())
	}
	return &user, nil
}

func UpdateUser(db *gorm.DB, user models.User) (*models.User, *helpers.CustomError) {
	var existingUser models.User
	var count struct{ Count int }
	db.Raw("select count(*) from users where user_id = ?", user.UserId).Find(&count)
	if count.Count == 0 {
		return nil, helpers.ValidationError("A user with this id does not exist")
	}
	db.Model(&existingUser).Updates(user)
	return &existingUser, nil
}

func DeleteUser(db *gorm.DB, userID uint) *helpers.CustomError {
	res := db.Exec("delete from leads where owner_id = ?", userID)
	if res.Error != nil {
		return helpers.ValidationError(res.Error.Error())
	}
	res = db.Exec("delete from users where user_id = ?", userID)
	if res.Error != nil {
		return helpers.ValidationError(res.Error.Error())
	}
	return nil
}

func UserList(db *gorm.DB, limit uint, offset uint, filters schemas.FilterUser) (*schemas.UserListResponse, *helpers.CustomError) {
	var users []schemas.UserResponse
	var params []interface{}
	query := "select * from users"
	if filters.Search != "" {
		query = makeSqlQuery(query, "name", "ilike")
		query = makeOrSqlQuery(query, "email", "ilike")
		query = makeOrSqlQuery(query, "user_type", "ilike")
		params = append(params, "%"+filters.Search+"%",
			"%"+filters.Search+"%",
			"%"+filters.Search+"%")
	}
	if filters.UserType != "" {
		userType := filters.UserType
		query = makeSqlQuery(query, "user_type", "ilike")
		params = append(params, userType)
	}
	var resCount struct{ Count int }
	countQuery := strings.Replace(query, "*",
		"count(*)", 1)
	res := db.Raw(fmt.Sprintf("%s;", countQuery), params...).Find(&resCount)
	query = fmt.Sprintf("%s offset ? limit ? ;", query)
	params = append(params, offset, limit)
	res = db.Raw(query, params...).Find(&users)
	if res.Error != nil {
		return nil, helpers.ValidationError(res.Error.Error())
	}
	return &schemas.UserListResponse{
		Count:   resCount.Count,
		Results: users,
	}, nil
}
