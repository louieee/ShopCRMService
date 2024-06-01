package repositories

import (
	"ShopService/helpers"
	"ShopService/schemas"
	"github.com/jinzhu/gorm"
)

func GetUser(db *gorm.DB, userID uint) (*schemas.UserResponse, *helpers.CustomError) {
	var users []schemas.UserResponse
	db.Raw("select * from users where id = ? ;", userID).Find(&users)
	if len(users) == 0 {
		return nil, helpers.NotFoundError("This user does not exist")
	}
	return &users[0], nil
}
