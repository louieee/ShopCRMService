package controllers

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type DBController struct {
	DB *gorm.DB
}

func NewController(db *gorm.DB) *DBController {
	return &DBController{DB: db}
}

// # Trying generics
//
// ConvertToUint prints slices using the [fmt.Println] function.
// The current implementation prints the following slices:
//   - []int{}
//   - []string{}
//
// For more information about Go doc comments, see [Go Doc Comments] at tip.golang.org.
//
// [Go Doc Comments]: https://tip.golang.org/doc/comment
func ConvertToUint(s string) uint {
	// Convert string to uint (assuming it's a uint in the model)
	// Handle errors appropriately in a production environment
	// This is a simplified example, you might want to use strconv.ParseUint
	result, _ := strconv.ParseUint(s, 10, 64)
	return uint(result)
}
