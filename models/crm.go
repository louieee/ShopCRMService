package models

import (
	"ShopService/helpers"
	"github.com/jinzhu/gorm"
	"time"
)

type Contact struct {
	gorm.Model

	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
	Age         uint   `json:"age"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	OwnerID     uint   `json:"owner_id"`
}

func (c Contact) IsValid() (error *helpers.CustomError) {
	if len(c.Name) < 3 {
		return helpers.ValidationError("Name must contain more than 3 characters")
	}
	if c.Age < 18 {
		return helpers.ValidationError("Contact must be above 18 years")
	}
	return nil
}

type Company struct {
	gorm.Model

	Name     string `json:"name"`
	Industry string `json:"industry"`
	Size     string `json:"size"`
}

type Lead struct {
	gorm.Model

	Title            string    `json:"title" validate:"min=3"`
	Description      string    `json:"description" validate:"min=3"`
	ContactInfo      Contact   `json:"contact_info"`
	ContactID        uint      `json:"contact_id" validate:"gte=1"`
	OwnerID          uint      `json:"owner_id" validate:"gte=1"`
	InterestLevel    uint      `json:"interest_level"`
	Source           string    `json:"source" `
	PurchaseTimeline time.Time `json:"purchase_timeline"`
	NurturingStatus  string    `json:"nurturing_status"`
	Budget           float64   `json:"budget"`
	BudgetCurrency   string    `json:"budget_currency"`
	IsDeal           bool      `json:"is_deal" default:"true"`
	ConversionDate   time.Time `json:"conversion_date"`
	CompanyID        int       `json:"company_id" validate:"gte=1" `
	Company          Company   `json:"company"`
}
