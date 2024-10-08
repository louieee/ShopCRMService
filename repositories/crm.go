package repositories

import (
	"ShopService/core/rabbitmq"
	"ShopService/core/rabbitmq/payloads"
	"ShopService/helpers"
	"ShopService/models"
	"ShopService/schemas"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"

	RabbitMQHelpers "ShopService/core/rabbitmq/helpers"
)

/*
	create, retrieve, list, update and delete Leads
	create, retrieve, list, update and delete contacts
	filter
	convert lead to deal
	convert contact to lead

*/

func CreateLead(db *gorm.DB, lead models.Lead) (*models.Lead, *helpers.CustomError) {
	var count struct{ Count int }
	db.Raw("select count(*) from leads where title ilike ? and company_id = ?", "%"+lead.Title+"%", lead.CompanyID).Find(&count)
	if count.Count > 0 {
		return nil, helpers.ValidationError("A lead with this title already exist")
	}
	res := db.Create(&lead)
	if res.Error != nil {
		return nil, helpers.ValidationError(res.Error.Error())
	}
	server := rabbitmq.RabbitMQServer
	server = server.Connect(os.Getenv("RABBIT_MQ_HOST"), os.Getenv("REDIS_URL"))
	lead_payload := payloads.LeadPayload{
		Id:              int(lead.ID),
		Title:           lead.Title,
		ContactId:       int(lead.ContactID),
		OwnerId:         int(lead.OwnerID),
		Source:          lead.Source,
		NurturingStatus: lead.NurturingStatus,
		IsDeal:          lead.IsDeal,
		Company:         lead.Company.Name,
		ConversionDate:  lead.ConversionDate,
	}
	lead_payload_string := helpers.StructToString(lead_payload)
	if lead_payload_string != nil {
		payload := RabbitMQHelpers.Payload{
			Action:   "create",
			DataType: "lead",
			Data:     *lead_payload_string,
		}
		server.Publish([]string{"report_queue"}, payload)
	}

	return &lead, nil
}

func RetrieveLead(db *gorm.DB, leadInt uint) (*schemas.LeadResponse, *helpers.CustomError) {
	var leads []schemas.LeadResponse
	var count struct{ Count int }
	query := "select count(*) from leads_view where id = ?"
	res := db.Raw(query, leadInt).Find(&count)
	println("count:: ", count.Count)
	if res.Error != nil || count.Count == 0 {
		return nil, helpers.NotFoundError("This lead does not exist")
	}
	res = db.Raw(strings.Replace(query, "count(*)", "*", 1), leadInt).Find(&leads)
	return &leads[0], nil
}

func LeadList(db *gorm.DB, limit uint, offset uint, filters schemas.FilterLead) (*schemas.LeadListResponse, *helpers.CustomError) {
	var leads []schemas.LeadListItem
	var params []interface{}
	query := "select id, title, contact, budget, budget_currency, is_deal from leads_view"
	if len(filters.CompanyIDs) > 0 {
		for _, val := range filters.CompanyIDs {
			println("Param: ", val)
		}
		query = makeSqlQuery(query, "company_id", "in")
		params = append(params, filters.CompanyIDs)
	}
	if filters.Search != "" {
		query = makeSqlQuery(query, "title", "ilike")
		query = makeOrSqlQuery(query, "description", "ilike")
		query = makeOrSqlQuery(query, "contact", "ilike")
		query = makeOrSqlQuery(query, "owner", "ilike")
		params = append(params, "%"+filters.Search+"%", "%"+filters.Search+"%",
			"%"+filters.Search+"%",
			"%"+filters.Search+"%")
	}
	if len(filters.ContactIDs) > 0 {
		for _, val := range filters.ContactIDs {
			println("Param: ", val)
		}
		query = makeSqlQuery(query, "contact_id", "in")
		params = append(params, filters.ContactIDs)
	}
	if filters.StartPurchaseTimeline != nil && filters.EndPurchaseTimeline != nil {
		startPT := filters.StartPurchaseTimeline
		endPT := filters.EndPurchaseTimeline
		query = makeSqlQuery(query, "purchase_timeline", "between")
		query = fmt.Sprintf("%s and ?", query)
		params = append(params, *startPT, *endPT)
	}
	if filters.StartConversionDate != nil && filters.EndConversionDate != nil {
		startCD := filters.StartConversionDate
		endCD := filters.EndConversionDate
		query = makeSqlQuery(query, "conversion_date", "between")
		query = fmt.Sprintf("%s and ?", query)
		params = append(params, *startCD, *endCD)
	}
	if filters.InterestLevel != nil {
		interestLevel := filters.InterestLevel
		query = makeSqlQuery(query, "interest_level", "=")
		params = append(params, *interestLevel)
	}
	if filters.IsDeal != nil {
		isDeal := filters.IsDeal
		query = makeSqlQuery(query, "is_deal", "=")
		params = append(params, *isDeal)
	}
	if len(filters.NurturingStatuses) > 0 {
		for _, val := range filters.NurturingStatuses {
			println("Param: ", val)
		}
		query = makeSqlQuery(query, "nurturing_status", "in")
		params = append(params, filters.NurturingStatuses)
	}
	if len(filters.OwnerIDs) > 0 {
		for _, val := range filters.OwnerIDs {
			println("Param: ", val)
		}
		query = makeSqlQuery(query, "owner_id", "in")
		params = append(params, filters.OwnerIDs)
	}
	if len(filters.Sources) > 0 {
		for _, val := range filters.Sources {
			println("Param: ", val)
		}
		query = makeSqlQuery(query, "source", "in")
		params = append(params, filters.Sources)
	}
	var resCount struct{ Count int }
	countQuery := strings.Replace(query, "id, title, contact, budget, budget_currency, is_deal",
		"count(*)", 1)
	res := db.Raw(fmt.Sprintf("%s;", countQuery), params...).Find(&resCount)
	query = fmt.Sprintf("%s offset ? limit ? ;", query)
	params = append(params, offset, limit)
	res = db.Raw(query, params...).Find(&leads)
	if res.Error != nil {
		return nil, helpers.ValidationError(res.Error.Error())
	}
	return &schemas.LeadListResponse{
		Count:   resCount.Count,
		Results: leads,
	}, nil
}

func UpdateLead(db *gorm.DB, leadID uint, lead models.Lead) (*models.Lead, *helpers.CustomError) {
	var existingLead models.Lead
	if res := db.First(&existingLead, leadID); res.Error != nil {
		return nil, helpers.NotFoundError("This lead does not exist")
	}
	var count struct{ Count int }
	db.Raw("select count(*) from leads_view where title ilike ? and id != ? and company_id = ?", "%"+lead.Title+"%",
		leadID, lead.CompanyID).Find(&count)
	println("count: ", count.Count)
	if count.Count > 0 {
		return nil, helpers.ValidationError("A lead with this title already exists")
	}
	db.Model(&existingLead).Updates(lead)
	server := rabbitmq.RabbitMQServer
	server = server.Connect(os.Getenv("RABBIT_MQ_HOST"), os.Getenv("REDIS_URL"))
	lead_payload := payloads.LeadPayload{
		Id:              int(lead.ID),
		Title:           lead.Title,
		ContactId:       int(lead.ContactID),
		OwnerId:         int(lead.OwnerID),
		Source:          lead.Source,
		NurturingStatus: lead.NurturingStatus,
		IsDeal:          lead.IsDeal,
		Company:         lead.Company.Name,
		ConversionDate:  lead.ConversionDate,
	}
	lead_payload_string := helpers.StructToString(lead_payload)
	if lead_payload_string != nil {
		payload := RabbitMQHelpers.Payload{
			Action:   "update",
			DataType: "lead",
			Data:     *lead_payload_string,
		}
		server.Publish([]string{"report_queue"}, payload)
	}

	return &existingLead, nil
}

func DeleteLead(db *gorm.DB, leadID uint) *helpers.CustomError {
	lead_res, err := RetrieveLead(db, leadID)
	lead := *lead_res
	if err != nil {
		return helpers.NotFoundError("This lead does not exist")
	}
	res := db.Exec("delete from leads where id = ?", leadID)
	if res.Error != nil {
		return helpers.ValidationError(res.Error.Error())
	}
	server := rabbitmq.RabbitMQServer
	server = server.Connect(os.Getenv("RABBIT_MQ_HOST"), os.Getenv("REDIS_URL"))
	lead_payload := payloads.LeadPayload{
		Id:              int(leadID),
		Title:           lead.Title,
		ContactId:       int(lead.ContactID),
		OwnerId:         int(lead.OwnerID),
		Source:          lead.Source,
		NurturingStatus: lead.NurturingStatus,
		IsDeal:          lead.IsDeal,
		Company:         lead.Company,
		ConversionDate:  &lead.ConversionDate,
	}
	lead_payload_string := helpers.StructToString(lead_payload)
	if lead_payload_string != nil {
		payload := RabbitMQHelpers.Payload{
			Action:   "delete",
			DataType: "lead",
			Data:     *lead_payload_string,
		}
		server.Publish([]string{"report_queue"}, payload)
	}

	return nil
}

func CreateContact(db *gorm.DB, contact *models.Contact) (*models.Contact, *helpers.CustomError) {
	var countRes struct{ Count int }
	db.Raw("select count(*) from contacts where name ilike ? ;", contact.Name).Find(&countRes)
	if countRes.Count > 0 {
		return nil, helpers.ValidationError("A contact with this name already exists")
	}

	res := db.Create(&contact)
	if res.Error != nil {
		return nil, helpers.ValidationError(res.Error.Error())
	}
	server := rabbitmq.RabbitMQServer
	server = server.Connect(os.Getenv("RABBIT_MQ_HOST"), os.Getenv("REDIS_URL"))
	contact_payload := payloads.ContactPayload{
		Id:              int(contact.ID),
		OwnerId: int(contact.OwnerID),
		Name: contact.Name,
		Email: contact.Email,
		PhoneNumber: contact.PhoneNumber,
		Address: contact.Address,
	}
	contact_payload_string := helpers.StructToString(contact_payload)
	if contact_payload_string != nil {
		payload := RabbitMQHelpers.Payload{
			Action:   "create",
			DataType: "contact",
			Data:     *contact_payload_string,
		}
		server.Publish([]string{"report_queue"}, payload)
	}

	return contact, nil
}

func RetrieveContact(db *gorm.DB, contactInt uint) (*models.Contact, *helpers.CustomError) {
	var contact models.Contact
	res := db.First(&contact, contactInt)
	if res.Error != nil {
		return nil, helpers.NotFoundError("This contact does not exist")
	}
	return &contact, nil
}

func ContactList(db *gorm.DB, limit uint, offset uint, filters schemas.FilterContact) (*schemas.ContactListResponse, *helpers.CustomError) {
	var contacts []schemas.ContactItem
	var params []interface{}
	query := "select id, name, email, phone_number, gender from contacts"
	if len(filters.OwnerIDs) > 0 {
		query = makeSqlQuery(query, "owner_id", "in")
		params = append(params, filters.OwnerIDs)
	}
	if filters.Gender != "" {
		query = makeSqlQuery(query, "gender", "ilike")
		params = append(params, "%"+filters.Gender+"%")
	}
	if filters.Country != "" {
		query = makeSqlQuery(query, "gender", "ilike")
		params = append(params, "%"+filters.Country+"%")
	}
	if filters.StartAge > 0 && filters.EndAge > 0 {
		query = makeSqlQuery(query, "age", "between")
		query = fmt.Sprintf("%s and ?", query)
		params = append(params, filters.StartAge, filters.EndAge)
	}
	if filters.State != "" {
		query = makeSqlQuery(query, "state", "ilike")
		params = append(params, "%"+filters.State+"%")
	}
	var result struct{ Count int }
	countQuery := strings.Replace(fmt.Sprintf("%s;", query),
		"id, name, email, phone_number, gender", "count(*)", 1)
	db.Raw(countQuery, params...).Find(&result)
	query = fmt.Sprintf("%s offset ? limit ? ;", query)
	params = append(params, offset, limit)
	db.Raw(query, params...).Find(&contacts)
	return &schemas.ContactListResponse{
		Count:   result.Count,
		Results: contacts,
	}, nil
}

func UpdateContact(db *gorm.DB, contactID uint, contact *models.Contact) (*models.Contact, *helpers.CustomError) {
	var existingContact models.Contact
	if res := db.First(&existingContact, contactID); res.Error != nil {
		return nil, helpers.NotFoundError("This contact does not exist")
	}
	var count struct{ Count int }
	db.Raw("select count(*) from contacts where name ilike ? and id != ?", "%"+contact.Name+"%",
		contactID).Find(&count)
	if count.Count > 0 {
		return nil, helpers.NotFoundError("A contact with this name already exists")
	}
	db.Model(&existingContact).Updates(*contact)
	server := rabbitmq.RabbitMQServer
	server = server.Connect(os.Getenv("RABBIT_MQ_HOST"), os.Getenv("REDIS_URL"))
	contact_obj := *contact
	contact_payload := payloads.ContactPayload{
		Id:              int(contact_obj.ID),
		OwnerId: int(contact_obj.OwnerID),
		Name: contact_obj.Name,
		Email: contact_obj.Email,
		PhoneNumber: contact_obj.PhoneNumber,
		Address: contact_obj.Address,
	}
	contact_payload_string := helpers.StructToString(contact_payload)
	if contact_payload_string != nil {
		payload := RabbitMQHelpers.Payload{
			Action:   "update",
			DataType: "contact",
			Data:     *contact_payload_string,
		}
		server.Publish([]string{"report_queue"}, payload)
	}
	return &existingContact, nil
}

func DeleteContact(db *gorm.DB, contactID uint) *helpers.CustomError {
	contact_res, err := RetrieveContact(db, contactID)
	contact :=  *contact_res
	if err != nil {
		return helpers.NotFoundError("This contact does not exist")
	}
	res := db.Exec("delete from contacts where id = ?", contactID)
	if res.Error != nil {
		return helpers.ValidationError(res.Error.Error())
	}
	server := rabbitmq.RabbitMQServer
	server = server.Connect(os.Getenv("RABBIT_MQ_HOST"), os.Getenv("REDIS_URL"))
	contact_payload := payloads.ContactPayload{
		Id:              int(contact.ID),
		OwnerId: int(contact.OwnerID),
		Name: contact.Name,
		Email: contact.Email,
		PhoneNumber: contact.PhoneNumber,
		Address: contact.Address,
	}
	contact_payload_string := helpers.StructToString(contact_payload)
	if contact_payload_string != nil {
		payload := RabbitMQHelpers.Payload{
			Action:   "delete",
			DataType: "contact",
			Data:     *contact_payload_string,
		}
		server.Publish([]string{"report_queue"}, payload)
	}
	return nil
}

func CreateCompany(db *gorm.DB, company schemas.CreateCompany) *helpers.CustomError {
	var res struct{ Count int }
	db.Raw("select count(*) from companies where name ilike ?", "%"+company.Name+"%").Find(&res)
	if res.Count > 0 {
		return helpers.NotFoundError("A company with this name already exists")
	}
	res2 := db.Create(&models.Company{
		Name:     company.Name,
		Industry: company.Industry,
		Size:     company.Size,
	})
	if res2.Error != nil {
		return helpers.ValidationError(res2.Error.Error())
	}
	return nil
}

func RetrieveCompany(db *gorm.DB, companyInt uint) (*schemas.CompanyItem, *helpers.CustomError) {
	var company schemas.CompanyItem
	db.Raw("select * from companies where id = ? ;", companyInt).Find(&company)

	if company.Id == 0 {
		return nil, helpers.NotFoundError("this company does not exist")
	}
	return &company, nil
}

func CompanyList(db *gorm.DB, limit uint, offset uint, filters schemas.FilterCompany) (*schemas.CompanyListResponse, *helpers.CustomError) {
	var companies []schemas.CompanySchema
	var count struct{ Count int }
	var params []interface{}
	query := "select * from companies"

	if filters.Industry != "" {
		query = makeSqlQuery(query, "industry", "ilike")
		params = append(params, "%"+filters.Industry+"%")
	}
	cQuery := strings.Replace(query, "select *", "select count(*)", 1)
	res := db.Raw(cQuery, params...).Find(&count)
	if res.Error != nil {
		return nil, helpers.ValidationError(res.Error.Error())
	}
	if filters.Search != "" {
		query = makeSqlQuery(query, "search", "ilike")
		params = append(params, "%"+filters.Search+"%")
	}
	query = fmt.Sprintf("%s offset ? limit ?", query)
	println("query: ", offset, "  ", limit)
	params = append(params, offset, limit)
	res = db.Raw(query, params...).Find(&companies)
	if res.Error != nil {
		return nil, helpers.ValidationError(res.Error.Error())
	}
	return &schemas.CompanyListResponse{
		Count:   count.Count,
		Results: companies,
	}, nil
}

func UpdateCompany(db *gorm.DB, companyID uint, company models.Company) (*models.Company, *helpers.CustomError) {
	var existingCompany models.Company
	if res := db.First(&existingCompany, companyID); res.Error != nil {
		println("Error: ", res.Error.Error())
		return nil, helpers.NotFoundError("This company does not exist")
	}
	var count struct{ Count int }
	db.Raw("select count(*) from companies where name ilike ? and id != ?", "%"+company.Name+"%",
		companyID).Find(&count)
	if count.Count > 0 {
		return nil, helpers.NotFoundError("A company with this name already exists")
	}
	db.Model(&existingCompany).Updates(company)
	return &company, nil
}

func DeleteCompany(db *gorm.DB, companyID uint) *helpers.CustomError {
	var query struct{ Count int }
	db.Raw("select count(*) from companies where id=?;", companyID).Find(&query)
	if query.Count == 0 {
		return helpers.NotFoundError("This company does not exist")
	}
	res := db.Exec("delete from companies where id=?;", companyID)
	if res.Error != nil {
		return helpers.ValidationError(res.Error.Error())
	}
	return nil
}
