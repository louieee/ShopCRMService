package schemas

import (
	"time"
)

/*

	nurturing_statuses = [
    "New Lead",
    "Engaged Lead",
    "Contacted Lead",
    "Qualified Lead",
    "Proposal Sent",
    "Negotiating",
    "Pending Contract/Agreement",
    "Closed - Won",
    "Closed - Lost",
    "On Hold",
    "Re-Engagement"
]
lead_sources = [
    "Website Form",
    "Social Media",
    "Email Campaign",
    "Referral",
    "Event",
    "Cold Call",
    "Content Marketing",
    "Advertisement",
    "Networking",
    "Trade Show",
    "Word of Mouth"
]

*/

type PageFilter struct {
	Page     int `json:"page" validate:"required,gte=1" default:"1"`
	PageSize int `json:"page_size" validate:"required,gte=1" default:"10"`
}

type FilterLead struct {
	*PageFilter
	Search        string   `form:"search"`
	OwnerIDs      []uint   `form:"owner_ids,omitempty"`
	ContactIDs    []uint   `form:"contact_ids,omitempty"`
	InterestLevel *uint    `form:"interest_level,omitempty" validate:"gte=0" default:"0"`
	Sources       []string `form:"sources,omitempty"`
	// filters by purchase time
	StartPurchaseTimeline *time.Time `form:"start_purchase_timeline,omitempty" default:"2006-01-02T15:04:05Z"`
	EndPurchaseTimeline   *time.Time `form:"end_purchase_timeline,omitempty" default:"2006-01-02T15:04:05Z"`
	NurturingStatuses     []string   `form:"nurturing_statuses,omitempty" enums:"New Lead, Engaged Lead, Contacted Lead, Qualified Lead, Proposal Sent, Negotiating, Pending Contract/Agreement, Closed - Won, Closed - Lost, On Hold, Re-Engagement"`
	IsDeal                *bool      `form:"is_deal,omitempty"`
	StartConversionDate   *time.Time `form:"start_conversion_date,omitempty" default:"2006-01-02T15:04:05Z"`
	EndConversionDate     *time.Time `form:"end_conversion_date,omitempty" default:"2006-01-02T15:04:05Z"`
	CompanyIDs            []uint     `form:"company_ids,omitempty"`
}

type FilterContact struct {
	*PageFilter
	Gender   string `json:"gender"`
	StartAge int    `json:"start_age" validate:"gte=18"`
	EndAge   int    `json:"end_age" validate:"gte=18" `
	State    string `json:"state"`
	Country  string `json:"country"`
	OwnerIDs []int  `json:"owner_ids"`
}

type FilterCompany struct {
	*PageFilter
	Search   string `json:"search"`
	Industry string `json:"industry"`
}

type LeadListItem struct {
	Id             uint    `json:"id"`
	Title          string  `json:"title"`
	Contact        string  `json:"contact"`
	Budget         float64 `json:"budget"`
	BudgetCurrency string  `json:"budgetCurrency"`
	IsDeal         bool    `json:"isDeal"`
}

type LeadResponse struct {
	Id               uint      `json:"id"`
	Title            string    `json:"title"`
	ContactID        uint      `json:"contact_id"`
	Contact          string    `json:"contact"`
	OwnerID          uint      `json:"owner_id"`
	Owner            string    `json:"owner"`
	InterestLevel    uint      `json:"interest_level"`
	Source           string    `json:"source"`
	PurchaseTimeline time.Time `json:"purchase_timeline"`
	NurturingStatus  string    `json:"nurturing_status"`
	Budget           float64   `json:"budget"`
	BudgetCurrency   string    `json:"budget_currency"`
	IsDeal           bool      `json:"is_deal"`
	ConversionDate   time.Time `json:"conversion_date"`
	CompanyID        int       `json:"company_id"`
	Company          string    `json:"company"`
	CreatedAt        string    `json:"createdAt"`
	UpdatedAt        string    `json:"updatedAt"`
}

type LeadListResponse struct {
	Count   int            `json:"count"`
	Results []LeadListItem `json:"results"`
}

type ContactResponse struct {
	Id          uint   `json:"id"`
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

type ContactItem struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
}

type ContactListResponse struct {
	Count   int           `json:"count"`
	Results []ContactItem `json:"results"`
}

type CompanyItem struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Industry string `json:"industry"`
	Size     string `json:"size"`
}

type CompanyListResponse struct {
	Count   int             `json:"count"`
	Results []CompanySchema `json:"results"`
}
type CompanySchema struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Industry string `json:"industry"`
}

type CreateLead struct {
	Title         string `json:"title" validate:"min=3""`
	Description   string `json:"description" validate:"min=3"`
	ContactID     uint   `json:"contact_id" validate:"gte=1" json:"contact_id,omitempty"`
	InterestLevel uint   `json:"interest_level" json:"interest_level,omitempty"`
	Source        string `json:"source" enums:"Website Form, Social Media, Email Campaign, Referral, Event, Cold Call, Content Marketing, Advertisement, Networking, Trade Show, Word of Mouth" json:"source,omitempty"`
	//ProductIDs       []uint    `json:"product_ids"`
	PurchaseTimeline time.Time `json:"purchase_timeline" json:"purchase_timeline" default:"2020-01-02T15:04:05Z"`
	NurturingStatus  string    `json:"nurturing_status" enums:"New Lead, Engaged Lead, Contacted Lead, Qualified Lead, Proposal Sent, Negotiating, Pending Contract/Agreement, Closed - Won, Closed - Lost, On Hold, Re-Engagement" json:"nurturing_status,omitempty"`
	Budget           float64   `json:"budget" json:"budget,omitempty"`
	BudgetCurrency   string    `json:"budget_currency" enums:"usd,ngn,gbp" json:"budget_currency,omitempty"`
	IsDeal           bool      `json:"is_deal" json:"is_deal,omitempty" default:"false"`
	ConversionDate   time.Time `json:"conversion_date" json:"conversion_date" default:"2020-01-02T15:04:05Z"`
	CompanyID        int       `json:"company_id" validate:"gte=1" json:"company_id,omitempty"`
}

type CreateContact struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
	Age         uint   `json:"age" validate:"gt=5"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	OwnerID     uint   `json:"owner_id"`
}

type CreateCompany struct {
	Name     string `json:"name"`
	Industry string `json:"industry"`
	Size     string `json:"size"`
}
