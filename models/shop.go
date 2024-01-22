package models

import (
	"github.com/jinzhu/gorm"
)

type Shop struct {
	gorm.Model
	Title       string
	Description string
	UserID      uint // Foreign key referencing User's ID
}
