package core

import (
	"ShopService/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	var err error
	config := fmt.Sprintf(""+
		"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", DBConfig["HOST"],
		DBConfig["PORT"], DBConfig["USER"], DBConfig["NAME"],
		DBConfig["SSL_MODE"], DBConfig["PASSWORD"])
	db, err = gorm.Open("postgres", config)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{}, &models.Lead{}, &models.Contact{}, &models.Company{})
	return db
}
