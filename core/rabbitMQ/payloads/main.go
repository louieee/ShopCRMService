package payloads

import "time"


type LeadPayload struct{
	Id int `json:"id"`
	Title string `json:"title"`
	ContactId int `json:"contact_id"`
	OwnerId int `json:"owner_id"`
	Source string `json:"source"`
	NurturingStatus string `json:"nurturing_status"`
	IsDeal bool `json:"is_deal"`
	Company string `json:"company"`
	ConversionDate *time.Time `json:"conversion_date"`
}

type ContactPayload struct{
	Id int `json:"id"`
	OwnerId int `json:"owner_id"`
	Name string `json:"name"`
	Email string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address string `json:"address"`
}