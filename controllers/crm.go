package controllers

import "C"
import (
	_ "ShopService/core"
	"ShopService/helpers"
	"ShopService/models"
	"ShopService/repositories"
	"ShopService/schemas"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// CreateLead creates a lead
// @Summary creates a lead
// @Router /crm/leads [post]
// @Param lead body schemas.CreateLead True "lead"
// @Tags         leads
// @Produce      json
// @Success      200  {object} schemas.LeadResponse
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) CreateLead(c *gin.Context) {
	authUser, _ := c.Get("user")
	user := authUser.(schemas.UserResponse)
	var lead models.Lead
	if err := c.ShouldBindJSON(&lead); err != nil {
		helpers.FailureResponse(c, *helpers.ValidationError(err.Error()))
		return
	}
	lead.OwnerID = user.UserId
	createdLead, err := repositories.CreateLead(dc.DB, lead)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{Message: helpers.StrPtr("Lead creation is successful"),
		Data: helpers.AnyPtr(createdLead), Status: helpers.IntPtr(201)})

}

// RetrieveLead retrieves a lead
// @Summary retrieves a lead
// @Router /crm/leads/{lead_id} [get]
// @Param lead_id path int True "lead"
// @Tags         leads
// @Produce      json
// @Success      200  {object} schemas.LeadResponse
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) RetrieveLead(c *gin.Context) {
	leadIDParam, _ := strconv.Atoi(c.Param("lead_id"))
	leadID := uint(leadIDParam)
	lead, err := repositories.RetrieveLead(dc.DB, leadID)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Data: helpers.AnyPtr(lead), Status: helpers.IntPtr(200)})

}

// LeadList retrieves all leads
// @Summary retrieves all lead
// @Router /crm/leads/ [get]
// @Param filter query schemas.FilterLead true "filters"
// @Tags         leads
// @Produce      json
// @Success      200  {array} schemas.LeadResponse[]
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) LeadList(c *gin.Context) {

	params := schemas.FilterLead{
		Search:                c.Query("search"),
		OwnerIDs:              ConvertStringSliceToUintSlice(strings.Split(c.Query("owner_ids"), ",")),
		ContactIDs:            ConvertStringSliceToUintSlice(strings.Split(c.Query("contact_ids"), ",")),
		InterestLevel:         ConvertStringToUintPtr(c.Query("interest_level")),
		Sources:               ConvertStringToSlice(c.Query("sources")),
		StartPurchaseTimeline: ConvertStringToTime(c.Query("start_purchase_timeline"), false),
		EndPurchaseTimeline:   ConvertStringToTime(c.Query("end_purchase_timeline"), false),
		NurturingStatuses:     ConvertStringToSlice(c.Query("nurturing_statuses")),
		IsDeal:                ConvertStringToBool(c.Query("is_deal")),
		StartConversionDate:   ConvertStringToTime(c.Query("start_conversion_date"), false),
		EndConversionDate:     ConvertStringToTime(c.Query("end_conversion_date"), false),
		CompanyIDs:            ConvertStringSliceToUintSlice(strings.Split(c.Query("company_ids"), ",")),
	}

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	offset := uint((page - 1) * pageSize)
	limit := uint(pageSize)
	leads, err := repositories.LeadList(dc.DB, limit, offset, params)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Data: helpers.AnyPtr(leads), Status: helpers.IntPtr(200)})
}

// UpdateLead updates a lead
// @Summary updates a lead
// @Router /crm/leads/{lead_id} [put]
// @Param lead body schemas.CreateLead True "lead"
// @Param lead_id path int True "lead id"
// @Tags         leads
// @Produce      json
// @Success      200  {object} schemas.LeadResponse
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) UpdateLead(c *gin.Context) {
	leadIDParam, _ := strconv.Atoi(c.Param("lead_id"))
	var requestBody models.Lead
	err2 := c.ShouldBindJSON(&requestBody)
	if err2 != nil {
		helpers.FailureResponse(c, *helpers.ValidationError(err2.Error()))
		return
	}
	leadID := uint(leadIDParam)
	_, err := repositories.UpdateLead(dc.DB, leadID, requestBody)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Status: helpers.IntPtr(200)})

}

// DeleteLead deletes a lead
// @Summary deletes a lead
// @Router /crm/leads/{lead_id} [delete]
// @Param lead_id path int True "lead id"
// @Tags         leads
// @Produce      json
// @Success      204  {object} nil
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) DeleteLead(c *gin.Context) {
	leadIDParam, _ := strconv.Atoi(c.Param("lead_id"))
	leadID := uint(leadIDParam)
	err := repositories.DeleteLead(dc.DB, leadID)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Message: helpers.StrPtr("Lead deleted successfully"), Status: helpers.IntPtr(204)})

}

// CreateContact creates a contact
// @Summary creates a contact
// @Router /crm/contacts/ [post]
// @Param contact body schemas.CreateContact True "contact"
// @Tags         contacts
// @Produce      json
// @Success      201  {object} schemas.ContactResponse
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) CreateContact(c *gin.Context) {
	var contact models.Contact
	err := c.ShouldBindJSON(&contact)
	if err != nil {
		helpers.FailureResponse(c, *helpers.ValidationError(err.Error()))
		return
	}
	err2 := contact.IsValid()
	if err2 != nil {
		helpers.FailureResponse(c, *err2)
		return
	}

	createdContact, err2 := repositories.CreateContact(dc.DB, &contact)
	if err2 != nil {
		helpers.FailureResponse(c, *err2)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{Message: helpers.StrPtr("Contact creation is successful"),
		Data: helpers.AnyPtr(createdContact), Status: helpers.IntPtr(201)})
}

// RetrieveContact retrieves a contact
// @Summary retrieves a contact
// @Router /crm/contacts/{contact_id} [get]
// @Param contact_id path int True "contact id"
// @Tags         contacts
// @Produce      json
// @Success      200  {object} schemas.ContactResponse
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) RetrieveContact(c *gin.Context) {
	contactIDParam, _ := strconv.Atoi(c.Param("contact_id"))
	contactID := uint(contactIDParam)
	contact, err := repositories.RetrieveContact(dc.DB, contactID)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Data: helpers.AnyPtr(contact), Status: helpers.IntPtr(200)})
}

// ContactList retrieves the list of contacts
// @Summary retrieves the list of contacts
// @Router /crm/contacts [get]
// @Param filter query schemas.FilterContact False "filters"
// @Tags         contacts
// @Produce      json
// @Success      200  {array} schemas.ContactResponse[]
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) ContactList(c *gin.Context) {
	startAge, _ := strconv.Atoi(c.Query("start_age"))
	endAge, _ := strconv.Atoi(c.Query("end_age"))
	params := schemas.FilterContact{
		Gender:   c.Query("gender"),
		StartAge: startAge,
		EndAge:   endAge,
		State:    c.Query("state"),
		Country:  c.Query("country"),
		OwnerIDs: ConvertStringSliceToIntSlice(strings.Split(c.Query("owner_id"), ",")),
	}

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	offset := uint((page - 1) * (pageSize))
	limit := uint(pageSize)
	contacts, err := repositories.ContactList(dc.DB, limit, offset, params)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Data: helpers.AnyPtr(contacts), Status: helpers.IntPtr(200)})
}

// UpdateContact updates a contact
// @Summary updates a contact
// @Router /crm/contacts/{contact_id} [put]
// @Param contact body schemas.CreateContact True "contact"
// @Param contact_id path int True "contact id"
// @Tags         contacts
// @Produce      json
// @Success      200  {object} schemas.ContactResponse
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) UpdateContact(c *gin.Context) {
	contactIDParam, _ := strconv.Atoi(c.Param("contact_id"))
	var requestBody models.Contact
	err1 := c.ShouldBindJSON(&requestBody)
	if err1 != nil {
		helpers.FailureResponse(c, *helpers.ValidationError(err1.Error()))
	}
	contactID := uint(contactIDParam)
	contact, err := repositories.UpdateContact(dc.DB, contactID, &requestBody)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Data: helpers.AnyPtr(contact), Status: helpers.IntPtr(200)})
}

// DeleteContact deletes a contact
// @Summary deletes a contact
// @Router /crm/contacts/{contact_id} [delete]
// @Param contact_id path int True "contact id"
// @Tags         contacts
// @Produce      json
// @Success      204  {object} nil
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) DeleteContact(c *gin.Context) {
	contactIDParam, _ := strconv.Atoi(c.Param("contact_id"))
	contactID := uint(contactIDParam)
	err := repositories.DeleteContact(dc.DB, contactID)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Message: helpers.StrPtr("Contact deleted successfully"), Status: helpers.IntPtr(204)})
}

// CreateCompany creates a company
// @Summary creates a company
// @Router /crm/companies [post]
// @Param company body schemas.CreateCompany True "company"
// @Tags         companies
// @Produce      json
// @Success      201  {object} schemas.CompanyItem
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) CreateCompany(c *gin.Context) {
	var company schemas.CreateCompany
	err := c.ShouldBindJSON(&company)
	if err != nil {
		helpers.FailureResponse(c, *helpers.ValidationError(err.Error()))
		return
	}
	err2 := repositories.CreateCompany(dc.DB, company)
	if err != nil {
		helpers.FailureResponse(c, *err2)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{Message: helpers.StrPtr("Company creation is successful"),
		Status: helpers.IntPtr(201)})
}

// RetrieveCompany retrieves a company
// @Summary retrieves a company
// @Router /crm/companies/{company_id} [get]
// @Param company_id path int True "company id"
// @Tags         companies
// @Produce      json
// @Success      200  {object} schemas.CompanyItem
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) RetrieveCompany(c *gin.Context) {
	companyIDParam, _ := strconv.Atoi(c.Param("company_id"))
	companyID := uint(companyIDParam)
	company, err := repositories.RetrieveCompany(dc.DB, companyID)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Data: &company, Status: helpers.IntPtr(200)})
}

// CompanyList retrieves the list of companies
// @Summary retrieves the list of companies
// @Router /crm/companies/ [get]
// @Param filter query schemas.FilterCompany false "Filters"
// @Tags         companies
// @Produce      json
// @Success      200  {array} schemas.CompanyListResponse
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) CompanyList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	params := schemas.FilterCompany{
		Industry: c.Query("industry"),
		Search:   c.Query("search"),
	}
	offset := uint((page - 1) * pageSize)
	limit := uint(pageSize)

	companies, err := repositories.CompanyList(dc.DB, limit, offset, params)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Data: helpers.AnyPtr(companies), Status: helpers.IntPtr(200)})
}

// UpdateCompany updates a company details
// @Summary updates a company details
// @Router /crm/companies/{company_id} [put]
// @Param company_id path int true "Company Id"
// @Param company body schemas.CreateCompany true "Company"
// @Tags         companies
// @Produce      json
// @Success      200  {object} schemas.CompanyItem
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) UpdateCompany(c *gin.Context) {
	companyIDParam, _ := strconv.Atoi(c.Param("company_id"))
	var requestBody models.Company
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		helpers.FailureResponse(c, *helpers.ValidationError(err.Error()))
		return
	}
	companyID := uint(companyIDParam)
	company, err2 := repositories.UpdateCompany(dc.DB, companyID, requestBody)
	if err != nil {
		helpers.FailureResponse(c, *err2)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Data: helpers.AnyPtr(company), Status: helpers.IntPtr(200)})
}

// DeleteCompany deletes a company
// @Summary deletes a company
// @Router /crm/companies/{company_id} [delete]
// @Param company_id path int true "Company Id"
// @Tags         companies
// @Produce      json
// @Success      204 {object} nil
// @Failure      400  {object}  helpers.APIFailure
// @Failure      404  {object} helpers.APIFailure
func (dc *DBController) DeleteCompany(c *gin.Context) {
	companyIDParam, _ := strconv.Atoi(c.Param("company_id"))
	companyID := uint(companyIDParam)
	err := repositories.DeleteCompany(dc.DB, companyID)
	if err != nil {
		helpers.FailureResponse(c, *err)
		return
	}
	helpers.SuccessResponse(c, helpers.APISuccess{
		Message: helpers.StrPtr("Company deleted successfully"), Status: helpers.IntPtr(204)})
}
